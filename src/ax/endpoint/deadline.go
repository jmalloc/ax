package endpoint

import (
	"context"
)

// Deadline is an inbound pipeline that sets a context deadline before
// forwarding on to the next stage.
type Deadline struct {
	Deadline time.Time
	Cancel   context.CancelFunc
	Next     endpoint.InboundPipeline
}

const DefaultDeadlineMs = 1000

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (d Deadline) Initialize(ctx context.Context, ep *Endpoint) error {
	if d.Deadline.IsZero() {
		d.Deadline = DefaultDeadlineMs * time.Millisecond
	}
	return i.Next.Initialize(ctx, ep)
}

// Accept forwards an inbound message through the pipeline until
// it is handled by some application-defined message handler(s).
func (d Deadline) Accept(ctx context.Context, sink MessageSink, env InboundEnvelope) error {
	ctx, d.Cancel = context.WithDeadline(ctx, d.Deadline)
	go cancelDeadline(ctx, d)
	return i.Next.Accept(ctx, sink, env)
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
