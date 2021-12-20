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

type Getter interface {
	Get() ([]*Content, error)
}

type Sharer interface {
	Share() error
}

const TWITCH = "twitch"
const YOUTUBE = "youtube"

func NewGetter(s config.Source) (Getter, error) {
	switch s.Platform {
	case TWITCH:
		return NewTwitchGetter(s.Credentials, s.Query)
	default:
		return nil, fmt.Errorf("No such content-getter for platform '%s'", s.Platform)
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
