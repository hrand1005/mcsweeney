package main

/* TODO:
- Consistent error handling
- Testing
- Architecture
    - Think about params, pointers, interfaces, etc.
- Goroutines
- Decide on a standard for where to use pointers vs struct values
*/

//TODO: maybe we don't need the entire packages?
import (
	"fmt"
	"log"
	"mcsweeney/config"
	"mcsweeney/db"
	"mcsweeney/edit"
	"mcsweeney/get"
	"mcsweeney/share"
	//"sync"
)

const (
	RawVidsDir       = "tmp/raw"
	ProcessedVidsDir = "tmp/processed"
)

func main() {
	// Get config from yaml file
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

	getIntf, err := get.NewContentGetter(*c, dbIntf)
	if err != nil {
		fmt.Println("Couldn't create content-getter.")
		log.Fatal(err)
	}

	content, err := getIntf.GetContent()
	if err != nil {
		fmt.Println("Couldn't get new content.")
		log.Fatal(err)
	}

	//s.EditContent()
	err = edit.ApplyOverlay(content)
	if err != nil {
		fmt.Println("Couldn't edit some clips")
		log.Fatal(err)
	}

	//s.CompileContent()
	err = edit.Compile()
	if err != nil {
		fmt.Println("Couldn't compile clips")
		log.Fatal(err)
	}

	shareIntf, err := share.NewContentSharer(*c, "compiled-vid.mp4")
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
