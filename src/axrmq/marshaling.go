package axrmq

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/streadway/amqp"
)

// marshalMessage marshals a message envelope to an AMQP "publishing" message.
func marshalMessage(ep string, env endpoint.OutboundEnvelope) (amqp.Publishing, error) {
	pub := amqp.Publishing{
		AppId:         ep,
		MessageId:     env.MessageID.Get(),
		ReplyTo:       env.CausationID.Get(), // hijack reply-to for causation
		CorrelationId: env.CorrelationID.Get(),
		Timestamp:     env.Time,
		Type:          ax.TypeOf(env.Message).Name,
	}

	var err error
	pub.ContentType, pub.Body, err = ax.MarshalMessage(env.Message)

	return pub, err
}

// unmarshalMessage unmarshals a message envelope from an AMQP "delivery"
// message.
func unmarshalMessage(del amqp.Delivery) (endpoint.InboundEnvelope, error) {
	env := endpoint.InboundEnvelope{
		SourceEndpoint: del.AppId,
	}

	if err := env.MessageID.Parse(del.MessageId); err != nil {
		return env, err
	}

	if err := env.CausationID.Parse(del.ReplyTo); err != nil {
		return env, err
	}

	if err := env.CorrelationID.Parse(del.CorrelationId); err != nil {
		return env, err
	}

	env.Time = del.Timestamp

	var err error
	env.Message, err = ax.UnmarshalMessage(del.ContentType, del.Body)

	return env, err
}
