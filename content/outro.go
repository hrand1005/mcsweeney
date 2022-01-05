package content

// Outro represents a video outro component.
type Outro struct {
	Description string
	Duration    float64
	Path        string
}

// Accept implements the component interface for Outro.
func (o *Outro) Accept(v Visitor) error {
	return v.VisitOutro(o)
}
