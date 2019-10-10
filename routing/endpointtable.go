package routing

import (
	"errors"
	"sort"
	"strings"

	"github.com/jmalloc/ax"
)

// EndpointTable is set of rules that determine the destination endpoint for
// unicast messages.
type EndpointTable []route

// NewEndpointTable returns an endpoint table that choses a destination endpoint
// based on a mapping of protocol buffers message name to endpoint name.
//
// r is a sequence of key/value pairs that form the routing rules. The keys may
// be a fully-qualified protocol buffers message name, or a protocol buffers
// package name. The corresponding value is the endpoint name that such messages
// are routed to.
func NewEndpointTable(r ...string) (EndpointTable, error) {
	size := len(r)

	if size%2 != 0 {
		return nil, errors.New("r must contain an even number of strings")
	}

	t := make(EndpointTable, 0, size/2)

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
func (t EndpointTable) Lookup(mt ax.MessageType) (ep string, ok bool) {
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
