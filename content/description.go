package content

import (
	"fmt"
	"strings"
)

// Describer generates descriptions for the video components it
// visits.
type Describer struct {
	cursor      float64
	description []string
}

// String returns a formatted string of generated description. The description
// is the aggregate of all visited elements, also reflecting visit order.
func (d *Describer) String() string {
	return strings.Join(d.description, "")
}

// VisitClip implements the visitor interface for Describer. Appends
// a formatted timestamp, title, broadcaster, and author of the clip.
func (d *Describer) VisitClip(c *Clip) error {
	// if the clip has no duration, do nothing with it
	if c.Duration == 0.0 {
		return ErrNoDuration
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
	// concatenate description, increment cursor
	d.description = append(d.description, fmt.Sprintf("\n\n%s '%s'\nStreamed by %s at %s\nClipped by %s\n", timestamp, c.Title, c.Broadcaster, c.Channel(), c.Author))
	d.cursor += c.Duration

	return nil
}

// VisitIntro implements the visitor interface for Describer. Appends
// a faithful duplicate of the intro's description field.
func (d *Describer) VisitIntro(i *Intro) error {
	d.description = append(d.description, i.Description)
	d.cursor += i.Duration
	return nil
}

// VisitOutro implements the visitor interface for Describer. Appends
// a faithful duplicate of the outro's description field.
func (d *Describer) VisitOutro(o *Outro) error {
	d.description = append(d.description, o.Description)
	d.cursor += o.Duration
	return nil
}
