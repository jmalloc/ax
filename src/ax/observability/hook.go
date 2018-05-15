package observability

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jmalloc/ax/src/ax/bus"
)

// InboundHook is an inbound pipeline stage that invokes hook methods
// on a set of observers.
type InboundHook struct {
	Observers []Observer
	Next      bus.InboundPipeline

	before []BeforeInboundObserver
	after  []AfterInboundObserver
}

// Initialize is called after the transport is initialized. It can be used
// to inspect or configure the transport as per the needs of the pipeline.
func (o *InboundHook) Initialize(ctx context.Context, t bus.Transport) error {
	for _, v := range o.Observers {
		used := false

		if ob, ok := v.(BeforeInboundObserver); ok {
			used = true
			o.before = append(o.before, ob)
		}

		if ob, ok := v.(AfterInboundObserver); ok {
			used = true
			o.after = append(o.after, ob)
		}

		if !used {
			panic(fmt.Sprintf(
				"%s does not implement either of BeforeInboundObserver or AfterInboundObserver",
				reflect.TypeOf(v),
			))
		}
	}

	return o.Next.Initialize(ctx, t)
}

// Accept forwards an inbound message through the pipeline until
// it is handled by some application-defined message handler(s).
func (o *InboundHook) Accept(ctx context.Context, s bus.MessageSink, env bus.InboundEnvelope) error {
	for _, ob := range o.before {
		ob.BeforeInbound(ctx, env)
	}

	err := o.Next.Accept(ctx, s, env)

	for _, ob := range o.after {
		ob.AfterInbound(ctx, env, err)
	}

	return err
}

// OutboundHook is an outbound pipeline stage that invokes hook methods
// on a set of observers.
type OutboundHook struct {
	Observers []Observer
	Next      bus.OutboundPipeline

	before []BeforeOutboundObserver
	after  []AfterOutboundObserver
}

// Initialize is called after the transport is initialized. It can be used
// to inspect or configure the transport as per the needs of the pipeline.
func (o *OutboundHook) Initialize(ctx context.Context, t bus.Transport) error {
	for _, v := range o.Observers {
		used := false

		if ob, ok := v.(BeforeOutboundObserver); ok {
			used = true
			o.before = append(o.before, ob)
		}

		if ob, ok := v.(AfterOutboundObserver); ok {
			used = true
			o.after = append(o.after, ob)
		}

		if !used {
			panic(fmt.Sprintf(
				"%s does not implement either of BeforeOutboundObserver or AfterOutboundObserver",
				reflect.TypeOf(v),
			))
		}
	}

	return o.Next.Initialize(ctx, t)
}

// Accept processes the message encapsulated in env.
func (o *OutboundHook) Accept(ctx context.Context, env bus.OutboundEnvelope) error {
	for _, ob := range o.before {
		ob.BeforeOutbound(ctx, env)
	}

	err := o.Next.Accept(ctx, env)

	for _, ob := range o.after {
		ob.AfterOutbound(ctx, env, err)
	}

	return err
}
