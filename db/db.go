package db

import (
	"fmt"
)

const TWITCH = "twitch"

type ContentDB interface {
	Create() error
	Insert(string) error
	Exists(string) (bool, error)
}

func NewContentDB(source string) (ContentDB, error) {
	switch source {
	case TWITCH:
		return &TwitchDB{}, nil
	default:
		return nil, fmt.Errorf("DB %s not found", source)
	}
}
