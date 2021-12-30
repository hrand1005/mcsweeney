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

// Path implements the component interface for intro.
func (i *Intro) Path() string {
	return i.path
}
