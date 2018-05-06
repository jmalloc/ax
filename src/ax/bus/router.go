package bus

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/jmalloc/ax/src/ax"
)

// Router is an outbound pipeline stage that choses a destination endpoint for
// unicast messages based on the message type.
type Router struct {
	// Routes is the table used to determine the destination endpoint.
	Routes RoutingTable

	// Next is the next stage in the pipeline.
	Next OutboundPipeline

	// cache is a map of message type to endpoint, used to bypass scanning the
	// routing rules once a route has already been established.
	cache sync.Map // map[string]string
}

// Initialize is called when the transport is initialized.
func (r *Router) Initialize(ctx context.Context, t Transport) error {
	return r.Next.Initialize(ctx, t)
}

// SendMessage populates the m.DestinationEndpoint field of unicast messages that
// do not already have a DestinationEndpoint specified.
func (r *Router) SendMessage(ctx context.Context, m OutboundEnvelope) error {
	if err := r.ensureDestination(&m); err != nil {
		return err
	}

	return r.Next.SendMessage(ctx, m)
}

// ensureDestintion ensures that m.DestinationEndpoint is set if required.
func (r *Router) ensureDestination(m *OutboundEnvelope) error {
	if m.Operation != OpSendUnicast || m.DestinationEndpoint != "" {
		return nil
	}

	mt := ax.TypeOf(m.Message)

	if ep, ok := r.cache.Load(mt.Name); ok {
		m.DestinationEndpoint = ep.(string)
		return nil
	}

	if ep, ok := r.lookupDestination(mt); ok {
		r.cache.Store(mt.Name, ep)
		m.DestinationEndpoint = ep
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

// RoutingTable is set of rules that determine the destination endpoint for
// unicast messages.
type RoutingTable []route

// NewRoutingTable returns a routing table that choses a destination endpoint
// based on a mapping of protocol buffers message name to endpoint name.
//
// r is a sequence of key/value pairs that form the routing rules. The keys may
// be a fully-qualified protocol buffers message name, or a protocol buffers
// package name. The corresponding value is the endpoint name that such messages
// are routed to.
func NewRoutingTable(r ...string) (RoutingTable, error) {
	size := len(r)

	if size%2 != 0 {
		return nil, errors.New("r must contain an even number of strings")
	}

	t := make(RoutingTable, 0, size/2)

	for i := 0; i < size; i += 2 {
		t = append(
			t,
			route{r[i+0], r[i+1]},
		)
	}

	sort.Slice(t, func(i, j int) bool {
		return len(t[i].prefix) > len(t[j].prefix)
	})

	return t, nil
}

// Lookup returns the endpoint that messages of type mt are routed to.
//
// Routes are chosen based on the longest match, hence an exact match to the
// message type is favored, then the longest matching package name, and finally
// the empty string, which becomes the default route.
//
// If there is no default route and mt is a top-level message, that is, a
// message with no package name, then ok is false, indicating no route could be
// found.
func (t RoutingTable) Lookup(mt ax.MessageType) (ep string, ok bool) {
	for _, r := range t {
		if r.IsMatch(mt) {
			return r.endpoint, true
		}
	}

	return "", false
}

// route is a single routing rule that matches message types based on a prefix.
type route struct {
	prefix   string
	endpoint string
}

// IsMatch returns true if mt should be routed to r.endpoint.
func (r route) IsMatch(mt ax.MessageType) bool {
	return r.prefix == mt.Name ||
		r.prefix == "" ||
		strings.HasPrefix(mt.Name, r.prefix+".")
}
