package ax

// MessageContext contains information about the context in which a message is
// handled.
type MessageContext struct {
	Envelope Envelope
	Sender   Sender
}
