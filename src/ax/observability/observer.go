package observability

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// BeforeInboundObserver is an interface for observers that are notified before
// the inbound pipeline accepts a message.
type BeforeInboundObserver interface {
	// BeforeInbound is called before a message is passed to the next pipeline stage.
	// The returned context is passed to the next pipeline stage.
	BeforeInbound(context.Context, ax.Envelope) (context.Context, error)
}

// AfterInboundObserver is an interface for observers that are notified after
// the inbound pipeline accepts a message.
type AfterInboundObserver interface {
	// AfterInbound is called after a message is accepted by the next pipeline stage.
	// It is passed the error that occurred, if any,
	AfterInbound(context.Context, ax.Envelope, error) error
}

// BeforeOutboundObserver is an interface for observers that are notified before
// the outbound pipeline accepts a message.
type BeforeOutboundObserver interface {
	// BeforeOutbound is called before a message is passed to the next pipeline stage.
	// The returned context is passed to the next pipeline stage.
	BeforeOutbound(context.Context, ax.Envelope) (context.Context, error)
}

// AfterOutboundObserver is an interface for observers that are notified after
// the outbound pipeline accepts a message.
type AfterOutboundObserver interface {
	// AfterOutbound is called after a message is accepted by the next pipeline stage.
	// It is passed the error that occurred, if any,
	AfterOutbound(context.Context, ax.Envelope, error) error
}
