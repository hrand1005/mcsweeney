package twitch

import (
	"github.com/nicklaw5/helix"
    "fmt"
)

func GetClips(clientID string, token string, gameID string) error {

	client, err := helix.NewClient(&helix.Options{
		ClientID: clientID,
	})
	if err != nil {
        return err	
	}

	client.SetUserAccessToken(token)
	defer client.SetUserAccessToken("")

	// Define query for clips
	clipParams := &helix.ClipsParams{
		GameID: gameID,
	}

	// Execute query for clips
	twitchResp, err := client.GetClips(clipParams)
	if err != nil {
	    return err	
	}

	fmt.Printf("Melee:\n%+v\n", twitchResp)
	return nil
}
