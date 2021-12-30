package content

import (
	"fmt"
)

// Composite is defined by methods for adding, removing, and retrieving components
type Composite interface {
	Append(c Component) error
	Prepend(c Component) error
}

// Component is defined by an interface for accepting visitors
type Component interface {
	Accept( /*v *Visitor*/ )
}

// Getter is defined by a method for retrieving new components
type Getter interface {
	Get() ([]Component, error)
}

// Sharer is defined by a method to share a component
type Sharer interface {
	Share(Component) error
}

// Query defines fields for retrieving content from external sources.
type Query struct {
	GameID string
	First  int
	Days   int
}

// NewGetter returns new getter interface to the user to suit their platform.
// The credentials string should be a path to a file containing token and client
// info. TODO: simplify this interface accross getters when supporting new
// content sources.
func NewGetter(platform Platform, credentials string, query Query) (Getter, error) {
	switch platform {
	case TWITCH:
		return newTwitchGetter(credentials, query)
	default:
		return nil, fmt.Errorf("No such content getter for platform '%s'", platform)
	}
}
