package persistence

import (
	"context"

	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/internal/reflectx"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// InboundInjector is an implementation of endpoint.InboundPipeline that injects
// a data store into the context.
type InboundInjector struct {
	DataStore DataStore
	Next      endpoint.InboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (i *InboundInjector) Initialize(ctx context.Context, ep *endpoint.Endpoint) error {
	ctx = WithDataStore(ctx, i.DataStore)
	return i.Next.Initialize(ctx, ep)
}

// Accept calls i.Next.Accept() with a context derived from ctx
// and containing i.DataStore.
func (i *InboundInjector) Accept(
	ctx context.Context,
	s endpoint.MessageSink,
	env endpoint.InboundEnvelope,
) error {
	ctx = WithDataStore(ctx, i.DataStore)

	traceInject(
		opentracing.SpanFromContext(ctx),
		i.DataStore,
	)

	return i.Next.Accept(ctx, s, env)
}

// OutboundInjector is an implementation of endpoint.OutboundPipeline that injects
// a data store into the context.
type OutboundInjector struct {
	DataStore DataStore
	Next      endpoint.OutboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (i *OutboundInjector) Initialize(ctx context.Context, ep *endpoint.Endpoint) error {
	ctx = WithDataStore(ctx, i.DataStore)
	return i.Next.Initialize(ctx, ep)
}

// Accept calls i.Next.Accept() with a context derived from ctx
// and containing i.DataStore.
func (i *OutboundInjector) Accept(
	ctx context.Context,
	env endpoint.OutboundEnvelope,
) error {
	ctx = WithDataStore(ctx, i.DataStore)

	traceInject(
		opentracing.SpanFromContext(ctx),
		i.DataStore,
	)

	return i.Next.Accept(ctx, env)
}

func traceInject(span opentracing.Span, ds DataStore) {
	if span == nil {
		return
	}

	span.LogFields(
		log.String("event", "persistence.inject"),
		log.String("message", "added data-store to context, forwarding message to the next pipeline stage"),
		log.String("data-store", reflectx.PrettyTypeName(ds)),
	)
}
