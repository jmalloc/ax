package endpoint

import (
	"context"
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
	for _, v := range i.Validators {
		if err := v.Validate(ctx, env.Message); err != nil {
			return err
		}
	}

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
	for _, v := range o.Validators {
		if err := v.Validate(ctx, env.Message); err != nil {
			return err
		}
	}

	return o.Next.Accept(ctx, env)
}
