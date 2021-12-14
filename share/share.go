package share

import (
	"fmt"
	"mcsweeney/config"
	"mcsweeney/content"
)

const YOUTUBE = "youtube"

type ContentSharer interface {
	Share() error
}

// TODO: generic content object instead of path?
func NewContentSharer(c *config.Config, v *content.ContentObj) (ContentSharer, error) {
	switch c.Destination {
	case YOUTUBE:
		return NewYoutubeSharer(c, v)
	default:
		return nil, fmt.Errorf("No such content-sharer '%s'", c.Destination)
	}
}
