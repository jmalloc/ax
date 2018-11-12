package endpoint

import (
	"context"
	"reflect"

	"github.com/jmalloc/ax/src/ax"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

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
			string(ext.Component):   "ax",
			"message.pipeline":      "inbound",
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
	env ax.Envelope,
	tr opentracing.Tracer,
) opentracing.Span {
	opts := []opentracing.StartSpanOption{
		ext.SpanKindProducer,
		spanTagsForEnvelope(env),
		opentracing.Tags{
			string(ext.Component): "ax",
			"message.pipeline":    "outbound",
		},
	}

	if p := opentracing.SpanFromContext(ctx); p != nil {
		opts = append(
			opts,
			opentracing.ChildOf(p.Context()),
		)
	}

	return getTracer(tr).StartSpan(
		env.Type().String(),
		opts...,
	)
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

// traceReceive adds a log event representing the time at which an inbound
// message was received.
func traceReceive(span opentracing.Span) {
	span.LogFields(log.String("event", "receive"))
}

// traceAck adds a log event representing the time at which an inbound message
// was acknowledged.
func traceAck(span opentracing.Span) {
	span.LogFields(log.String("event", "ack"))
}

// traceRetry adds a log event representing the time at which an inbound message
// was was returned to the queue to be retried.
func traceRetry(span opentracing.Span) {
	span.LogFields(log.String("event", "retry"))
}

// traceReject adds a log event representing the time at which an inbound
// message was was rejected due to the retry policy.
func traceReject(span opentracing.Span) {
	span.LogFields(log.String("event", "reject"))
}

// traceSend adds a log event representing the time at which an outbound message
// was sent.
func traceSend(span opentracing.Span) {
	span.LogFields(log.String("event", "send"))
}

// getTracer returns tr if it is non-nil, otherwise it returns a no-op tracer.
func getTracer(tr opentracing.Tracer) opentracing.Tracer {
	if tr == nil {
		return opentracing.NoopTracer{}
	}

	return tr
}
