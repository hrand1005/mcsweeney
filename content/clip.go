package content

import (
	"fmt"
)

type Platform string

const (
	TWITCH  Platform = "twitch"
	YOUTUBE Platform = "youtube"
)

// Clip represents a video clip retrieved from some external source.
type Clip struct {
	Author      string
	Broadcaster string
	Duration    float64
	Language    string
	Path        string
	Platform    Platform
	Title       string
}

// Accept implements the component interface for Clip.
func (c *Clip) Accept() {
	fmt.Println("Accept not implemented for Clip.")
	return
}
