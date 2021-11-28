package main

import (
    "mcsweeney/twitch"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type context struct {
	ClientID string `yaml:"clientID"`
	Token    string `yaml:"token"`
}

func main() {
	// Get context from yaml file
	c := context{}

	err := loadContext("example.yaml", &c)
	if err != nil {
		fmt.Printf("Couldn't load context.")
		log.Fatal(err)
	}

	fmt.Printf("ClientID: %s", c.ClientID)

    twitch.GetClips(c.ClientID, c.Token, "16282")
    /*
	client, err := helix.NewClient(&helix.Options{
		ClientID: c.ClientID,
	})
	if err != nil {
		log.Fatal(err)
	}

	client.SetUserAccessToken(c.Token)
	defer client.SetUserAccessToken("")

	// Define query for clips
	clipParams := &helix.ClipsParams{
		GameID: "16282",
	}

	// Execute query for clips
	twitchResp, err := client.GetClips(clipParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Melee:\n%+v\n", twitchResp)
    */
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
    
