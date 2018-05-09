package observability

import (
	"context"

	"github.com/jmalloc/ax/src/ax/bus"
)

// InboundHook is an inbound pipeline stage that invokes hook methods
// on a set of observers.
type InboundHook struct {
	Observers []interface{}
	Next      bus.InboundPipeline
}

// Initialize is called after the transport is initialized. It can be used
// to inspect or configure the transport as per the needs of the pipeline.
func (o *InboundHook) Initialize(ctx context.Context, t bus.Transport) error {
	return o.Next.Initialize(ctx, t)
}

// Accept forwards an inbound message through the pipeline until
// it is handled by some application-defined message handler(s).
func (o *InboundHook) Accept(ctx context.Context, s bus.MessageSink, env bus.InboundEnvelope) error {
	var err error

	for _, v := range o.Observers {
		if ob, ok := v.(BeforeInboundObserver); ok {
			ctx, err = ob.BeforeInbound(ctx, env.Envelope)
			if err != nil {
				return err
			}
		}
	}

	acceptErr := o.Next.Accept(ctx, s, env)

	for _, v := range o.Observers {
		if ob, ok := v.(AfterInboundObserver); ok {
			err = ob.AfterInbound(ctx, env.Envelope, acceptErr)
			if err != nil {
				return err
			}
		}
	}

	return acceptErr
}

// OutboundHook is an outbound pipeline stage that invokes hook methods
// on a set of observers.
type OutboundHook struct {
	Observers []interface{}
	Next      bus.OutboundPipeline
}

// Initialize is called after the transport is initialized. It can be used
// to inspect or configure the transport as per the needs of the pipeline.
func (o *OutboundHook) Initialize(ctx context.Context, t bus.Transport) error {
	return o.Next.Initialize(ctx, t)
}

// Accept processes the message encapsulated in env.
func (o *OutboundHook) Accept(ctx context.Context, env bus.OutboundEnvelope) error {
	var err error

	for _, v := range o.Observers {
		if ob, ok := v.(BeforeOutboundObserver); ok {
			ctx, err = ob.BeforeOutbound(ctx, env.Envelope)
			if err != nil {
				return err
			}
		}
	}

	acceptErr := o.Next.Accept(ctx, env)

	for _, v := range o.Observers {
		if ob, ok := v.(AfterOutboundObserver); ok {
			err = ob.AfterOutbound(ctx, env.Envelope, acceptErr)
			if err != nil {
				return err
			}
		}
	}

	return acceptErr
}
