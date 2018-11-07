package axrmq

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/marshaling"
	"github.com/jmalloc/ax/src/internal/bufferpool"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
)

const (
	// createdAtHeader is the name of the AMQP message header that carries the
	// ax.Envelope.CreatedAt field.
	createdAtHeader = "ax-created-at"

	// sendAtHeader is the name of the AMQP message header that carries the
	// ax.Envelope.SendAt field.
	sendAtHeader = "ax-send-at"

	// spanContextHeader is the name of the AMQP message header that carries the
	// OpenTracing span-context.
	spanContextHeader = "ax-span-context"
)

// marshalMessage marshals a message envelope to an AMQP "publishing" message.
func marshalMessage(
	ep string,
	env endpoint.OutboundEnvelope,
	tr opentracing.Tracer,
) (amqp.Publishing, error) {
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

	marshalSpanContext(pub.Headers, env.SpanContext, tr)

	var err error
	pub.ContentType, pub.Body, err = ax.MarshalMessage(env.Message)

	return pub, err
}

// unmarshalMessage unmarshals a message envelope from an AMQP "delivery"
// message.
func unmarshalMessage(
	del amqp.Delivery,
	tr opentracing.Tracer,
) (endpoint.InboundEnvelope, error) {
	env := endpoint.InboundEnvelope{
		SourceEndpoint: del.AppId,
		AttemptID:      endpoint.GenerateAttemptID(),
		AttemptCount:   countAttempts(del),
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

	env.SpanContext = unmarshalSpanContext(del.Headers, tr)

	env.Message, err = ax.UnmarshalMessage(del.ContentType, del.Body)

	return env, err
}

// countAttempts returns the number of times the given message has been
// delievered. It returns zero if the count is unknown.
func countAttempts(del amqp.Delivery) uint {
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

// marshalSpanContext marshals an OpenTracing span context into AMQP headers.
func marshalSpanContext(
	headers amqp.Table,
	sc opentracing.SpanContext,
	tr opentracing.Tracer,
) {
	if tr == nil || sc == nil {
		return
	}

	// buf is not returned to the pool, as the underlying slice is retained in the
	// AMQP headers.
	buf := bufferpool.Get()

	if err := tr.Inject(
		sc,
		opentracing.Binary,
		buf,
	); err != nil {
		// TODO(jmalloc): add logging - https://github.com/jmalloc/ax/issues/144
		// A failed span propagation is worth logging about, but does not stop the
		// system from sending the message.
		return
	}

	headers[spanContextHeader] = buf.Bytes()
}

// unmarshalSpanContext unmarshals an OpenTracing span context from AMQP headers.
func unmarshalSpanContext(
	headers amqp.Table,
	tr opentracing.Tracer,
) opentracing.SpanContext {
	if tr == nil {
		return nil
	}

	v, ok := headers[spanContextHeader]
	if !ok {
		return nil
	}

	b, ok := v.([]byte)
	if !ok {
		// TODO(jmalloc): add logging - https://github.com/jmalloc/ax/issues/144
		// A failed span propagation is worth logging about, but does not stop the
		// system from sending the message.
		return nil
	}

	buf := bytes.NewBuffer(b)
	defer bufferpool.Put(buf)

	sc, err := tr.Extract(opentracing.Binary, buf)

	if err == opentracing.ErrSpanContextNotFound {
		return nil
	} else if err != nil {
		// TODO(jmalloc): add logging - https://github.com/jmalloc/ax/issues/144
		// A failed span propagation is worth logging about, but does not stop the
		// system from sending the message.
		return nil
	}

	return sc
}
