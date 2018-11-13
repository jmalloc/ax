package endpoint

import (
	"context"
	"reflect"
	"time"

	"github.com/jmalloc/ax/src/ax"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

// getTracer returns tr if it is non-nil, otherwise it returns a no-op tracer.
func getTracer(tr opentracing.Tracer) opentracing.Tracer {
	if tr == nil {
		return opentracing.NoopTracer{}
	}

	return tr
}

// spanTagsForEnvelope returns the OpenTracing span tags describing the given
// message envelope.
func spanTagsForEnvelope(env ax.Envelope) opentracing.Tags {
	return opentracing.Tags{
		"message.id":             env.MessageID.Get(),
		"message.causation_id":   env.CausationID.Get(),
		"message.correlation_id": env.CorrelationID.Get(),
		"message.created_at":     env.CreatedAt,
		"message.send_at":        env.SendAt,
		"message.description":    env.Message.MessageDescription(),
	}
}

// traceError marks the given span as an error, and includes error information
// as a log event.
func traceError(span opentracing.Span, err error) {
	ext.Error.Set(span, true)

	span.LogFields(
		log.String("event", "error"),
		log.String("message", err.Error()),
		log.String("error.kind", reflect.TypeOf(err).String()),
	)
}

// startInboundSpan starts an OpenTracing span representing an inbound message.
func startInboundSpan(
	ctx context.Context,
	env InboundEnvelope,
	tr opentracing.Tracer,
) opentracing.Span {
	opts := []opentracing.StartSpanOption{
		ext.SpanKindConsumer,
		spanTagsForEnvelope(env.Envelope),
		opentracing.Tags{
			string(ext.Component):   "ax.endpoint",
			"message.source":        env.SourceEndpoint,
			"message.attempt.id":    env.AttemptID.Get(),
			"message.attempt.count": env.AttemptCount,
		},
	}

	if env.SpanContext != nil {
		opts = append(
			opts,
			opentracing.ChildOf(env.SpanContext),
		)
	}

	return getTracer(tr).StartSpan(
		env.Type().String(),
		opts...,
	)
}

// startOutboundSpan starts an OpenTracing span representing an outbound message.
func startOutboundSpan(
	ctx context.Context,
	env OutboundEnvelope,
	tr opentracing.Tracer,
) opentracing.Span {
	opts := []opentracing.StartSpanOption{
		ext.SpanKindProducer,
		spanTagsForEnvelope(env.Envelope),
		opentracing.Tags{
			string(ext.Component): "ax.endpoint",
		},
	}

	if p := opentracing.SpanFromContext(ctx); p != nil {
		opts = append(
			opts,
			opentracing.ChildOf(p.Context()),
		)
	}

	span := getTracer(tr).StartSpan(
		env.Type().String(),
		opts...,
	)

	switch env.Operation {
	case OpSendUnicast:
		span.SetTag("message.operation", "send-unicast")
	case OpSendMulticast:
		span.SetTag("message.operation", "send-multicast")
	}

	return span
}

// traceReceive adds a log event representing the time at which an
// inbound message was received.
func traceReceive(span opentracing.Span) {
	span.LogFields(
		log.String("event", "endpoint.receive"),
		log.String("message", "the message is entering the inbound pipeline"),
	)
}

// traceAck adds a log event representing the time at which an inbound
// message was acknowledged.
func traceAck(span opentracing.Span) {
	span.LogFields(
		log.String("event", "endpoint.ack"),
		log.String("message", "the message has been processed successfully"),
	)
}

// traceRetry adds a log event representing the time at which an inbound
// message was was returned to the queue to be retried.
func traceRetry(span opentracing.Span, d time.Duration) {
	if d <= 0 {
		span.LogFields(
			log.String("event", "endpoint.retry"),
			log.String("message", "scheduling message for an immediate retry"),
		)
	} else {
		t := time.Now().Add(d)

		span.LogFields(
			log.String("event", "endpoint.retry"),
			log.String("message", "scheduling message for a delayed retry"),
			log.String("delay-for", d.String()),
			log.String("delay-until", t.Format(time.RFC3339Nano)),
		)
	}
}

// traceReject adds a log event representing the time at which an inbound
// message was was rejected due to the retry policy.
func traceReject(span opentracing.Span) {
	span.LogFields(
		log.String("event", "endpoint.reject"),
		log.String("message", "rejecting message due to retry policy"),
	)
}

// traceSend adds a log event representing the time at which an outbound
// message was sent.
func traceSend(span opentracing.Span, op Operation) {
	span.LogFields(
		log.String("event", "endpoint.send"),
		log.String("message", "the message is entering the outbound pipeline"),
	)
}

// tracingSink is an implementation of MessageSink that traces messages.
type tracingSink struct {
	Tracer opentracing.Tracer
	Next   MessageSink
}

// Accept processes the message encapsulated in env.
func (s tracingSink) Accept(ctx context.Context, env OutboundEnvelope) error {
	span := startOutboundSpan(ctx, env, s.Tracer)
	defer span.Finish()

	env.SpanContext = span.Context()
	traceSend(span, env.Operation)

	if err := s.Next.Accept(ctx, env); err != nil {
		traceError(span, err)
		return err
	}

	return nil
}
