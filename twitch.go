package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nicklaw5/helix"
	"golang.org/x/oauth2"
)

func GetTwitchClipChannel(c helix.Clip) string {
	return "https://twitch.tv/" + c.BroadcasterName
}

const (
	TwitchTokenFileKey    = "TWITCH_TOKEN_FILE"
	TwitchClientIDKey     = "TWITCH_CLIENT_ID"
	TwitchClientSecretKey = "TWITCH_CLIENT_SECRET"
	TwitchTokenRawKey     = "TWITCH_TOKEN_RAW"
)

// setTwitchToken attempst to set a valid access token on the given client.
// First looks for an existnig access token in the twitch token file env.
// If a token cannot be retrieved, requests a new token using the given client.
func setTwitchToken(client *helix.Client) error {

	// token, err := readTokenFromFile(os.Getenv(TwitchTokenFileKey))
	// log.Printf("Error reading token from file: %v", err)
	// if err == nil {
	// 	client.SetAppAccessToken(token.AccessToken)
	// 	return nil
	// }

	// token, err = getNewTwitchToken(client)
	// if err != nil {
	// 	return fmt.Errorf("failed to get new twitch token: %v", err)
	// }
	accessToken := os.Getenv(TwitchTokenRawKey)
	client.SetAppAccessToken(accessToken)

	return nil
	// return writeTokenToFile(os.Getenv(TwitchTokenFileKey), token)
}

func getNewTwitchToken(c *helix.Client) (*oauth2.Token, error) {
	fmt.Println("Requesting new twitch access token.")
	resp, err := c.RequestAppAccessToken(nil)
	if err != nil {
		return nil, fmt.Errorf("encountered error requesting new twitch token: %v\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got resp with status code %v and message %s when requesting twitch token", resp.StatusCode, resp.ErrorMessage)
	}

	return &oauth2.Token{
		AccessToken:  resp.Data.AccessToken,
		RefreshToken: resp.Data.RefreshToken,
	}, nil
}

// NewTwitchClient as defined in this package creates a Twitch client using environment variables
func NewTwitchClient() (*helix.Client, error) {
	cOpts := &helix.Options{
		ClientID:     os.Getenv(TwitchClientIDKey),
		ClientSecret: os.Getenv(TwitchClientSecretKey),
	}

	client, err := helix.NewClient(cOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create new twitch client: %v", err)
	}

	if err = setTwitchToken(client); err != nil {
		return nil, fmt.Errorf("failed to set token for new client: %v", err)
	}

	return client, nil
}
