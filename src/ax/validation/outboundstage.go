package validation

import (
	"context"

	"github.com/jmalloc/ax/src/ax/endpoint"
)

// OutboundStage is the validation stage in the
// message outbound pipeline
type OutboundStage struct {
	Validators []endpoint.Validator
	Next       endpoint.OutboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (s *OutboundStage) Initialize(
	ctx context.Context,
	ep *endpoint.Endpoint,
) error {
	return s.Next.Initialize(ctx, ep)
}

// Accept processes the message encapsulated in env.
func (s *OutboundStage) Accept(
	ctx context.Context,
	env endpoint.OutboundEnvelope,
) error {
	for _, v := range s.Validators {
		if err := v.Validate(ctx, env.Envelope.Message); err != nil {
			return err
		}
	}
	return s.Next.Accept(ctx, env)
}
