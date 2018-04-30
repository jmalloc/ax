package bus

import "github.com/jmalloc/ax/src/ax"

// MessageHandler is an interface for application-defined message handlers.
//
// Message handlers are typically the last stage in the inbound message
// pipeline. Each message handler registers interest in a specific set of
// message types and is notified when any matching message arrives.
type MessageHandler interface {
	// MessageTypes returns the set of messages that should be passed to
	// HandleMessage upon delivery.
	//
	// The return value should be constant as it may be cached by various
	// independent stages in the message pipeline.
	MessageTypes() ax.MessageTypeSet

	// HandleMessage invokes application-defined logic that handles m.
	//
	// It panics if m is not one of the types described by MessageTypes().
	HandleMessage(ctx ax.MessageContext, m ax.Message) error
}

// MessageHandlerFunc returns a MessageHandler that calls fn when a message
// in the type set mt is delivered.
func MessageHandlerFunc(
	mt ax.MessageTypeSet,
	fn func(ax.MessageContext, ax.Message) error,
) MessageHandler {
	return messageHandlerFunc{mt, fn}
}

// messageHandlerFunc is an implementation of MessageHandler that dispatches to
// a function.
type messageHandlerFunc struct {
	mt ax.MessageTypeSet
	fn func(ax.MessageContext, ax.Message) error
}

func (h messageHandlerFunc) MessageTypes() ax.MessageTypeSet {
	return h.mt
}

func (h messageHandlerFunc) HandleMessage(ctx ax.MessageContext, m ax.Message) error {
	return h.fn(ctx, m)
}
