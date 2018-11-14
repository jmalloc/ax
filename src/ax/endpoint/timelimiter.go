package endpoint

import (
	"context"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// DefaultTimeout is the default timeout duration to use if none is given.
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

	traceTimeLimit(
		opentracing.SpanFromContext(ctx),
		to,
	)

	return tl.Next.Accept(ctx, sink, env)
}

func traceTimeLimit(span opentracing.Span, d time.Duration) {
	if span == nil {
		return
	}

	span.LogFields(
		log.String("event", "time-limiter.limit"),
		log.String("message", "added timeout, forwarding message to the next pipeline stage"),
		log.String("timeout", d.String()),
	)
}
