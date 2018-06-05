package projection

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// MessageHandler exposes as Projector as a routing.MessageHandler.
type MessageHandler struct {
	Projector Projector
}

// MessageTypes returns the set of messages that the handler intends
// to handle.
//
// The return value should be constant as it may be cached by various
// independent stages in the message pipeline.
func (a *MessageHandler) MessageTypes() ax.MessageTypeSet {
	return a.Projector.MessageTypes()
}

// HandleMessage invokes application-defined logic that handles a
// message.
//
// It may panic if env.Message is not one of the types described by
// MessageTypes().
func (a *MessageHandler) HandleMessage(
	ctx context.Context,
	_ ax.Sender,
	env ax.Envelope,
) error {
	return a.Projector.ApplyMessage(ctx, env)
}
