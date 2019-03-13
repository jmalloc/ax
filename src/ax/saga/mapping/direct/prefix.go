package direct

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/saga"
)

// WithPrefix returns a mapper that wraps around m and adds a prefix string to the
// instance ID returned by m's Resolver.
//
// If m's Resolver does not return a valid instance ID, no prefix is added.
func WithPrefix(prefix string, m *Mapper) *Mapper {
	return &Mapper{
		withPrefixResolver{
			pre: prefix,
			wr:  m.Resolver,
		},
	}
}

type withPrefixResolver struct {
	pre string
	wr  Resolver
}

func (r withPrefixResolver) InstanceIDForMessage(env ax.Envelope) (saga.InstanceID, bool) {
	id, ok := r.wr.InstanceIDForMessage(env)
	if !ok {
		return id, false
	}
	return saga.MustParseInstanceID(r.pre + "-" + id.Get()), true
}
