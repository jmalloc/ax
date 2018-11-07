package observability

import (
	"context"

	"github.com/jmalloc/ax/src/ax/endpoint"
)

// Observer is an interface for types that observe the messages sent and
// received by an endpoint.
type Observer interface {
	// BeforeInbound is called before a message is passed to the next pipeline stage.
	BeforeInbound(ctx context.Context, env endpoint.InboundEnvelope)

	// AfterInbound is called after a message is accepted by the next pipeline stage.
	// err is the error returned by the next pipeline stage, which may be nil.
	AfterInbound(ctx context.Context, env endpoint.InboundEnvelope, err error)

	// BeforeOutbound is called before a message is passed to the next pipeline stage.
	BeforeOutbound(ctx context.Context, env endpoint.OutboundEnvelope)

	// AfterOutbound is called after a message is accepted by the next pipeline stage.
	// err is the error returned by the next pipeline stage, which may be nil.
	AfterOutbound(ctx context.Context, env endpoint.OutboundEnvelope, err error)
}

// NullObserver is an observer that does nothing.
//
// It can be embedded into observer implementations to avoid having to write
// empty methods.
type NullObserver struct{}

// BeforeInbound is called before a message is passed to the next pipeline stage.
func (NullObserver) BeforeInbound(context.Context, endpoint.InboundEnvelope) {}

// AfterInbound is called after a message is accepted by the next pipeline stage.
// err is the error returned by the next pipeline stage, which may be nil.
func (NullObserver) AfterInbound(context.Context, endpoint.InboundEnvelope, error) {}

// BeforeOutbound is called before a message is passed to the next pipeline stage.
func (NullObserver) BeforeOutbound(context.Context, endpoint.OutboundEnvelope) {}

// AfterOutbound is called after a message is accepted by the next pipeline stage.
// err is the error returned by the next pipeline stage, which may be nil.
func (NullObserver) AfterOutbound(context.Context, endpoint.OutboundEnvelope, error) {}
