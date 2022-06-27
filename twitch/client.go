package twitch 

import (
  "fmt"
  "os"

	"github.com/nicklaw5/helix"
)

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
