package saga

import "github.com/jmalloc/ax/src/ax/ident"

// InstanceID uniquely identifies a saga instance.
type InstanceID struct {
	ident.ID
}

// Instance is an interface for application-defined saga data.
type Instance interface {
	// InstanceID returns a unique identifier for the saga instance.
	InstanceID() InstanceID
}
