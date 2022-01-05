package content

// Video represents a composite object with a slice of component interfaces
type Video struct {
	components []Component
}

// Append implements the interface for Composite.Append(). It requires a component
// interface, which is added to the end of the Composite video.
func (v *Video) Append(c Component) error {
	v.components = append(v.components, c)
	return nil
}

// Prepend implements the interface for Composite.Prepend(). It requires a
// component interface, which is added to the beginning of the Composite video.
func (v *Video) Prepend(c Component) error {
	v.components = append([]Component{c}, v.components...)
	return nil
}

// Accept implements the component interface for Video. It calls accept on its
// child components.
func (v *Video) Accept(t Visitor) error {
	for _, c := range v.components {
		err := c.Accept(t)
		if err != nil {
			return err
		}
	}
	return nil
}
