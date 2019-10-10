package endpoint

import (
	"context"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/internal/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// OutboundTracer is an implementation of OutboundPipeline that starts a new
// OpenTracing span before forwarding to the next stage.
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
	span := tracing.StartChildOf(
		ctx,
		s.Tracer,
		env.Type().String(),
		ext.SpanKindProducer,
		spanTagsForEnvelope(env.Envelope),
	)
	defer span.Finish()

	switch env.Operation {
	case OpSendUnicast:
		span.SetTag("operation", "send-unicast")
	case OpSendMulticast:
		span.SetTag("operation", "send-multicast")
	}

	ctx = opentracing.ContextWithSpan(ctx, span)

	if err := s.Next.Accept(ctx, env); err != nil {
		tracing.LogErrorS(span, err)
		return err
	}

	return nil
}

// spanTagsForEnvelope returns the OpenTracing span tags describing a message
// envelope. It is used for both inbound and outbound traces.
func spanTagsForEnvelope(env ax.Envelope) opentracing.Tags {
	return opentracing.Tags{
		"message_id":           env.MessageID.Get(),
		"causation_id":         env.CausationID.Get(),
		"correlation_id":       env.CorrelationID.Get(),
		"message_short_id":     env.MessageID.String(),
		"causation_short_id":   env.CausationID.String(),
		"correlation_short_id": env.CorrelationID.String(),
		"description":          env.Message.MessageDescription(),
		"created_at":           env.CreatedAt,
		"send_at":              env.SendAt,
		"delay":                env.SendAt.Sub(env.CreatedAt).String(),
	}
}
