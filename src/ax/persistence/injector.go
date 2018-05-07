package persistence

import (
	"context"

	"github.com/jmalloc/ax/src/ax/bus"
)

// Injector is an implementation of bus.InboundPipeline that injects a data
// store into the context.
type Injector struct {
	DataStore DataStore
	Next      bus.InboundPipeline
}

// Initialize calls i.Next.Initialize() with a context derived from ctx and
// containing i.DataStore.
func (i *Injector) Initialize(ctx context.Context, t bus.Transport) error {
	ctx = WithDataStore(ctx, i.DataStore)
	return i.Next.Initialize(ctx, t)
}

// Accept calls i.Next.Accept() with a context derived from ctx
// and containing i.DataStore.
func (i *Injector) Accept(
	ctx context.Context,
	s bus.MessageSink,
	env bus.InboundEnvelope,
) error {
	ctx = WithDataStore(ctx, i.DataStore)
	return i.Next.Accept(ctx, s, env)
}
