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

func NewContentGetter(s config.Source) (ContentGetter, error) {
	switch s.Platform {
	case TWITCH:
		return NewTwitchGetter(s.Credentials, s.Query)
	default:
		return nil, fmt.Errorf("No such content-getter for platform '%s'", s.Platform)
	}
}
