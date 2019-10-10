package ax

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

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

	// CreatedAt is the time at which the message was created. This typically
	// correlates to the time at which the message was passed to a Sender.
	//
	// It is populated via the regular Go system clock, and as such there are
	// almost no guarantees about the accuracy of the time. It must not be
	// assumed that messages will arrive in any chronological order.
	//
	// Depending on the application this field may not be appropriate for use as an
	// "occurred" time. Care must be taken to choose appropriately between
	// CreatedAt and SendAt for each use case.
	CreatedAt time.Time

	// SendAt is the time at which the message should be sent by the endpoint,
	// which may be after the CreatedAt time.
	//
	// Depending on the application this field may not be appropriate for use as an
	// "occurred" time. Care must be taken to choose appropriately between
	// CreatedAt and SendAt for each use case.
	SendAt time.Time

	// Message is the application-defined message encapsulated by the envelope.
	Message Message
}

// NewEnvelope creates a new message envelope containing m.
//
// It generates a UUID-based message ID and configures the envelope such that m
// is at the root of a new tree of messages.
func NewEnvelope(m Message) Envelope {
	id := GenerateMessageID()
	t := time.Now()

	return Envelope{
		MessageID:     id,
		CausationID:   id,
		CorrelationID: id,
		CreatedAt:     t,
		SendAt:        t,
		Message:       m,
	}
}

// NewEnvelopeFromProto returns a new envelope from its protocol-buffers
// representation.
func NewEnvelopeFromProto(env *EnvelopeProto) (Envelope, error) {
	messageID, err := ParseMessageID(env.MessageId)
	if err != nil {
		return Envelope{}, err
	}

	causationID, err := ParseMessageID(env.CausationId)
	if err != nil {
		return Envelope{}, err
	}

	correlationID, err := ParseMessageID(env.CorrelationId)
	if err != nil {
		return Envelope{}, err
	}

	createdAt, err := ptypes.Timestamp(env.CreatedAt)
	if err != nil {
		return Envelope{}, err
	}

	sendAt, err := ptypes.Timestamp(env.SendAt)
	if err != nil {
		return Envelope{}, err
	}

	var any ptypes.DynamicAny
	err = ptypes.UnmarshalAny(env.Message, &any)
	if err != nil {
		return Envelope{}, err
	}

	message, ok := any.Message.(Message)
	if !ok {
		return Envelope{}, fmt.Errorf(
			"%s does not implement Message",
			reflect.TypeOf(any.Message),
		)
	}

	return Envelope{
		MessageID:     messageID,
		CausationID:   causationID,
		CorrelationID: correlationID,
		CreatedAt:     createdAt,
		SendAt:        sendAt,
		Message:       message,
	}, nil
}

// NewChild returns a new message envelope containing m.
//
// It generates a UUID-based message ID and configures the envelope such that
// m is a child of e.Message within an existing tree of messages.
func (e Envelope) NewChild(m Message) Envelope {
	t := time.Now()

	return Envelope{
		MessageID:     GenerateMessageID(),
		CorrelationID: e.CorrelationID,
		CausationID:   e.MessageID,
		CreatedAt:     t,
		SendAt:        t,
		Message:       m,
	}
}

// Type returns the message type of the message contained in the envelope.
func (e Envelope) Type() MessageType {
	return TypeOf(e.Message)
}

// Delay returns the delay between the messages creation time and the time at
// which it is to be sent.
func (e Envelope) Delay() time.Duration {
	if e.SendAt.After(e.CreatedAt) {
		return e.SendAt.Sub(e.CreatedAt)
	}

	return 0
}

// Equal returns true if e and env contain the same data.
func (e Envelope) Equal(env Envelope) bool {
	return e.MessageID == env.MessageID &&
		e.CorrelationID == env.CorrelationID &&
		e.CausationID == env.CausationID &&
		e.CreatedAt.Equal(env.CreatedAt) &&
		e.SendAt.Equal(env.SendAt) &&
		proto.Equal(e.Message, env.Message)
}

// AsProto returns a Protocol Buffers representation of the envelope.
func (e Envelope) AsProto() (*EnvelopeProto, error) {
	createdAt, err := ptypes.TimestampProto(e.CreatedAt)
	if err != nil {
		return nil, err
	}

	sendAt, err := ptypes.TimestampProto(e.SendAt)
	if err != nil {
		return nil, err
	}

	message, err := ptypes.MarshalAny(e.Message)
	if err != nil {
		return nil, err
	}

	return &EnvelopeProto{
		MessageId:     e.MessageID.String(),
		CausationId:   e.CausationID.String(),
		CorrelationId: e.CorrelationID.String(),
		CreatedAt:     createdAt,
		SendAt:        sendAt,
		Message:       message,
	}, nil
}

// MarshalEnvelope marshals env to a binary representation.
func MarshalEnvelope(env Envelope) ([]byte, error) {
	v, err := env.AsProto()
	if err != nil {
		return nil, err
	}

	return proto.Marshal(v)
}

// UnmarshalEnvelope unmarshals an envelope from its serialized representation.
func UnmarshalEnvelope(data []byte) (Envelope, error) {
	var v EnvelopeProto
	if err := proto.Unmarshal(data, &v); err != nil {
		return Envelope{}, err
	}

	return NewEnvelopeFromProto(&v)
}
