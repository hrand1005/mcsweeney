package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	// "google.golang.org/api/youtube/v3"
	"github.com/hrand1005/mcsweeney/twitch"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
)

var config = flag.String("config", "", "Path to configuration file defining mcsweeney options")
var maxEncoders = flag.Int("max-encoders", 1, "Maximum number of video encodings that can occur concurrently")
var twitchCredentials = flag.String("twitch-credentials", "", "Path to file defining twitch credentials as environmnet variables, may be overwritten")
var youtubeCredentials = flag.String("youtube-credentials", "", "Path to file defining youtube credentials as environmnet variables")

const (
	clipScraperTimeout = time.Second * 5
	outVideo           = "vidout.mp4"
)

func main() {
	flag.Parse()
	if *twitchCredentials == "" || *config == "" {
		flag.Usage()
		return
	}

	err := godotenv.Load(*twitchCredentials, *youtubeCredentials)
	if err != nil {
		log.Fatalf("Failed to load env variables for API access: %v", err)
	}
	// write updated credentials back to env file in case tokens expired/updated
	defer godotenv.Write(twitch.Credentials(), *twitchCredentials)

	tConf, err := LoadTwitchConfig(*config)
	if err != nil {
		log.Fatal("Encountered error loading twitch config: " + err.Error())
	}

	db, err := initClipDB(tConf.DB)
	if err != nil {
		log.Fatal("Encountered error initializng database: " + err.Error())
	}

	clipChan, scrapeCancel := startScrapingService(tConf, db)
	mp4Chan, encodeReportChan, encodeCancel := startEncodingService()

	clipTargetCount := tConf.First
	// mp4ToClipData maps clip mp4s to their clip metadata, which is useful
	// for referencing clip metadata associated with the video file
	mp4ToClip := make(map[string]helix.Clip, clipTargetCount)
	encodeJobs := 0

scrape:
	for i := 0; i < clipTargetCount; i++ {
		select {
		case clip := <-clipChan:
			// transform the ThumbnailURL to get the raw mp4
			cURL := strings.SplitN(clip.ThumbnailURL, "-preview", 2)[0] + ".mp4"
			mp4ToClip[cURL] = clip
			mp4Chan <- cURL
			encodeJobs += 1
		case <-time.After(clipScraperTimeout):
			log.Println("Timed out waiting for clip. Sending done signal...")
			break scrape
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
			encodedClip := mp4ToClip[rep.Input]
			concatenator.AppendMKVFile(rep.Output)
			db.Insert(encodedClip)
			descriptionBuilder.AppendClip(encodedClip)
		}
	}
	encodeCancel <- true

	if err := concatenator.Concatenate(outVideo); err != nil {
		log.Fatalf("Encountered error writing video to file: %v", err)
	}
	defer cleanup()

	desc := descriptionBuilder.Result()
	log.Printf("Generated description for video:\n%s", desc)

	/*
		ytClient, err := NewYoutubeClient()
		if err != nil {
			log.Fatalf("Encountered error building youtube client: %v", err)
		}

		ytVideo := &youtube.Video{
			Snippet: &youtube.VideoSnippet{
				Title: "McSweeney Title",
				Description: desc,
			},
			Status: &youtube.VideoStatus{PrivacyStatus: "private"},
		}

		resp, err := ytClient.UploadVideo(outVideo, ytVideo)
		if err != nil {
			log.Fatalf("Encountered error uploading video: %v", err)
		}

		log.Printf("Uploading Video yielded HTTP Response:\n%#v\nStatus Code: %v", resp, resp.HTTPStatusCode)
	*/
}

// initClipDB initializes the clipDB from the given file. If the file doesn't
// exist, a new one is created.
func initClipDB(f string) (*clipDB, error) {
	// create db from conf, use it to enforce rules of the filter func
	handle, err := sqliteDB(f)
	if err != nil {
		log.Fatal("Encountered error initializng sqlite handle: " + err.Error())
	}
	return newClipDB(handle)
}

// startScrapingService initializes the twitch scraper with a client with the
// given configuration and clipDB for filtering. Returns two channels, one
// which sends clips, and a channel for canceling the scraping service.
func startScrapingService(conf twitchConfig, db *clipDB) (<-chan helix.Clip, chan<- bool) {
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
		if db.Exists(c) {
			return false
		}
		return true
	})
	doneChan := make(chan bool)
	clipChan := clipScraper.Scrape(clipFilter, doneChan)

	return clipChan, doneChan
}

// startEncodingService initializes a new MP4ToMKVEncoder. Returns a channel
// to push mp4s for processing, a channel for recieving results of encodings,
// and a channel to push a cancel signal for the encoding service.
func startEncodingService() (chan<- string, <-chan EncodeReport, chan<- bool) {
	mp4Encoder := NewMP4ToMKVEncoder(*maxEncoders) /*encoding options*/

	// mp4Chan is the channel the mp4Encoder expects to recieve mp4 filepaths
	// from, and then encode them
	mp4Chan := make(chan string)
	encodeDoneChan := make(chan bool)
	encodeReportChan := mp4Encoder.Encode(mp4Chan, encodeDoneChan)

	return mp4Chan, encodeReportChan, encodeDoneChan
}

// cleanup is meant to destroy intermediate files used in the video compilation.
// cleanup should only be called if video compilation didn't fail, otherwise
// keep intermediate files around for debugging.
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
