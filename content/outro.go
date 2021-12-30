package content

import (
	"fmt"
)

// Outro represents a video.
type Outro struct {
	path string
}

// Accept implements the component interface for Outro.
func (o *Outro) Accept() {
	fmt.Println("Accept not implemented for Outro.")
	return
}
