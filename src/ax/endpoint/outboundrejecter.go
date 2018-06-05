package endpoint

import (
	"context"
)

// OutboundRejecter is the validation stage in the message outbound pipeline
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

// Accept processes the message encapsulated in env.
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
