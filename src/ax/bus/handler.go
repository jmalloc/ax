package bus

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// MessageHandler is an interface for application-defined message handlers.
//
// Message handlers are typically the last stage in the inbound message
// pipeline. Each message handler declares its interest in a specific set
// of message types and is notified when any matching message arrives.
type MessageHandler interface {
	// MessageTypes returns the set of messages that the handler intends
	// to handle.
	//
	// The return value should be constant as it may be cached by various
	// independent stages in the message pipeline.
	MessageTypes() ax.MessageTypeSet

	// HandleMessage invokes application-defined logic that handles a
	// message.
	//
	// It may panic if env.Message is not one of the types described by
	// MessageTypes().
	HandleMessage(ctx context.Context, s ax.Sender, env ax.Envelope) error
}
