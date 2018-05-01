package axrmq

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/marshaling"
	"github.com/streadway/amqp"
)

// marshalMessage marshals a message envelope to an AMQP "publishing" message.
func marshalMessage(ep string, env bus.OutboundEnvelope, pub *amqp.Publishing) error {
	pub.AppId = ep
	pub.MessageId = env.MessageID.Get()
	pub.ReplyTo = env.CausationID.Get() // hijack reply-to for causation
	pub.CorrelationId = env.CorrelationID.Get()
	pub.Timestamp = env.Time
	pub.Type = ax.TypeOf(env.Message).Name

	var err error
	pub.ContentType, pub.Body, err = marshaling.MarshalMessage(env.Message)

	return err
}

// unmarshalMessage unmarshals a message envelope from an AMQP "delivery"
// message.
func unmarshalMessage(del amqp.Delivery, env *bus.InboundEnvelope) error {
	env.SourceEndpoint = del.AppId

	if err := env.MessageID.Parse(del.MessageId); err != nil {
		return err
	}

	if err := env.CausationID.Parse(del.ReplyTo); err != nil {
		return err
	}

	if err := env.CorrelationID.Parse(del.CorrelationId); err != nil {
		return err
	}

	env.Time = del.Timestamp

	var err error
	env.Message, err = marshaling.UnmarshalMessage(del.ContentType, del.Body)

	return err
}
