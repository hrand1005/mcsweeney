package twitch

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nicklaw5/helix"
)

func GetClipChannel(c helix.Clip) string {
	return "https://twitch.tv/" + c.BroadcasterName
}

const (
	AppTokenEnvKey     = "TWITCH_APP_TOKEN"
	ClientIDEnvKey     = "CLIENT_ID"
	ClientSecretEnvKey = "CLIENT_SECRET"
)

// Credentials returns a map of environment variables specific
// to twitch authentication/ app access
func Credentials() map[string]string {
	return map[string]string{
		ClientIDEnvKey:     os.Getenv(ClientIDEnvKey),
		ClientSecretEnvKey: os.Getenv(ClientSecretEnvKey),
		AppTokenEnvKey:     os.Getenv(AppTokenEnvKey),
	}
}

// UpdateAppToken generates a new API token and sets the client, environment
func UpdateAppToken(client *helix.Client) error {
	resp, err := client.RequestAppAccessToken(nil)
	if resp.StatusCode != http.StatusOK || err != nil {
		return fmt.Errorf("Encountered error updating app token: %v\n", err)
	}

	// set new App token in the client and environment
	client.SetAppAccessToken(resp.Data.AccessToken)
	os.Setenv(AppTokenEnvKey, resp.Data.AccessToken)

	return nil
}

// NewClient as defined in this package creates a Twitch client using environment variables
func NewClient() (*helix.Client, error) {
	cOpts := &helix.Options{
		ClientID:     os.Getenv(ClientIDEnvKey),
		ClientSecret: os.Getenv(ClientSecretEnvKey),
	}

	client, err := helix.NewClient(cOpts)
	if err != nil {
		return nil, fmt.Errorf("NewClient: failed to create new twitch client: %v", err)
	}

	client.SetAppAccessToken(os.Getenv(AppTokenEnvKey))

	return client, nil
}
