package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hrand1005/mcsweeney/twitch"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
)

var twitchConf = flag.String("twitch-config", "", "Path to twitch scraper configuration file")

const clipScraperTimeout = time.Second * 3

func main() {
	flag.Parse()
	if *twitchConf == "" {
		flag.Usage()
		return
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env variables for API access: %v", err)
	}

	tConf, err := LoadTwitchConfig(*twitchConf)
	if err != nil {
		log.Fatalf("Encountered error loading twitch config: " + err.Error())
	}
	clipScraper, err := ConstructTwitchScraper(tConf)
	if err != nil {
		log.Fatalf("Encountered error constructing twitch scraper: " + err.Error())
	}

	clipFilter := twitch.ClipFilter(func(c helix.Clip) bool {
		return true
	})
	clipChan := make(chan helix.Clip, 10)
	doneChan := make(chan bool, 1)

	go clipScraper.Scrape(clipFilter, clipChan, doneChan)

	select {
	case clip := <-clipChan:
		log.Printf("Scraper returned a clip: %+v", clip)
	case <-time.After(clipScraperTimeout):
		log.Println("Timed out waiting for clip. Sending done signal...")
		doneChan <- true
	}

	log.Println("Finished.")
}

func ConstructTwitchScraper(conf twitchConfig) (twitch.Scraper, error) {
	cOpts := &helix.Options{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	c, err := helix.NewClient(cOpts)
	if err != nil {
		return nil, err
	}
	// TODO: don't generate new access token every time
	// request app access token

	resp, err := c.RequestAppAccessToken(nil)
	if err != nil {
		return nil, fmt.Errorf("ConstructTwitchScraper: couldn't request app access token: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ConstructTwitchScraper: requesting app access token returned response with status code %v, err: %v", resp.StatusCode, resp.ErrorMessage)
	}

	c.SetAppAccessToken(resp.Data.AccessToken)

	q := helix.ClipsParams{
		GameID: conf.GameID,
		First:  conf.First,
		// start date -- counts backwards from 'days' in config
		StartedAt: helix.Time{
			Time: time.Now().AddDate(0, 0, -1*conf.Days),
		},
	}

	return twitch.NewScraper(c, q)
}
