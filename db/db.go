package db

import (
	"fmt"
	"mcsweeney/content"
)

const TWITCH = "twitch"

type ContentDB interface {
	Insert(*content.ContentObj) error
	Exists(*content.ContentObj) (bool, error)
}

func NewContentDB(source string, name string) (ContentDB, error) {
	switch source {
	case TWITCH:
		return NewTwitchDB(name)
	default:
		return nil, fmt.Errorf("DB %s not found", source)
	}
}
