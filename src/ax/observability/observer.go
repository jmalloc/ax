package observability

import (
	"context"

	"github.com/jmalloc/ax/src/ax/endpoint"
)

// InboundObserver is an interface for types that observe inbound messages.
type InboundObserver interface {
	// InitializeInbound initializes the observer for inbound messages.
	InitializeInbound(ctx context.Context, ep *endpoint.Endpoint) error

	// BeforeInbound is called before a message is passed to the next pipeline stage.
	BeforeInbound(ctx context.Context, env endpoint.InboundEnvelope)

	// AfterInbound is called after a message is accepted by the next pipeline stage.
	// err is the error returned by the next pipeline stage, which may be nil.
	AfterInbound(ctx context.Context, env endpoint.InboundEnvelope, err error)
}

// OutboundObserver is an interface for types that observe outbound messages.
type OutboundObserver interface {
	// InitializeOutbound initializes the observer for outbound messages.
	InitializeOutbound(ctx context.Context, ep *endpoint.Endpoint) error

	// BeforeOutbound is called before a message is passed to the next pipeline stage.
	BeforeOutbound(ctx context.Context, env endpoint.OutboundEnvelope)

	// AfterOutbound is called after a message is accepted by the next pipeline stage.
	// err is the error returned by the next pipeline stage, which may be nil.
	AfterOutbound(ctx context.Context, env endpoint.OutboundEnvelope, err error)
}
