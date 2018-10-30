package ax

// MessageContext provides context about the message being handled.
type MessageContext struct {
	// Envelope is the message envelope containing the message to be handled.
	Envelope Envelope
}

// NewMessageContext returns a message context for the given envelope.
func NewMessageContext(
	env Envelope,
) MessageContext {
	return MessageContext{
		Envelope: env,
	}
}
