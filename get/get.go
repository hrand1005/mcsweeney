package get

import (
	"fmt"
	"mcsweeney/config"
	"mcsweeney/content"
)

// TODO: remove this duplicate
const TWITCH = "twitch"

type ContentGetter interface {
	GetContent() ([]*content.ContentObj, error)
}

func NewContentGetter(c *config.Config) (ContentGetter, error) {
	switch c.Source {
	case TWITCH:
		return NewTwitchGetter(c.ClientID, c.Query, c.Token)
	default:
		return nil, fmt.Errorf("No such content-getter '%s'", c.Source)
	}
}
