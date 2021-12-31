package content

// Composite is defined by methods for adding, removing, and retrieving components
type Composite interface {
	Append(Component) error
	Prepend(Component) error
}

// Component is defined by an interface for accepting visitors
type Component interface {
	Accept(Visitor)
}
