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
	"mcsweeney/share"
	//"sync"
)

const (
	RawVidsDir       = "tmp/raw"
	ProcessedVidsDir = "tmp/processed"
)

func main() {
	c, err := config.NewConfig("config/example.yaml")
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

	contentObjs, err := getIntf.GetContent(dbIntf)
	if err != nil {
		fmt.Println("Couldn't get new content.")
		log.Fatal(err)
	}

	for _, v := range contentObjs {
		// TODO: go func() for all this
		err = v.ApplyOverlay(RawVidsDir)
		if err != nil {
			fmt.Println("Couldn't download content.")
			log.Fatal(err)
		}
	}

	finalProduct, err := content.Compile(contentObjs)
	if err != nil {
		fmt.Println("Couldn't compile content.")
		log.Fatal(err)
	}

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

	fmt.Println("Content shared successfully!")

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
