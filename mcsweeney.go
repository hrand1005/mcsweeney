package main

/* TODO:
- Consistent error handling
- Testing
- Architecture
    - Think about params, pointers, interfaces, etc.
- Goroutines
- Decide on a standard for where to use pointers vs struct values
- Url or URL, not both
- Content obj methods? Or limit info passed around?
- Encode videos consistently
- Take data streams and do things with them for faster editing?
- content should not rely on how config happens to be implemented
    - content should be functional as a standalone package
*/

import (
	"fmt"
	"log"
	"mcsweeney/config"
	"mcsweeney/content"
	"mcsweeney/db"
	"os"
	"os/exec"
)

func main() {
	// TODO: add command line parsing
	c, err := config.LoadConfig(os.Args[1])
	if err != nil {
		fmt.Println("Couldn't load config.")
		log.Fatal(err)
	}

	dbIntf, err := db.New("mcsweeney-test.db")
	if err != nil {
		fmt.Println("Couldn't create content-db.")
		log.Fatal(err)
	}

	query := content.Query(c.Source.Query)
	getIntf, err := content.NewGetter(c.Source.Platform, c.Source.Credentials, query)
	if err != nil {
		fmt.Println("Couldn't create content-getter.")
		log.Fatal(err)
	}

	tries := 0
	clips := make([]*content.Clip, 0, c.Source.Query.First+2)
	for len(clips) < c.Source.Query.First {
		tries++
		fmt.Printf("Have: %v, Want: %v\nGetting more content.\n", len(clips), c.Source.Query.First)
		dirtyContent, err := getIntf.Get()
		if err != nil {
			fmt.Println("Couldn't get new content.")
			log.Fatal(err)
		}
		if len(dirtyContent) == 0 {
			fmt.Println("Content getter dry...")
			break
		}

		for _, v := range dirtyContent {
			exists, err := dbIntf.Exists(v)
			if err != nil {
				fmt.Println("Couldn't check exists for dbIntf.")
				log.Fatal(err)
			}
			if !exists && len(clips) < c.Source.Query.First {
				valid := Filter(v, c.Filters)
				if valid {
					// Log this...
					clips = append(clips, v)
				}
			}
			// Log this...
			/*else {
			    fmt.Printf("Content exists: %s\n", v.Url)
			}*/
		}
	}

	if len(clips) == 0 {
		fmt.Println("Unable to find new content.\nExiting...")
		return
	}
	// Log this...
	fmt.Printf("Was able to retrieve %v content objects.\n", len(clips))
	fmt.Println("Number of tries: ", tries)

	// create composite video object from clips
	video := &content.Video{
		Description: c.Destination.Description,
	}

	// append the clips to the video
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

	encoder := &content.Encoder{Path: "encoded.txt"}
	fmt.Printf("About to encode video components...")
	video.Accept(encoder)
	describer := &content.Describer{}
	video.Accept(describer)
	fmt.Printf("Video's description:\n%s\n", describer.String())

	/*
		compiledVid, err := content.Concatenate(contentObjs, "compiled-vid.mp4")
		if err != nil {
			fmt.Println("Couldn't compile content.")
			log.Fatal(err)
		}

		err = compiledVid.ApplyOverlay(contentObjs, c.Options.Overlay)
		if err != nil {
			fmt.Println("Couldn't apply overlay.")
			log.Fatal(err)
		}

		shareIntf, err := content.NewSharer(c.Destination.Platform, c.Destination.Credentials)
		if err != nil {
			fmt.Println("Couldn't create content-sharer.")
			log.Fatal(err)
		}

		// set final Content object's fields with config args
		compiledVid.Title = c.Destination.Title
		compiledVid.Description = c.Destination.Description + compiledVid.Description // appends the default credits description
		compiledVid.Keywords = c.Destination.Keywords
		compiledVid.Privacy = c.Destination.Privacy

		err = shareIntf.Share(compiledVid)
		if err != nil {
			fmt.Println("Couldn't share content.")
			os.Remove(c.Destination.TokenCache)
			fmt.Println("Retrying after clearing token cache...")
			err = shareIntf.Share(compiledVid)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("Content shared successfully!")
	*/

	// TODO: table / data for uploaded videos that can be updated at a later
	// time with analytics
	for _, v := range clips {
		err := dbIntf.Insert(v)
		if err != nil {
			fmt.Println("Couldn't insert to dbIntf.")
			log.Fatal(err)
		}
	}

	return
}

// Filter checks whether the given content object passes all filters. If
// yes, returns true, else false
func Filter(c *content.Clip, f config.Filters) bool {
	//TODO: find a way to iterate through all filters
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

// Consider this in the final version
/*
func main(){
    var get, edit, share

    // parse command line flags
    if -g then get = getContent(); etc.
    eachStrategy(get, edit, share)
}
*/
