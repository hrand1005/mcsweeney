package db

import (
	"fmt"
)

const TWITCH = "twitch"

type ContentDB interface {
	Insert(string) error
	Exists(string) (bool, error)
}

func NewContentDB(source string, name string) (ContentDB, error) {
	switch source {
	case TWITCH:
		return NewTwitchDB(name)
	default:
		return nil, fmt.Errorf("DB %s not found", source)
	}
}
