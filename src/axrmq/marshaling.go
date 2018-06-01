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
		DeliveryCount:  countDeliveries(del),
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

// countDeliveries attempts to return the number of times the given message has
// been delievered. It returns zero if the count is unknown.
func countDeliveries(del amqp.Delivery) uint {
	death, ok := del.Headers["x-death"]

	if !ok {
		// The message has been redelivered, but there is no x-death header.
		// This can occur when the AMQP connection drops out without an explicit
		// rejection. We can't know the actual count in this case.
		if del.Redelivered {
			return 0 // unknown count
		}

		// The message has not been redelivered, and there is no x-death header,
		// so this must be the first attempt.
		return 1
	}

	// the x-death header should contain an AMQP array of AMQP tables.
	slice, ok := death.([]interface{})
	if !ok {
		return 0 // unknown count, unexpected header type
	}

	// sum the total of the DLX counts
	var count uint
	for _, v := range slice {
		if t, ok := v.(amqp.Table); ok {
			if n, ok := t["count"].(int64); ok {
				count += uint(n)
			}
		}
	}

	return count
}
