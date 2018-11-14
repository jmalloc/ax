package routing

import (
	"context"
	"fmt"
	"sync"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	opentracing "github.com/opentracing/opentracing-go"
)

// Router is an outbound pipeline stage that choses a destination endpoint for
// unicast messages based on the message type.
type Router struct {
	// Routes is the table used to determine the destination endpoint.
	Routes EndpointTable

	// Next is the next stage in the pipeline.
	Next endpoint.OutboundPipeline

	// cache is a map of message type to endpoint, used to bypass scanning the
	// routing rules once a route has already been established.
	cache sync.Map // map[string]string
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (r *Router) Initialize(ctx context.Context, ep *endpoint.Endpoint) error {
	return r.Next.Initialize(ctx, ep)
}

// Accept populates the evn.DestinationEndpoint field of unicast messages that
// do not already have a DestinationEndpoint specified.
func (r *Router) Accept(ctx context.Context, env endpoint.OutboundEnvelope) error {
	span := opentracing.SpanFromContext(ctx)

	if env.Operation == endpoint.OpSendUnicast {
		if env.DestinationEndpoint == "" {
			if err := r.ensureDestination(&env); err != nil {
				return err
			}

			traceSelect(span, env.DestinationEndpoint)
		} else {
			tracePreserve(span, env.DestinationEndpoint)
		}
	} else {
		traceForward(span)
	}

	return r.Next.Accept(ctx, env)
}

// ensureDestintion ensures that env.DestinationEndpoint is set if required.
func (r *Router) ensureDestination(env *endpoint.OutboundEnvelope) error {
	mt := env.Type()

	if ep, ok := r.cache.Load(mt.Name); ok {
		env.DestinationEndpoint = ep.(string)
		return nil
	}

	if ep, ok := r.lookupDestination(mt); ok {
		r.cache.Store(mt.Name, ep)
		env.DestinationEndpoint = ep
		return nil
	}

	return fmt.Errorf(
		"no endpoint route is configured for outbound '%s' message",
		mt.Name,
	)
}

// lookupDestination returns the destination endpoint for mt, using the
// routing table if available, or otherwise routing to the endpoint with the
// same name as the message's protocol buffers package name.
func (r *Router) lookupDestination(mt ax.MessageType) (string, bool) {
	if ep, ok := r.Routes.Lookup(mt); ok {
		return ep, ok
	}

	if ep := mt.PackageName(); ep != "" {
		return ep, true
	}

	return "", false
}
