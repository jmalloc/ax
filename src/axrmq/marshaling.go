package axrmq

import (
	"fmt"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/marshaling"
	"github.com/streadway/amqp"
)

const (
	// createdAtHeader is the name of the AMQP message header that carries the
	// ax.Envelope.CreatedAt field.
	createdAtHeader = "ax-created-at"

	// sendAtHeader is the name of the AMQP message header that carries the
	// ax.Envelope.SendAt field.
	sendAtHeader = "ax-send-at"
)

// marshalMessage marshals a message envelope to an AMQP "publishing" message.
func marshalMessage(ep string, env endpoint.OutboundEnvelope) (amqp.Publishing, error) {
	pub := amqp.Publishing{
		AppId:         ep,
		MessageId:     env.MessageID.Get(),
		ReplyTo:       env.CausationID.Get(), // hijack reply-to for causation
		CorrelationId: env.CorrelationID.Get(),
		Timestamp:     time.Now(), // informational only, envelope times are in headers to retain TZ
		Type:          ax.TypeOf(env.Message).Name,
		Headers: amqp.Table{
			createdAtHeader: marshaling.MarshalTime(env.CreatedAt),
		},
	}

	// only add the delayed-until header if the message was actually delayed
	if env.SendAt.After(env.CreatedAt) {
		pub.Headers[sendAtHeader] = marshaling.MarshalTime(env.SendAt)
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
		DeliveryID:     endpoint.GenerateDeliveryID(),
		DeliveryCount:  countDeliveries(del),
	}

	if err := env.MessageID.Parse(del.MessageId); err != nil {
		return endpoint.InboundEnvelope{}, err
	}

	if err := env.CausationID.Parse(del.ReplyTo); err != nil {
		return endpoint.InboundEnvelope{}, err
	}

	if err := env.CorrelationID.Parse(del.CorrelationId); err != nil {
		return endpoint.InboundEnvelope{}, err
	}

	ok, err := unmarshalTimeFromHeader(del.Headers, createdAtHeader, &env.CreatedAt)
	if err != nil {
		return endpoint.InboundEnvelope{}, err
	} else if !ok {
		return endpoint.InboundEnvelope{}, fmt.Errorf("message %s does not contain a %s header", env.MessageID, createdAtHeader)
	}

	ok, err = unmarshalTimeFromHeader(del.Headers, sendAtHeader, &env.SendAt)
	if err != nil {
		return endpoint.InboundEnvelope{}, err
	} else if !ok {
		env.SendAt = env.CreatedAt
	}

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

	return count + 1
}

// unmarshalTimeFromHeader unmarshals a time from headers[n].
// It returns an error if the header is not a string, or contains an invalid time.
// It returns false if no header is present at all.
func unmarshalTimeFromHeader(headers amqp.Table, n string, t *time.Time) (bool, error) {
	v, ok := headers[n]
	if !ok {
		return false, nil
	}

	s, ok := v.(string)
	if !ok {
		return false, fmt.Errorf("%s header is not a string", n)
	}

	return true, marshaling.UnmarshalTime(s, t)
}
