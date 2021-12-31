package content

// Intro represents a video introduction component.
type Intro struct {
	Description string
	Duration    float64
	Path        string
}

// Accept implements the component interface for Intro.
func (i *Intro) Accept(v Visitor) {
	v.VisitIntro(i)
	return
}
