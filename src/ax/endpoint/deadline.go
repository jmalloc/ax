package endpoint

import (
	"context"
	"time"
)

// Deadline is an inbound pipeline that sets a context deadline before
// forwarding on to the next stage.
type Deadline struct {
	Deadline time.Duration
	Cancel   context.CancelFunc
	Next     InboundPipeline
}

const DefaultDeadline = 1000 * time.Millisecond

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (d Deadline) Initialize(ctx context.Context, ep *Endpoint) error {
	if d.Deadline == 0 {
		d.Deadline = DefaultDeadline
	}
	return d.Next.Initialize(ctx, ep)
}

// Accept forwards an inbound message through the pipeline until
// it is handled by some application-defined message handler(s).
func (d Deadline) Accept(ctx context.Context, sink MessageSink, env InboundEnvelope) error {
	ctx, d.Cancel = context.WithTimeout(ctx, d.Deadline)
	go cancelDeadline(ctx, d)
	return d.Next.Accept(ctx, sink, env)
}

func cancelDeadline(ctx context.Context, d Deadline) {
	defer d.Cancel()

	select {
	case <-time.After(d.Deadline):
		// deadline reached
	case <-ctx.Done():
		// log ctx.Err(), if non nil
	}
}
