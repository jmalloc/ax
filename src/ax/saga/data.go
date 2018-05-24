package saga

import (
	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax"
)

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

// EventedData is a specialization of Data for sagas that use events to update
// their state. Event-sourced sagas always use EventedData.
type EventedData interface {
	Data

	// ApplyEvent updates the data to reflect the fact that an event has
	// occurred.
	//
	// It may panic if env.Message does not implement ax.Event.
	ApplyEvent(env ax.Envelope)
}
