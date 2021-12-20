package content

import (
	"fmt"
	"mcsweeney/config"
)

type Content struct {
	CreatorName string
	Description string
	Duration    float64
	Path        string
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
	Share() error
}

const TWITCH = "twitch"
const YOUTUBE = "youtube"

func NewGetter(platform string, credentials string, query Query) (Getter, error) {
	switch platform {
	case TWITCH:
		return NewTwitchGetter(credentials, query)
	default:
		return nil, fmt.Errorf("No such content-getter for platform '%s'", platform)
	}
}

func NewSharer(c *config.Config, v *Content) (Sharer, error) {
	switch c.Destination.Platform {
	case YOUTUBE:
		return NewYoutubeSharer(c, v)
	default:
		return nil, fmt.Errorf("No such content-sharer '%s'", c.Destination)
	}
}
