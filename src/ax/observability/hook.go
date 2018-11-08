package observability

import (
	"context"

	"github.com/jmalloc/ax/src/ax/endpoint"
)

// InboundHook is an inbound pipeline stage that invokes hook methods
// on a set of observers.
type InboundHook struct {
	Observers []Observer
	Next      endpoint.InboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (h *InboundHook) Initialize(ctx context.Context, ep *endpoint.Endpoint) error {
	return h.Next.Initialize(ctx, ep)
}

// Accept forwards an inbound message through the pipeline until
// it is handled by some application-defined message handler(s).
func (h *InboundHook) Accept(ctx context.Context, s endpoint.MessageSink, env endpoint.InboundEnvelope) error {
	for _, o := range h.Observers {
		o.BeforeInbound(ctx, env)
	}

	err := h.Next.Accept(ctx, s, env)

	for _, o := range h.Observers {
		o.AfterInbound(ctx, env, err)
	}

	return err
}

// OutboundHook is an outbound pipeline stage that invokes hook methods
// on a set of observers.
type OutboundHook struct {
	Observers []Observer
	Next      endpoint.OutboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (h *OutboundHook) Initialize(ctx context.Context, ep *endpoint.Endpoint) error {
	return h.Next.Initialize(ctx, ep)
}

// Accept processes the message encapsulated in env.
func (h *OutboundHook) Accept(ctx context.Context, env endpoint.OutboundEnvelope) error {
	for _, o := range h.Observers {
		o.BeforeOutbound(ctx, env)
	}

	err := h.Next.Accept(ctx, env)

	for _, o := range h.Observers {
		o.AfterOutbound(ctx, env, err)
	}

	return err
}
