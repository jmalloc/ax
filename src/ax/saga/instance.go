package saga

import (
	"github.com/golang/protobuf/proto"
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
	InstanceID InstanceID
	Revision   uint64
	Data       Data
}

// Data is an interface for application-defined data associated with a saga
// instance.
type Data interface {
	proto.Message

	// SagaDescription returns a human-readable description of the saga
	// instance.
	//
	// Assume that the description will be used inside log messages or displayed
	// in audit logs.
	//
	// Follow the same conventions as for error messages:
	// https://github.com/golang/go/wiki/CodeReviewComments#error-strings
	SagaDescription() string
}
