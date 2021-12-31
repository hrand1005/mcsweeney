package content

import (
	"fmt"
)

// DescriptionGenerator generates descriptions for the video components it
// visits.
type DescriptionGenerator struct {
	cursor      float64
	description string
}

// String returns a formatted string of generated description. The description
// is the aggregate of all visited elements, also reflecting visit order.
func (d *DescriptionGenerator) String() string {
	return d.description
}

// VisitClip implements the visitor interface for DescriptionGenerator. Appends
// a formatted timestamp, title, broadcaster, and author of the clip.
func (d *DescriptionGenerator) VisitClip(c *Clip) {
	// if the clip has no duration, do nothing with it
	// TODO: raise error flag in DescriptionGenerator
	if c.Duration == 0.0 {
		return
	}
	// generate simple timestamp up to 59:59
	minutes := int(d.cursor) / 60
	seconds := int(d.cursor) % 60
	var timestamp string
	if seconds < 10 {
		timestamp = fmt.Sprintf("[%v:0%v]", minutes, seconds)
	} else {
		timestamp = fmt.Sprintf("[%v:%v]", minutes, seconds)
	}

	d.description += fmt.Sprintf("\n\n%s '%s'\nStreamed by %s at %s\nClipped by %s\n", timestamp, c.Title, c.Broadcaster, c.Channel(), c.Author)
	d.cursor += c.Duration
	return
}

// VisitIntro implements the visitor interface for DescriptionGenerator. Appends
// a faithful duplicate of the intro's description field.
func (d *DescriptionGenerator) VisitIntro(i *Intro) {
	d.description += i.Description
	d.cursor += i.Duration
	return
}

// VisitOutro implements the visitor interface for DescriptionGenerator. Appends
// a faithful duplicate of the outro's description field.
func (d *DescriptionGenerator) VisitOutro(o *Outro) {
	d.description += o.Description
	d.cursor += o.Duration
	return
}
