package content

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
