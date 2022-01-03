package content

type Payload struct {
	Title       string
	Path        string
	Description string
	Keywords    string
	Privacy     string
}

// Sharer is defined by a method to share a component
type Sharer interface {
	Share(Payload) error
}

// NewSharer returns new sharer interface to the user to suit their platform.
// The credentials string should be a path to a file containing token and client
// info. TODO: simplify this interface accross sharers when supporting new
// content sources.
func NewSharer(platform Platform, credentials string) (Sharer, error) {
	switch platform {
	case YOUTUBE:
		return NewYoutubeSharer(credentials)
	default:
		return nil, PlatformNotFound
	}
}
