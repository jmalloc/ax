package endpoint

import (
	"context"

	"github.com/jmalloc/ax/src/internal/reflectx"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// InboundRejecter is an inbound pipeline stage that validates messages before
// forwarding them to the next pipeline stage. It uses a set of validators
// distinct from those configured in the endpoint.
type InboundRejecter struct {
	Validators []Validator
	Next       InboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (i *InboundRejecter) Initialize(
	ctx context.Context,
	ep *Endpoint,
) error {
	return i.Next.Initialize(ctx, ep)
}

// Accept forwards an inbound message to the next pipeline stage only if it is
// successfully validated.
func (i *InboundRejecter) Accept(
	ctx context.Context,
	sink MessageSink,
	env InboundEnvelope,
) error {
	span := opentracing.SpanFromContext(ctx)

	for _, v := range i.Validators {
		traceValidate(span, v)

		if err := v.Validate(ctx, env.Message); err != nil {
			return err
		}
	}

	traceValidated(span)

	return i.Next.Accept(ctx, sink, env)
}

// OutboundRejecter is an outbound pipeline stage that validates messages before
// forwarding them to the next pipeline stage. It uses a set of validators
// distinct from those configured in the endpoint.
type OutboundRejecter struct {
	Validators []Validator
	Next       OutboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (o *OutboundRejecter) Initialize(
	ctx context.Context,
	ep *Endpoint,
) error {
	return o.Next.Initialize(ctx, ep)
}

// Accept forwards an outbound message to the next pipeline stage only if it is
// successfully validated.
func (o *OutboundRejecter) Accept(
	ctx context.Context,
	env OutboundEnvelope,
) error {
	span := opentracing.SpanFromContext(ctx)

	for _, v := range o.Validators {
		traceValidate(span, v)

		if err := v.Validate(ctx, env.Message); err != nil {
			return err
		}
	}

	traceValidated(span)

	return o.Next.Accept(ctx, env)
}

func traceValidate(span opentracing.Span, v Validator) {
	if span == nil {
		return
	}

	span.LogFields(
		log.String("event", "rejecter.validate"),
		log.String("message", "validating the message"),
		log.String("validator", reflectx.PrettyTypeName(v)),
	)
}

func traceValidated(span opentracing.Span) {
	if span == nil {
		return
	}

	span.LogFields(
		log.String("event", "rejecter.validated"),
		log.String("message", "the message is valid, forwarding message to the next pipeline stage"),
	)
}
