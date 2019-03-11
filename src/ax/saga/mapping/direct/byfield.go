package direct

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/ax/saga/mapping/internal/byfieldx"
)

// ByField returns a mapper that maps messages to instances using a set of
// fields within the message.
//
// All fields must be strings. If any of the fields are empty, the message is
// not routed to any instance.
func ByField(f ...string) *Mapper {
	if len(f) == 0 {
		panic("at least one field must be specified")
	}

	return &Mapper{
		byFieldResolver{f},
	}
}

type byFieldResolver struct {
	fields []string
}

func (r byFieldResolver) InstanceIDForMessage(env ax.Envelope) (saga.InstanceID, bool) {
	k, err := byfieldx.FieldsToKey(env.Message, r.fields)
	if err != nil {
		// specifying incorrect fields is a programmer error
		// TODO: this could be verified eagerly once
		// https://github.com/jmalloc/ax/issues/81 is done
		panic(err)
	}

	if k == "" {
		return saga.InstanceID{}, false
	}

	id := saga.MustParseInstanceID(k)

	return id, true
}

// WithPrefix returns a mapper that wraps around m and adds a prefix string to the
// instance ID returned by m's Resolver.
//
// If m's Resolver does not return a valid instance ID, no prefix is added.
func WithPrefix(prefix string, m *Mapper) *Mapper {
	return &Mapper{
		withPrefixResolver{
			pre: prefix,
			m:   m,
		},
	}
}

type withPrefixResolver struct {
	pre string
	m   *Mapper
}

func (r withPrefixResolver) InstanceIDForMessage(env ax.Envelope) (saga.InstanceID, bool) {
	id, ok := r.m.Resolver.InstanceIDForMessage(env)
	if !ok {
		return id, false
	}
	return saga.MustParseInstanceID(r.pre + "-" + id.Get()), true
}
