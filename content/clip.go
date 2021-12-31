package content

// Clip represents a video clip retrieved from some external source.
type Clip struct {
	Author      string
	Broadcaster string
	Duration    float64
	Language    string
	Path        string
	Platform    Platform
	Title       string
}

// Channel returns the source channel (broadcaster homepage) of the clip element
func (c *Clip) Channel() string {
	switch c.Platform {
	case TWITCH:
		return "https://twitch.tv/" + c.Broadcaster
	default:
		return ""
	}
}

// Accept implements the component interface for Clip.
func (c *Clip) Accept(v Visitor) {
	v.VisitClip(c)
	return
}
