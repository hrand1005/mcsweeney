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
        
    tries := 0
    contentObjs := make([]*content.Content, 0, c.Source.Query.First)
    for len(contentObjs) < c.Source.Query.First {
        tries ++
        fmt.Printf("Have: %v, Want: %v\nGetting more content.\n", len(contentObjs), c.Source.Query.First)
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
            if !exists && len(contentObjs) < c.Source.Query.First {
                // Log this...
                contentObjs = append(contentObjs, v)
            } 
            // Log this...
            /*else {
                fmt.Printf("Content exists: %s\n", v.Url)
            }*/
        }
    }

    if len(contentObjs) == 0 {
        fmt.Println("Unable to find new content.\nExiting...")
        return
    }
    // Log this...
    fmt.Printf("Was able to retrieve %v content objects.\n", len(contentObjs))
    fmt.Println("Number of tries: ", tries)

	compiledVid, err := content.Compile(contentObjs, "compiled-vid.mp4")
	if err != nil {
		fmt.Println("Couldn't compile content.")
		log.Fatal(err)
	}

	err = compiledVid.ApplyOverlay(contentObjs, c.Options)
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
