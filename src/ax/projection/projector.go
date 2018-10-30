package projection

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Projector is an interface for a specialized form of application-defined
// message handler which produces a "projection" of state from the messages it
// receives.
type Projector interface {
	// PersistenceKey returns a unique name for the projector.
	//
	// The persistence key is used to relate persisted data with the projector
	// implementation that owns it. Persistence keys should not be changed once
	// a projection has been started.
	PersistenceKey() string

	// MessageTypes returns the set of messages that the projector intends
	// to handle.
	//
	// The return value should be constant as it may be cached.
	MessageTypes() ax.MessageTypeSet

	// ApplyMessage invokes application-defined logic that updates the
	// application state to reflect the occurrence of a message.
	//
	// It may panic if env.Message is not one of the types described by
	// MessageTypes().
	ApplyMessage(ctx context.Context, mctx ax.MessageContext) error
}
