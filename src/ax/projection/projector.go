package projection

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Projector is an interface for a specialized form of application-defined
// message handler which produces a "projection" of state from the messages it
// receives.
type Projector interface {
	// ProjectorName returns a unique name for the projector.
	//
	// The projector's name is used to correlate persisted data with this
	// instance, so it should not be changed.
	ProjectorName() string

	// MessageTypes returns the set of messages that the projector intends
	// to handle.
	//
	// The return value should be constant as it may be cached.
	MessageTypes() ax.MessageTypeSet

	// ApplyMessage invokes application-defined logic that updates the
	// application state to reflect the delivery of a message.
	//
	// It may panic if env.Message is not one of the types described by
	// MessageTypes().
	ApplyMessage(ctx context.Context, env ax.Envelope) error
}
