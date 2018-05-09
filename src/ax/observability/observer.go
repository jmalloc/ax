package observability

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// BeforeInboundObserver is an interface for observers that are notified before
// the inbound pipeline accepts a message.
type BeforeInboundObserver interface {
	// BeforeInbound is called before a message is passed to the next pipeline stage.
	// The returned context replaces ctx in calls to other observers and the next pipeline stage.
	BeforeInbound(ctx context.Context, env ax.Envelope) context.Context
}

// AfterInboundObserver is an interface for observers that are notified after
// the inbound pipeline accepts a message.
type AfterInboundObserver interface {
	// AfterInbound is called after a message is accepted by the next pipeline stage.
	// err is the error returned by the next pipeline stage, which may be nil.
	AfterInbound(ctx context.Context, env ax.Envelope, err error)
}

// BeforeOutboundObserver is an interface for observers that are notified before
// the outbound pipeline accepts a message.
type BeforeOutboundObserver interface {
	// BeforeOutbound is called before a message is passed to the next pipeline stage.
	// The returned context replaces ctx in calls to other observers and the next pipeline stage.
	BeforeOutbound(ctx context.Context, env ax.Envelope) context.Context
}

// AfterOutboundObserver is an interface for observers that are notified after
// the outbound pipeline accepts a message.
type AfterOutboundObserver interface {
	// AfterOutbound is called after a message is accepted by the next pipeline stage.
	// err is the error returned by the next pipeline stage, which may be nil.
	AfterOutbound(ctx context.Context, env ax.Envelope, err error)
}
