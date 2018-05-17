package endpoint

import (
	"context"
)

// Fanout is an inbound pipeline that forwards inbound messages to zero-or-more
// pipeline stages in parallel.
type Fanout []InboundPipeline

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (f Fanout) Initialize(ctx context.Context, ep *Endpoint) error {
	for _, p := range f {
		if err := p.Initialize(ctx, ep); err != nil {
			return err
		}
	}

	return nil
}

// Accept forwards an inbound message through the pipeline until
// it is handled by some application-defined message handler(s).
func (f Fanout) Accept(ctx context.Context, sink MessageSink, env InboundEnvelope) error {
	for _, p := range f {
		if err := p.Accept(ctx, sink, env); err != nil {
			return err
		}
	}

	return nil
}
