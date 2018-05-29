package saga

import (
	"github.com/golang/protobuf/proto"
)

// Data is an interface for application-defined data associated with a saga
// instance.
type Data interface {
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
