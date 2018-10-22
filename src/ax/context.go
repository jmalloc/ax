package ax

// MessageContext provides context about the message being handled.
type MessageContext struct {
	// Envelope is the message envelope containing the message to be handled.
	Envelope Envelope
}
