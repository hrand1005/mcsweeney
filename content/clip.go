package content

import (
	"fmt"
)

// Clip represents a video clip retrieved from some external source.
type Clip struct {
	title       string
	broadcaster string
	author      string
	path        string
}

// Accept implements the component interface for Clip.
func (c *Clip) Accept() {
	fmt.Println("Accept not implemented for Clip.")
	return
}
