package ax

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
	MessageTypes() MessageTypeSet

	// HandleMessage invokes application-defined logic that handles m.
	//
	// It panics if m is not one of the types described by MessageTypes().
	HandleMessage(ctx MessageContext, m Message) error
}
