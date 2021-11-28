package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"mcsweeney/twitch"
)

type context struct {
	ClientID string `yaml:"clientID"`
	Token    string `yaml:"token"`
    GameID string `yaml:"gameID"` 
    First int `yaml:"first"`
}

func main() {
	// Get context from yaml file
	c := context{}

	err := loadContext("example.yaml", &c)
	if err != nil {
		fmt.Println("Couldn't load context.")
		log.Fatal(err)
	}

	// Remember, this is a strategy, so it will be more like s.GetContent()
	err = twitch.GetClips(c.ClientID, c.Token, c.GameID, c.First)
	if err != nil {
		fmt.Println("Couldn't get content.")
		log.Fatal(err)
	}

	//s.EditContent()
	//s.CompileContent()
	//s.ShareContent()

	return
}

// It may be appropriate to get more information than just a token
func loadContext(path string, c *context) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read %s: %w", path, err)
	}

	err = yaml.Unmarshal(raw, c)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal: %v", err)
	}

	return nil
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
