package main

import (
	"fmt"

	"github.com/hrand1005/mcsweeney/twitch"
	"github.com/nicklaw5/helix"
)

// Generates a string description from an ordered list of twitch clips
// The offset param shifts the starting timestamp by offset seconds
func DescriptionFromTwitchClips(clips []helix.Clip, offset float64) string {
	var cursor float64 = offset
	description := ""
	for _, v := range clips {
		ts := timestamp(cursor)
		description += fmt.Sprintf("%s %q\nStreamed by %s at %s\nClipped by %s\n\n", ts, v.Title, v.BroadcasterName, twitch.GetClipChannel(v), v.CreatorName)
		cursor += v.Duration
	}

	return description
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
