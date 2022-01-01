package content

import (

)

// OverlayGenerator generates video overlays for the video components it visits.
type OverlayGenerator struct {
    cursor float64
    overlay string
}

// String returns a string of the generated overlay. The overlay is the 
// aggregate of all visited elements, also reflecting visit order.
func (o *OverlayGenerator) String() string {
    return o.overlay
}

// VisitClip implements the visitor interface for OverlayGenerator. 
func (o *OverlayGenerator) VisitClip(c *Clip) {
    return
}

// VisitIntro implements the visitor interface for OverlayGenerator. 
// Increments internal cursor, but does not alter the overlay string.
func (o *OverlayGenerator) VisitIntro(i *Intro) {
	o.cursor += i.Duration
	return
}

// VisitOutro implements the visitor interface for OverlayGenerator. 
// Increments internal cursor, but does not alter the overlay string.
func (o *OverlayGenerator) VisitOutro(u *Outro) {
	o.cursor += u.Duration
	return
}
