package endpoint

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/internal/reflectx"
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
		"message.delay":          env.SendAt.Sub(env.CreatedAt).String(),
		"message.description":    env.Message.MessageDescription(),
	}
}

// traceError marks the given span as an error, and includes error information
// as a log event.
func traceError(span opentracing.Span, err error) {
	if span == nil {
		return
	}

	ext.Error.Set(span, true)

	span.LogFields(
		log.String("event", "error"),
		log.String("message", err.Error()),
		log.String("error.kind", reflectx.PrettyTypeName(err)),
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
			string(ext.Component):   "ax",
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
			string(ext.Component): "ax",
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

// traceInboundAccept adds a log event representing the time at which an
// inbound message was received.
func traceInboundAccept(span opentracing.Span) {
	if span == nil {
		return
	}

	span.LogFields(
		log.String("event", "endpoint.accept-inbound"),
		log.String("message", "the message is entering the inbound pipeline"),
	)
}

// traceInboundAck adds a log event representing the time at which an inbound
// message was acknowledged.
func traceInboundAck(span opentracing.Span) {
	if span == nil {
		return
	}

	span.LogFields(
		log.String("event", "transport.ack"),
		log.String("message", "the message has been processed successfully"),
	)
}

// traceInboundRetry adds a log event representing the time at which an inbound
// message was was returned to the queue to be retried.
func traceInboundRetry(span opentracing.Span, d time.Duration) {
	if span == nil {
		return
	}

	if d <= 0 {
		span.LogFields(
			log.String("event", "transport.retry"),
			log.String("message", "scheduling message for an immediate retry"),
		)
	} else {
		t := time.Now().Add(d)

		span.LogFields(
			log.String("event", "transport.retry"),
			log.String("message", "scheduling message for a delayed retry"),
			log.String("delay-for", d.String()),
			log.String("delay-until", t.Format(time.RFC3339Nano)),
		)
	}
}

// traceInboundReject adds a log event representing the time at which an inbound
// message was was rejected due to the retry policy.
func traceInboundReject(span opentracing.Span) {
	if span == nil {
		return
	}

	span.LogFields(
		log.String("event", "transport.reject"),
		log.String("message", "rejecting message due to retry policy"),
	)
}

// traceOutboundAccept adds a log event representing the time at which an outbound
// message was sent.
func traceOutboundAccept(span opentracing.Span) {
	if span == nil {
		return
	}

	span.LogFields(
		log.String("event", "endpoint.accept-outbound"),
		log.String("message", "the message is entering the outbound pipeline"),
	)
}

// traceOutboundSend adds a log event representing the time at which an outbound message
// is sent via the transport.
func traceOutboundSend(span opentracing.Span) {
	if span == nil {
		return
	}

	span.LogFields(
		log.String("event", "transport.send"),
		log.String("message", "sending the message via the transport"),
	)
}

// OutboundTracer is an implementation of OutboundPipeline that traces messages.
type OutboundTracer struct {
	Tracer opentracing.Tracer
	Next   OutboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further
// configure the endpoint as per the needs of the pipeline.
func (s OutboundTracer) Initialize(ctx context.Context, ep *Endpoint) error {
	return s.Next.Initialize(ctx, ep)
}

// Accept processes the message encapsulated in env.
func (s OutboundTracer) Accept(ctx context.Context, env OutboundEnvelope) error {
	span := startOutboundSpan(ctx, env, s.Tracer)
	defer span.Finish()

	env.SpanContext = span.Context()
	traceOutboundAccept(span)

	ctx = opentracing.ContextWithSpan(ctx, span)

	if err := s.Next.Accept(ctx, env); err != nil {
		traceError(span, err)
		return err
	}

	return nil
}
