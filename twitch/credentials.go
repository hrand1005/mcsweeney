package twitch

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nicklaw5/helix"
)

const (
	AppTokenEnvKey = "TWITCH_APP_TOKEN"
	ClientIDEnvKey = "CLIENT_ID"
	ClientSecretEnvKey = "CLIENT_SECRET"
)

// UpdateAppToken generates a new API token and sets the client, environment, and 
// tokenFile 
func UpdateAppToken(client *helix.Client, tokenFile string) error {
	// updates the access token and writes to the token file
	resp, err := client.RequestAppAccessToken(nil)
	if resp.StatusCode != http.StatusOK || err != nil {
		return fmt.Errorf("Encountered error updating app token: %v\n", err)
	}

	f, err := os.Create(tokenFile)
	if err != nil {
		return fmt.Errorf("Scrape: encountered error overwriting token file: %v", err)
	}
	defer f.Close()

	client.SetAppAccessToken(resp.Data.AccessToken)

	os.Setenv(AppTokenEnvKey, resp.Data.AccessToken)
	_, err = f.WriteString(fmt.Sprintf("%s=%s", AppTokenEnvKey, resp.Data.AccessToken))

	return err
}
