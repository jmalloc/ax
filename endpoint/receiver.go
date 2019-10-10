package endpoint

import (
	"context"
	"time"

	"github.com/jmalloc/ax/internal/servicegroup"
	"github.com/jmalloc/ax/internal/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// receiver receives a message from a transport, forwards it to the inbound
// pipeline, then acknowledges the message.
type receiver struct {
	Transport        InboundTransport
	InboundPipeline  InboundPipeline
	OutboundPipeline OutboundPipeline
	RetryPolicy      RetryPolicy
	Tracer           opentracing.Tracer

	wg *servicegroup.Group
}

// Run processes inbound messages until ctx is canceled or an error occurrs.
func (r *receiver) Run(ctx context.Context) error {
	if r.RetryPolicy == nil {
		r.RetryPolicy = DefaultRetryPolicy
	}

	r.wg = servicegroup.NewGroup(ctx)

	if err := r.wg.Go(r.receive); err != nil {
		return err
	}

	return r.wg.Wait()
}

// receive starts a new goroutine to process each inbound message.
func (r *receiver) receive(ctx context.Context) error {
	for {
		env, ack, err := r.Transport.Receive(ctx)
		if err != nil {
			return err
		}

		if err := r.wg.Go(func(ctx context.Context) error {
			return r.process(ctx, env, ack)
		}); err != nil {
			return err
		}
	}
}

func (r *receiver) process(
	ctx context.Context,
	env InboundEnvelope,
	ack Acknowledger,
) error {
	span := startInboundSpan(env, r.Tracer)
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)

	tracing.LogEventS(
		span,
		"receive",
		"the message has been received from the transport",
	)

	acceptErr := r.InboundPipeline.Accept(
		WithEnvelope(ctx, env),
		r.OutboundPipeline,
		env,
	)

	if acceptErr != nil {
		tracing.LogErrorS(span, acceptErr)
	}

	if err := r.ack(ctx, ack, env, acceptErr); err != nil {
		tracing.LogErrorS(span, err)
		return err
	}

	return nil
}

func (r *receiver) ack(
	ctx context.Context,
	ack Acknowledger,
	env InboundEnvelope,
	err error,
) error {
	if err == nil {
		tracing.LogEvent(
			ctx,
			"ack",
			"acknowledging successfully processed message",
		)

		return ack.Ack(ctx)
	}

	if d, ok := r.RetryPolicy(env, err); ok {
		tracing.LogEvent(
			ctx,
			"retry",
			"scheduling failed message for retry",
			tracing.Duration("delay_for", d),
			tracing.Time("delay_until", time.Now().Add(d)),
		)

		return ack.Retry(ctx, err, d)
	}

	tracing.LogEvent(
		ctx,
		"reject",
		"rejecting failed message as per retry policy",
	)

	return ack.Reject(ctx, err)
}

// startInboundSpan starts an OpenTracing span representing an inbound message.
func startInboundSpan(env InboundEnvelope, tr opentracing.Tracer) opentracing.Span {
	opts := []opentracing.StartSpanOption{
		ext.SpanKindConsumer,
		spanTagsForEnvelope(env.Envelope),
		opentracing.Tags{
			string(ext.Component): "ax",
			"source_endpoint":     env.SourceEndpoint,
			"attempt_id":          env.AttemptID.Get(),
			"attempt_short_id":    env.AttemptID.String(),
			"attempt_count":       env.AttemptCount,
		},
	}

	if env.SpanContext != nil {
		opts = append(
			opts,
			opentracing.ChildOf(env.SpanContext),
		)
	}

	return tracing.StartSpan(
		tr,
		env.Type().String(),
		opts...,
	)
}
