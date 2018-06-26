package saga

import (
	"github.com/jmalloc/ax/src/ax/ident"
)

// InstanceID uniquely identifies a saga instance.
type InstanceID struct {
	ident.ID
}

// GenerateInstanceID generates a new unique identifier for a saga instance.
func GenerateInstanceID() InstanceID {
	var id InstanceID
	id.GenerateUUID()
	return id
}

// ParseInstanceID parses s into a saga instance ID and returns it. It returns
// an error if s is empty.
func ParseInstanceID(s string) (InstanceID, error) {
	var id InstanceID
	err := id.Parse(s)
	return id, err
}

// MustParseInstanceID parses s into a saga instance ID and returns it. It
// panics if s is empty.
func MustParseInstanceID(s string) InstanceID {
	var id InstanceID
	id.MustParse(s)
	return id
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
