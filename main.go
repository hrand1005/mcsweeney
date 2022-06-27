package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/hrand1005/mcsweeney/twitch"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
)

var twitchConf = flag.String("twitch-config", "", "Path to twitch scraper configuration file")
var tokensFile = flag.String("tokens", "", "Path to file containing app tokens, tokens be overwritten")

const clipScraperTimeout = time.Second * 5

func main() {
	flag.Parse()
	if *twitchConf == "" || *tokensFile == "" {
		flag.Usage()
		return
	}

	if err := godotenv.Load(".env", *tokensFile); err != nil {
		log.Fatalf("failed to load env variables for API access: %v", err)
	}

	tConf, err := LoadTwitchConfig(*twitchConf)
	if err != nil {
		log.Fatalf("Encountered error loading twitch config: " + err.Error())
	}

	// TODO: maybe we can eliminate mutliple uses of tokensFile? 
	clipScraper, err := ConstructTwitchScraper(tConf, *tokensFile)
	if err != nil {
		log.Fatalf("Encountered error constructing twitch scraper: " + err.Error())
	}

	clipFilter := twitch.ClipFilter(func(c helix.Clip) bool {
		return true
	})
	clipChan := make(chan helix.Clip)
	doneChan := make(chan bool)

	go clipScraper.Scrape(clipFilter, clipChan, doneChan)

	// first 5 clips meeting criteria
	for i := 0; i < 5; i++ {
		select {
		case clip := <-clipChan:
			log.Printf("Scraper returned a clip: %+v", clip)
		case <-time.After(clipScraperTimeout):
			log.Println("Timed out waiting for clip. Sending done signal...")
			doneChan <- true
		}
	}
	log.Println("Finished.")
}

func ConstructTwitchScraper(conf twitchConfig, tokenFile string) (twitch.Scraper, error) {
	cOpts := &helix.Options{
		ClientID:     os.Getenv(twitch.ClientIDEnvKey),
		ClientSecret: os.Getenv(twitch.ClientSecretEnvKey),
	}

	query := helix.ClipsParams{
		GameID: conf.GameID,
		First:  conf.First,
		// start date -- counts backwards from 'days' in config
		StartedAt: helix.Time{
			Time: time.Now().AddDate(0, 0, -1*conf.Days),
		},
	}

	return twitch.NewScraper(cOpts, query, tokenFile)
}
