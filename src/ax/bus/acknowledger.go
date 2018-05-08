package bus

import (
	"context"
)

// RetryPolicy returns true if the message should be retried.
type RetryPolicy func(InboundEnvelope) bool

// DefaultRetryPolicy is a RetryPolicy that rejects a messages after it has been
// attempted three (3) times.
func DefaultRetryPolicy(m InboundEnvelope) bool {
	return m.DeliveryCount < 3
}

// Acknowledger is an implementation of InboundPipeline that forwards a message
// on to the next stage and then acknowledges the message.
type Acknowledger struct {
	RetryPolicy RetryPolicy
	Next        InboundPipeline
}

// Initialize sets a default retry policy and forwards the message on to the
// next stage for initialization.
func (a *Acknowledger) Initialize(ctx context.Context, t Transport) error {
	if a.RetryPolicy == nil {
		a.RetryPolicy = DefaultRetryPolicy
	}
	return a.Next.Initialize(ctx, t)
}

// Accept forwards the message on to the next stage and then marks it as
// done based on the result of that next stage.
func (a *Acknowledger) Accept(
	ctx context.Context,
	s MessageSink,
	m InboundEnvelope,
) error {
	op := a.accept(ctx, s, m)
	return m.Done(ctx, op)
}

func (a *Acknowledger) accept(
	ctx context.Context,
	s MessageSink,
	m InboundEnvelope,
) InboundOperation {
	if err := a.Next.Accept(ctx, s, m); err != nil {
		if a.RetryPolicy(m) {
			return OpRetry
		}
		return OpReject
	}
	return OpAck
}
