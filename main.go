package main

import (
	"flag"
	// "fmt"
	"log"
	"strings"
	"time"

	"github.com/hrand1005/mcsweeney/twitch"
	"github.com/hrand1005/mcsweeney/video"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
)

var env = flag.String("env", "", "Path to file defining environment variables, may be overwritten")
var twitchConf = flag.String("twitch-config", "", "Path to twitch scraper configuration file")

const clipScraperTimeout = time.Second * 5

func main() {
	flag.Parse()
	if *env == "" || *twitchConf == "" {
		flag.Usage()
		return
	}

	// load twitch environment variables from env file
	err := godotenv.Load(*env)
	if err != nil {
		log.Fatalf("failed to load env variables for API access: %v", err)
	}
	// write updated credentials back to env file in case tokens expired/updated
	defer godotenv.Write(twitch.Credentials(), *env)

	tConf, err := LoadTwitchConfig(*twitchConf)
	if err != nil {
		log.Fatalf("Encountered error loading twitch config: " + err.Error())
	}

	clipScraper, err := ConstructTwitchScraper(tConf)
	if err != nil {
		log.Fatalf("Encountered error constructing twitch scraper: " + err.Error())
	}

	// define filter and channels to be used by the scraper
	clipFilter := twitch.ClipFilter(func(c helix.Clip) bool {
		return true
	})
	clipChan := make(chan helix.Clip)
	doneChan := make(chan bool)

	go clipScraper.Scrape(clipFilter, clipChan, doneChan)

	// keep a slice of clip mp4s to create an outfile
	clips := make([]helix.Clip, 0, 10)
	clipMP4s := make([]string, 0, 10)

	// first 5 clips meeting criteria
	for i := 0; i < 5; i++ {
		select {
		case clip := <-clipChan:
			log.Printf("Scraper returned a clip: %+v", clip)
			clips = append(clips, clip)
			cURL := strings.SplitN(clip.ThumbnailURL, "-preview", 2)[0] + ".mp4"
			clipMP4s = append(clipMP4s, cURL)
		case <-time.After(clipScraperTimeout):
			log.Println("Timed out waiting for clip. Sending done signal...")
			doneChan <- true
		}
	}
	if err := video.ConcatenateMP4Files(clipMP4s, "vidout.mp4"); err != nil {
		log.Printf("Encountered error writing video to file: %v", err)
	}
	log.Printf("Generated description for video:\n%s", DescriptionFromTwitchClips(clips, 0))
	log.Println("Finished.")
}

func ConstructTwitchScraper(conf twitchConfig) (twitch.Scraper, error) {
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
		return nil, err
	}

	return twitch.NewScraper(c, query)
}
