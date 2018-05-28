package endpoint

import (
	"context"
	"time"
)

// The default timeout duration if none given.
const DefaultTimeout = 5 * time.Second

// TimeLimiter is an inbound pipeline that sets a context timeout before
// forwarding on to the next stage.
type TimeLimiter struct {
	Timeout time.Duration
	Next    InboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (tl TimeLimiter) Initialize(ctx context.Context, ep *Endpoint) error {
	return tl.Next.Initialize(ctx, ep)
}

// Accept forwards an inbound message through the pipeline until
// it is handled by some application-defined message handler(s).
func (tl TimeLimiter) Accept(ctx context.Context, sink MessageSink, env InboundEnvelope) error {
	to := tl.Timeout
	if to == 0 {
		to = DefaultTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, to)
	defer cancel()
	return tl.Next.Accept(ctx, sink, env)
}
