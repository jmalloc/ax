package ax

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax/ident"
	"github.com/jmalloc/ax/src/ax/marshaling"
)

// MessageID uniquely identifies a message.
type MessageID struct {
	ident.ID
}

// Message is a unit of communication.
type Message interface {
	proto.Message

	// MessageDescription returns a human-readable description of the message.
	//
	// Assume that the description will be used inside log messages or displayed
	// in audit logs.
	//
	// Follow the same conventions as for error messages:
	// https://github.com/golang/go/wiki/CodeReviewComments#error-strings
	MessageDescription() string
}

// Command is a message that requests some action take place.
//
// Commands are always sent to a single handler within a single end-point.
type Command interface {
	Message

	// IsCommand() is a "marker method" used to indicate that a message is
	// intended to be used as a command.
	IsCommand()
}

// Event is a message that indicates some action has already taken place.
//
// Events are published by one endpoint and (potentially) consumed by many.
type Event interface {
	Message

	// IsEvent() is a "marker method" used to indicate that a message is
	// intended to be used as an event.
	IsEvent()
}

var (
	commandType = reflect.TypeOf((*Command)(nil)).Elem()
	eventType   = reflect.TypeOf((*Event)(nil)).Elem()
)

// MarshalMessage marshals m to a binary representation.
func MarshalMessage(m Message) (contentType string, data []byte, err error) {
	return marshaling.MarshalProtobuf(m)
}

// UnmarshalMessage unmarshals an Ax message from some serialized
// representation. ct is the MIME content-type for the binary data.
func UnmarshalMessage(ct string, data []byte) (Message, error) {
	v, err := marshaling.Unmarshal(ct, data)
	if err != nil {
		return nil, err
	}

	if m, ok := v.(Message); ok {
		return m, nil
	}

	return nil, fmt.Errorf(
		"can not unmarshal '%s', content-type is not an Ax message",
		ct,
	)
}
