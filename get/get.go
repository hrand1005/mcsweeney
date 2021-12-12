package get

import (
	"fmt"
	"mcsweeney/config"
	"mcsweeney/content"
	"mcsweeney/db"
)

// TODO: remove this duplicate
const TWITCH = "twitch"

// TODO: generic content object or interface
type ContentGetter interface {
	GetContent(db.ContentDB) ([]*content.ContentObj, error)
}

func NewContentGetter(c config.Config) (ContentGetter, error) {
	switch c.Source {
	case TWITCH:
		return NewTwitchGetter(c)
	default:
		return nil, fmt.Errorf("No such content-getter '%s'", c.Source)
	}
}
