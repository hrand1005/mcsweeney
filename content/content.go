package content

import (
	"errors"
)

type Platform string

const (
	TWITCH  Platform = "twitch"
	YOUTUBE Platform = "youtube"
)

// PlatformNotFound is returned when attempting an operation on an invalid
// platform argument
var PlatformNotFound error = errors.New("Platform not found.")

// Query defines fields for retrieving content from external sources.
type Query struct {
	GameID string
	First  int
	Days   int
}
