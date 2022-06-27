package twitch

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nicklaw5/helix"
)

const (
	AppTokenEnvKey     = "TWITCH_APP_TOKEN"
	ClientIDEnvKey     = "CLIENT_ID"
	ClientSecretEnvKey = "CLIENT_SECRET"
)

// Credentials returns a map of environment variables specific
// to twitch authentication/ app access
func Credentials() map[string]string {
	envMap := map[string]string{
		ClientIDEnvKey:     os.Getenv(ClientIDEnvKey),
		ClientSecretEnvKey: os.Getenv(ClientSecretEnvKey),
		AppTokenEnvKey:     os.Getenv(AppTokenEnvKey),
	}
	return envMap
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

	return err
}
