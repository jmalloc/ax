package ax

import "time"

// Envelope is a container for a message and its associated meta-data.
type Envelope struct {
	// MessageID is a globally unique identifier for a single message.
	//
	// Among other things, the message ID is often used during message
	// de-duplicatation in order to provide exactly-one handling semantics.
	MessageID MessageID

	// CausationID is the ID of the message that directly caused this message to
	// occur.
	//
	// Messages can be thought of as occurring within a tree of messages. The
	// CausationID identifies the direct parent message within that tree.
	//
	// When a message is injected into the messaging system via a MessageBus,
	// the CausationID is set to the MessageID, that is, the message is its own
	// cause.
	//
	// When a message is sent via a MessageContext, the CausationID is
	// automatically set to the MessageID of the message being handled in that
	// context.
	CausationID MessageID

	// CorrelationID is the ID of the message that (perhaps indirectly) caused
	// this message to occur.
	//
	// Messages can be thought of as occurring within a tree of messages. The
	// CorrelationID identifies the message at the root of that tree.
	//
	// When a message is injected into the message system via a MessageBus,
	// the CorrelationID is set to the MessageID, that is, the message is at the
	// root of the tree.
	//
	// When a message is sent via a MessageContext, the CorrelationID is
	// automatically set to the CorrelationID of the message being handled in
	// that context.
	CorrelationID MessageID

	// Time is the time at which the message was created.
	//
	// Time is populated via the regular Go system clock, and as such there are
	// almost no guarantees about the accuracy of the time. It must not be
	// assumed that messages will arrive or even be delivered in any
	// chronological order.
	Time time.Time

	// Message is the application-defined message encapsulated by the envelope.
	Message Message
}

// NewEnvelope creates a new message envelope containing m.
//
// It generates a UUID-based message ID and configures the envelope such that m
// is at the root of a new tree of messages.
func NewEnvelope(m Message) Envelope {
	id := GenerateMessageID()

	env := Envelope{
		MessageID:     id,
		CausationID:   id,
		CorrelationID: id,
		Time:          time.Now(),
		Message:       m,
	}

	return env
}

// NewChild returns a new message envelope containing m.
//
// It generates a UUID-based message ID and configures the envelope such that
// m is a child of e.Message within an existing tree of messages.
func (e Envelope) NewChild(m Message) Envelope {
	env := Envelope{
		MessageID:     GenerateMessageID(),
		CorrelationID: e.CorrelationID,
		CausationID:   e.MessageID,
		Time:          time.Now(),
		Message:       m,
	}

	return env
}

// Type returns the message type of the message contained in the envelope.
func (e Envelope) Type() MessageType {
	return TypeOf(e.Message)
}
