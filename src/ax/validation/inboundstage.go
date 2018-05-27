package validation

import (
	"context"

	"github.com/jmalloc/ax/src/ax/endpoint"
)

// InboundStage is the validation stage in the
// message inbound pipeline
type InboundStage struct {
	Validators []Validator
	Next       endpoint.InboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (s *InboundStage) Initialize(
	ctx context.Context,
	ep *endpoint.Endpoint,
) error {
	return s.Next.Initialize(ctx, ep)
}

// Accept forwards an inbound message through the pipeline until
// it is handled by some application-defined message handler(s).
func (s *InboundStage) Accept(
	ctx context.Context,
	sink endpoint.MessageSink,
	env endpoint.InboundEnvelope,
) error {

	for _, v := range s.Validators {
		if err := v.Validate(ctx, env.Envelope.Message); err != nil {
			return err
		}
	}
	return s.Next.Accept(ctx, sink, env)
}
