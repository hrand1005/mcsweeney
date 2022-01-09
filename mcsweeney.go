package main

/* TODO:
- logging
- Goroutines
*/

import (
	"fmt"
	"log"
	"mcsweeney/config"
	"mcsweeney/content"
	"mcsweeney/db"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	CONCAT  = "concat.mp4"
	FINAL   = "final.mp4"
	ENCODED = "encoded.txt"
)

func main() {
	// load configs from command line arg
	c, err := config.LoadConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// init db
	dbIntf, err := db.New(c.Name + ".db")
	if err != nil {
		log.Fatal(err)
	}

	// init getter
	query := content.Query(c.Source.Query)
	getIntf, err := content.NewGetter(c.Source.Platform, c.Source.Credentials, query)
	if err != nil {
		log.Fatal(err)
	}

	// get new clips with a retry strategy
	tries := 0
	clips := make([]*content.Clip, 0, c.Source.Query.First+2)
	for len(clips) < c.Source.Query.First {
		tries++
		fmt.Printf("Have: %v, Want: %v\nGetting more content.\n", len(clips), c.Source.Query.First)
		dirtyContent, err := getIntf.Get()
		if err != nil {
			log.Fatal(err)
		}
		if len(dirtyContent) == 0 {
			fmt.Println("Content getter dry...")
			break
		}

		for _, v := range dirtyContent {
			exists, err := dbIntf.Exists(v)
			if err != nil {
				log.Fatal(err)
			}
			if !exists && len(clips) < c.Source.Query.First {
				valid := filter(v, c.Filters)
				if valid {
					clips = append(clips, v)
				}
			}
		}
	}

	if len(clips) == 0 {
		fmt.Println("Unable to find new content.\nExiting...")
		return
	}
	fmt.Printf("Was able to retrieve %v content objects.\nNumber of tries: %v\n", len(clips), tries)

	// create composite video object from clips
	video := &content.Video{}
	for _, v := range clips {
		video.Append(v)
	}

	// check for intro, create and append to video if applicable
	if c.Intro != (config.Intro{}) {
		intro := &content.Intro{
			Path:     c.Intro.Path,
			Duration: c.Intro.Duration,
		}
		video.Prepend(intro)
	}

	// check for outro, create and append to video if applicable
	if c.Outro != (config.Outro{}) {
		outro := &content.Outro{
			Path:     c.Outro.Path,
			Duration: c.Outro.Duration,
		}
		video.Append(outro)
	}

	// clean up existing files
	removeTempFiles()

	// create encoder, encode video components
	encoder := &content.Encoder{Path: ENCODED}
	err = video.Accept(encoder)
	if err != nil {
		log.Fatal(err)
	}

	// concatenate encoded components into one mp4 file
	concatCmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", ENCODED, CONCAT)
	err = concatCmd.Run()
	if err != nil {
		fmt.Printf("Command string\n%s\n", concatCmd)
		log.Fatal(err)
	}

	// create video overlay
	overlayer := &content.Overlayer{
		Font:       c.Options.Overlay.Font,
		Background: c.Options.Overlay.Background,
	}
	err = video.Accept(overlayer)
	if err != nil {
		log.Fatal(err)
	}

	// apply overlay with dateOverlay
	args := append([]string{"-i", CONCAT}, overlayer.Slice()...)
	dateOverlay := createDateOverlay(c.Intro, c.Source.Query.Days)
	args[len(args)-1] = args[len(args)-1] + "," + dateOverlay
	args = append(args, FINAL)
	overlayCmd := exec.Command("ffmpeg", args...)

	fmt.Println("Applying overlay")
	err = overlayCmd.Run()
	if err != nil {
		fmt.Printf("Command string\n%s\n", overlayCmd)
		log.Fatal(err)
	}

	// generate description for video
	describer := &content.Describer{}
	err = video.Accept(describer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Video's description:\n%s\n", describer)

	shareIntf, err := content.NewSharer(c.Destination.Platform, c.Destination.Credentials)
	if err != nil {
		log.Fatal(err)
	}

	// prepare payload to be shared
	payload := content.Payload{
		Title:       c.Destination.Title,
		Path:        FINAL,
		Description: c.Destination.Description + describer.String(), // prepends custom description from config file
		CategoryID:  c.Destination.CategoryID,
		Keywords:    c.Destination.Keywords,
		Privacy:     string(c.Destination.Privacy),
	}

	// share payload, check that the token cache hasn't expired
	statusCode, err := shareIntf.Share(payload)
	if err != nil {
		fmt.Printf("Couldn't share content: %v\n", err)
		fmt.Println("Retrying after clearing token cache...")
		os.Remove(c.Destination.TokenCache)
		statusCode, err = shareIntf.Share(payload)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Content shared!\nResponse Code: %v\n", statusCode)

	// insert data for retrieved components to db
	for _, v := range clips {
		err := dbIntf.Insert(v)
		if err != nil {
			log.Fatal(err)
		}
	}

	return
}

func createDateOverlay(i config.Intro, days int) string {
	// create timerange string
	const prettyDateFull = "January 2, 2006"
	now := time.Now()
	last := now.AddDate(0, 0, -1*days)
	nowFormatted := now.Format(prettyDateFull)
	lastFormatted := last.Format(prettyDateFull)
	// use abbreviated version if years match
	if now.Year() == last.Year() {
		lastFormatted = strings.Split(last.Format(prettyDateFull), ",")[0]
	}
	timeRange := lastFormatted + " - " + nowFormatted
	fmt.Printf("Range: %s\n", timeRange)
	// escape , which is invalid in ffmpeg drawtext
	escapeText := strings.ReplaceAll(timeRange, `,`, `\\\,`)

	return fmt.Sprintf("drawtext=enable='between(t,%f,%f)':fontfile=%s:text=%s:fontsize=112:fontcolor=ffffff:x=(w-text_w)/2:y=(h-text_h)/2", i.OverlayStart, i.Duration, i.Font, escapeText)
}

func filter(c *content.Clip, f config.Filters) bool {
	for _, v := range f.Blacklist {
		if c.Broadcaster == v {
			return false
		}
	}
	return c.Language == f.Language
}

func removeTempFiles() {
	// No need to check errors, as they may appear if no temp files exist
	cmd := exec.Command("/bin/sh", "./cleanup.sh")
	cmd.Run()
}
