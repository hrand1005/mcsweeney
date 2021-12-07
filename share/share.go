package share

import (
	"fmt"
	"mcsweeney/config"
)

const YOUTUBE = "youtube"

type ContentSharer interface {
	Share() error
}

// TODO: generic content object instead of path?
func NewContentSharer(c config.Config, contentPath string) (ContentSharer, error) {
	switch c.Destination {
	case YOUTUBE:
		return NewYoutubeSharer(c, contentPath)
	default:
		return nil, fmt.Errorf("No such content-sharer '%s'", c.Destination)
	}
}
