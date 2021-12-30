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
	author      string
	broadcaster string
	duration    float64
	language    string
	path        string
	platform    Platform
	title       string
}

// Accept implements the component interface for Clip.
func (c *Clip) Accept() {
	fmt.Println("Accept not implemented for Clip.")
	return
}

// Path implements the component interface for Clip.
func (c *Clip) Path() string {
	return c.path
}
