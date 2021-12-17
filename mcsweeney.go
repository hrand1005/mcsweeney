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
*/

//TODO: maybe we don't need the entire packages?
import (
	"fmt"
	"log"
	"mcsweeney/config"
	"mcsweeney/content"
	"mcsweeney/db"
	"mcsweeney/get"
	//"mcsweeney/share"
	"os"
)

const (
	RawVidsDir       = "tmp/raw"
	ProcessedVidsDir = "tmp/processed"
)

func main() {
	// TODO: add command line parsing
	c, err := config.NewConfig(os.Args[1])
	if err != nil {
		fmt.Println("Couldn't load config.")
		log.Fatal(err)
	}

	dbIntf, err := db.NewContentDB(c.Source, "mcsweeney-test.db")
	if err != nil {
		fmt.Println("Couldn't create content-db.")
		log.Fatal(err)
	}

	getIntf, err := get.NewContentGetter(c)
	if err != nil {
		fmt.Println("Couldn't create content-getter.")
		log.Fatal(err)
	}


    dirtyContent, err := getIntf.GetContent()
    if err != nil {
        fmt.Println("Couldn't get new content.")
        log.Fatal(err)
    }

    contentObjs := make([]*content.ContentObj, 0, len(dirtyContent))
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

	compiledVid, err := content.Compile(contentObjs, "compiled-vid.mp4")
	if err != nil {
		fmt.Println("Couldn't compile content.")
		log.Fatal(err)
	}

	err = content.ApplyOverlayWithFade(contentObjs, compiledVid.Path)
	if err != nil {
		fmt.Println("Couldn't apply overlay.")
		log.Fatal(err)
	}

    /*
    shareIntf, err := share.NewContentSharer(c, finalProduct)
    if err != nil {
        fmt.Println("Couldn't create content-sharer.")
        log.Fatal(err)
    }

    err = shareIntf.Share()
    if err != nil {
        fmt.Println("Couldn't share content.")
        log.Fatal(err)
    }
    */

    //fmt.Println("Content shared successfully!")

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
