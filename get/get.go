package get

import (
	"fmt"
	"github.com/nicklaw5/helix"
	"mcsweeney/config"
	"mcsweeney/db"
)

// TODO: remove this duplicate
const TWITCH = "twitch"

// TODO: generic content object or interface
type ContentGetter interface {
	GetContent() ([]helix.Clip, error)
}

func NewContentGetter(c config.Config, db db.ContentDB) (ContentGetter, error) {
	switch c.Source {
	case TWITCH:
		return NewTwitchGetter(c, db)
	default:
		return nil, fmt.Errorf("No such content-getter '%s'", c.Source)
	}
}
