package saga

import (
	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax/ident"
)

// InstanceID uniquely identifies a saga instance.
type InstanceID struct {
	ident.ID
}

// Instance is an interface for application-defined saga data.
type Instance interface {
	proto.Message

	// InstanceDescription returns a human-readable description of the saga
	// instance.
	//
	// Assume that the description will be used inside log messages or displayed
	// in audit logs.
	//
	// Follow the same conventions as for error messages:
	// https://github.com/golang/go/wiki/CodeReviewComments#error-strings
	InstanceDescription() string
}

// // InstanceEnvelope is a container for an instance and its associated meta-data.
// // TODO: name
// type InstanceEnvelope struct {
// 	InstanceID InstanceID
// 	Revision   uint64
// 	Data       Data
// }
