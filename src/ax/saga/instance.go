package saga

import (
	"github.com/jmalloc/ax/src/ax/ident"
)

// InstanceID uniquely identifies a saga instance.
type InstanceID struct {
	ident.ID
}

// Instance is an instance of a saga.
//
// It encapsulates the application-defined saga data and its meta-data.
type Instance struct {
	// InstanceID is a globally unique identifier for the saga instance.
	InstanceID InstanceID

	// Data is the application-defined data associated with this instance.
	Data Data

	// Revision is the version of the instance that the Data field reflects.
	// A value of zero indicates that the instance has not yet been persisted.
	Revision Revision
}

// Revision is a one-based version of a saga instance.
// An instance with a revision of zero has not yet been persisted.
type Revision uint64
