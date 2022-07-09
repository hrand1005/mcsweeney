package main

import (
	"fmt"
	"strings"

	"github.com/nicklaw5/helix"
)

type DescriptionBuilder struct {
	descriptions []string
	cursor       float64
}

func NewDescriptionBuilder() *DescriptionBuilder {
	return &DescriptionBuilder{}
}

// AppendCustom appends the custom string to the description. Does not
// affect the timestamps of any clip descriptions that may have been added.
func (d *DescriptionBuilder) AppendCustom(custom string) {
	d.descriptions = append(d.descriptions, custom)
}

// PrependCustom prepends the custom string to the description. Does not
// affect the timestamps of any clip descriptions that may have been added.
func (d *DescriptionBuilder) PrependCustom(custom string) {
	d.descriptions = append([]string{custom}, d.descriptions...)
}

// AppendClip appends a clip description. It's important that clips are appended
// in the order in which they may appear in a video so that the timestamps line up.
// Given this expectation, PrependClip is not yet supported.
func (d *DescriptionBuilder) AppendClip(c helix.Clip) {
	ts := timestamp(d.cursor)
	clipDesc := fmt.Sprintf("%s %q\nStreamed by %s at %s\nClipped by %s\n\n", ts, c.Title, c.BroadcasterName, GetTwitchClipChannel(c), c.CreatorName)
	d.descriptions = append(d.descriptions, clipDesc)
	d.cursor += c.Duration
}

func (d *DescriptionBuilder) Result() string {
	return strings.Join(d.descriptions, "")
}

// timestamp generates a timestamp string from the given float
func timestamp(t float64) string {
	minutes := int(t) / 60
	seconds := int(t) % 60
	if seconds < 10 {
		return fmt.Sprintf("[%v:0%v]", minutes, seconds)
	}
	return fmt.Sprintf("[%v:%v]", minutes, seconds)
}
