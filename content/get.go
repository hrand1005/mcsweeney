package content

// Getter is defined by a method for retrieving new components
type Getter interface {
	Get() ([]*Clip, error)
}

// NewGetter returns new getter interface to the user to suit their platform.
// The credentials string should be a path to a file containing token and client
// info. TODO: simplify this interface accross getters when supporting new
// content sources.
func NewGetter(platform Platform, credentials string, query Query) (Getter, error) {
	switch platform {
	case TWITCH:
		return NewTwitchGetter(credentials, query)
	default:
		return nil, PlatformNotFound
	}
}
