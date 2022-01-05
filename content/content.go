package content

import (
	"errors"
)

type Platform string

const (
	TWITCH  Platform = "twitch"
	YOUTUBE Platform = "youtube"
)

// ErrEmptyPath is returned when attempting an operation that requires a path
var ErrEmptyPath error = errors.New("Path not defined.")

// ErrNoDuration is returned when an operation cannot function on with an
// element with no duration.
var ErrNoDuration = errors.New("This element has no duration.")

// PlatformNotFound is returned when attempting an operation on an invalid
// platform argument
var ErrPlatformNotFound = errors.New("Platform not found.")

// Query defines fields for retrieving content from external sources.
type Query struct {
	GameID string
	First  int
	Days   int
}
