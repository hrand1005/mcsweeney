package content

import (
	"fmt"
)

// Intro represents a video introduction component.
type Intro struct {
	path string
}

// Accept implements the component interface for Intro.
func (i *Intro) Accept() {
	fmt.Println("Accept not implemented for Intro.")
	return
}
