package projection

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Projector is an interface for application-defined message handlers that are
// designed to construct read-models.
//
// Unlike a routing.MessageHandler, they do not accept an ax.Sender argument.
type Projector interface {
	// MessageTypes returns the set of messages that the projector intends
	// to handle.
	//
	// The return value should be constant as it may be cached.
	MessageTypes() ax.MessageTypeSet

	// HandleMessage invokes application-defined logic that handles a
	// message.
	//
	// It may panic if env.Message is not one of the types described by
	// MessageTypes().
	HandleMessage(ctx context.Context, env ax.Envelope) error
}
