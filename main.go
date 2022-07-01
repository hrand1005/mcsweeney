package main

import (
	"flag"
	"os"
	"path/filepath"

	"log"
	"strings"
	"time"

	"github.com/hrand1005/mcsweeney/twitch"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
)

var env = flag.String("env", "", "Path to file defining environment variables, may be overwritten")
var twitchConf = flag.String("twitch-config", "", "Path to twitch scraper configuration file")

const (
	clipScraperTimeout = time.Second * 5
	clipTargetCount    = 4
	encoderPool        = 3
	outVideo           = "vidout.mp4"
)

func main() {
	flag.Parse()
	if *env == "" || *twitchConf == "" {
		flag.Usage()
		return
	}

	err := godotenv.Load(*env)
	if err != nil {
		log.Fatalf("Failed to load env variables for API access: %v", err)
	}
	// write updated credentials back to env file in case tokens expired/updated
	defer godotenv.Write(twitch.Credentials(), *env)

	tConf, err := LoadTwitchConfig(*twitchConf)
	if err != nil {
		log.Fatalf("Encountered error loading twitch config: " + err.Error())
	}

	// start twitch clip scraping service
	clipChan, scrapeCancel := startScrapingService(tConf)

	// start mp4 encoding service
	mp4Chan, encodeReportChan, encodeCancel := startEncodingService()

	// mp4ToClipData maps clip mp4s to their clip metadata, which is useful
	// for referencing clip metadata associated with the video file
	mp4ToClip := make(map[string]helix.Clip, clipTargetCount)
	encodeJobs := 0
	for i := 0; i < clipTargetCount; i++ {
		select {
		case clip := <-clipChan:
			// log.Printf("Scraper returned a clip:\n%+v\n", clip)
			cURL := strings.SplitN(clip.ThumbnailURL, "-preview", 2)[0] + ".mp4"
			mp4ToClip[cURL] = clip
			mp4Chan <- cURL
			encodeJobs += 1
		case <-time.After(clipScraperTimeout):
			log.Println("Timed out waiting for clip. Sending done signal...")
			break
		}
	}
	scrapeCancel <- true

	// concatenate videos and build a description as encodeReports are returned
	concatenator := NewMKVToMP4Concatenator()
	descriptionBuilder := NewDescriptionBuilder()
	for i := 0; i < encodeJobs; i++ {
		rep := <-encodeReportChan
		if rep.Err != nil {
			log.Printf("Encountered error encoding %s:\nErr: %v\n", rep.Input, rep.Err)
		} else {
			log.Printf("Successfully encoded %s to %s\n", rep.Input, rep.Output)
			concatenator.AppendMKVFile(rep.Output)
			descriptionBuilder.AppendClip(mp4ToClip[rep.Input])
		}
	}
	encodeCancel <- true

	if err := concatenator.Concatenate(outVideo); err != nil {
		log.Fatalf("Encountered error writing video to file: %v", err)
	}
	defer cleanup()
	log.Printf("Generated description for video:\n%s", descriptionBuilder.Result())
}

func startScrapingService(conf twitchConfig) (<-chan helix.Clip, chan<- bool) {
	query := helix.ClipsParams{
		GameID: conf.GameID,
		First:  conf.First,
		// start date -- counts backwards from 'days' in config
		StartedAt: helix.Time{
			Time: time.Now().AddDate(0, 0, -1*conf.Days),
		},
	}

	c, err := twitch.NewClient()
	if err != nil {
		log.Fatalf("Failed to create twitch client: " + err.Error())
	}
	clipScraper, err := twitch.NewScraper(c, query)
	if err != nil {
		log.Fatalf("Encountered error constructing twitch scraper: " + err.Error())
	}

	// define filter and channels to be used by the scraper
	clipFilter := twitch.ClipFilter(func(c helix.Clip) bool {
		return true
	})
	doneChan := make(chan bool)
	clipChan := clipScraper.Scrape(clipFilter, doneChan)

	return clipChan, doneChan
}

func startEncodingService() (chan<- string, <-chan EncodeReport, chan<- bool) {
	mp4Encoder := NewMP4ToMKVEncoder(encoderPool) /*encoding options*/

	// mp4Chan is the channel the mp4Encoder expects to recieve mp4 filepaths
	// from, and then encode them
	mp4Chan := make(chan string, encoderPool)
	encodeDoneChan := make(chan bool)
	encodeReportChan := mp4Encoder.Encode(mp4Chan, encodeDoneChan)

	return mp4Chan, encodeReportChan, encodeDoneChan
}

// cleanup is meant to destroy intermediate files used in the video compilation.
// cleanup should only be called if video compilation didn't fail, otherwise
// keep intermediate files around for debugging. Currently hardcodes files to
// remove.
func cleanup() {
	intFiles, _ := filepath.Glob("*.mkv")
	for _, v := range intFiles {
		os.Remove(v)
	}
	listFiles, _ := filepath.Glob("*.txt")
	for _, v := range listFiles {
		os.Remove(v)
	}
}
