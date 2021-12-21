package content

import (
	"fmt"
)

type Content struct {
	CreatorName string
	Description string
	Duration    float64
	Keywords    string
	Path        string
	Privacy     string
	Title       string
	Url         string
}

type Query struct {
	GameID string
	First  int
	Days   int
}

type Getter interface {
	Get() ([]*Content, error)
}

type Sharer interface {
	Share(*Content) error
}

const (
	TWITCH  = "twitch"
	YOUTUBE = "youtube"
)

func NewGetter(platform string, credentials string, query Query) (Getter, error) {
	switch platform {
	case TWITCH:
		return NewTwitchGetter(credentials, query)
	default:
		return nil, fmt.Errorf("No such content-getter for platform '%s'", platform)
	}
}

func NewSharer(platform string, credentials string) (Sharer, error) {
	switch platform {
	case YOUTUBE:
		return NewYoutubeSharer(credentials)
	default:
		return nil, fmt.Errorf("No such content-sharer '%s'", platform)
	}
}
