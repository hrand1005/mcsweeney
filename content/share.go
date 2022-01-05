package content

// Payload defines fields for sharing content
type Payload struct {
	Title       string
	Path        string
	Description string
	CategoryID  string
	Keywords    string
	Privacy     string
}

// Sharer is defined by a method to share a payload. It returns an int for a
// HTTP Status Code and an error.
type Sharer interface {
	Share(Payload) (int, error)
}

// NewSharer returns new sharer interface to the user to suit their platform.
// The credentials string should be a path to a file containing token and client
// info.
func NewSharer(platform Platform, credentials string) (Sharer, error) {
	switch platform {
	case YOUTUBE:
		return NewYoutubeSharer(credentials)
	default:
		return nil, ErrPlatformNotFound
	}
}
