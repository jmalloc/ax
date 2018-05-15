package persistence

import (
	"context"

	"github.com/jmalloc/ax/src/ax/endpoint"
)

// Injector is an implementation of endpoint.InboundPipeline that injects a data
// store into the context.
type Injector struct {
	DataStore DataStore
	Next      endpoint.InboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (i *Injector) Initialize(ctx context.Context, ep *endpoint.Endpoint) error {
	ctx = WithDataStore(ctx, i.DataStore)
	return i.Next.Initialize(ctx, ep)
}

// Accept calls i.Next.Accept() with a context derived from ctx
// and containing i.DataStore.
func (i *Injector) Accept(
	ctx context.Context,
	s endpoint.MessageSink,
	env endpoint.InboundEnvelope,
) error {
	ctx = WithDataStore(ctx, i.DataStore)
	return i.Next.Accept(ctx, s, env)
}
