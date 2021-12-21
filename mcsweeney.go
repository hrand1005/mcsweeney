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

//TODO: maybe we don't need the entire packages?
import (
	"fmt"
	"log"
	"mcsweeney/config"
	"mcsweeney/content"
	"mcsweeney/db"
	"os"
)

const (
	RawVidsDir       = "tmp/raw"
	ProcessedVidsDir = "tmp/processed"
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

	dirtyContent, err := getIntf.Get()
	if err != nil {
		fmt.Println("Couldn't get new content.")
		log.Fatal(err)
	}

	contentObjs := make([]*content.Content, 0, len(dirtyContent))
	for _, v := range dirtyContent {
		exists, err := dbIntf.Exists(v)
		if err != nil {
			fmt.Println("Couldn't check exists for dbIntf.")
			log.Fatal(err)
		}
		if !exists {
			contentObjs = append(contentObjs, v)
		}
	}
	if len(contentObjs) == 0 {
		fmt.Println("No new content found.\nExiting...")
		return
	}

	compiledVid, err := content.Compile(contentObjs, "compiled-vid.mp4")
	if err != nil {
		fmt.Println("Couldn't compile content.")
		log.Fatal(err)
	}

	final, err := content.ApplyOverlay(contentObjs, c.Options, compiledVid.Path)
	if err != nil {
		fmt.Println("Couldn't apply overlay.")
		log.Fatal(err)
	}
	fmt.Println("Final output in file: ", final)

	shareIntf, err := content.NewSharer(c.Destination.Platform, c.Destination.Credentials)
	if err != nil {
		fmt.Println("Couldn't create content-sharer.")
		log.Fatal(err)
	}

	err = shareIntf.Share(&content.Content{
		Path:        final,
		Title:       c.Destination.Title,
		Description: c.Destination.Description,
		Keywords:    c.Destination.Keywords,
		Privacy:     c.Destination.Privacy,
	})
	if err != nil {
		fmt.Println("Couldn't share content.")
		log.Fatal(err)
	}

	fmt.Println("Content shared successfully!")

	// TODO: table / data for uploaded videos that can be updated at a later
	// time with analytics
	for _, v := range contentObjs {
		err := dbIntf.Insert(v)
		if err != nil {
			fmt.Println("Couldn't insert to dbIntf.")
			log.Fatal(err)
		}
	}

	return
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
