package keyset

import (
	"context"
	"strings"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/internal/reflectx"
)

// ByField returns a mapper that maps messages to instances using a set of
// fields within the message.
//
// IDs for new saga instances are produced by generating a random UUID.
//
// All fields must be strings. If any of the fields are empty, the message is
// not routed to any instance.
//
// The saga data must contain fields of the same name. If any of the data fields
// are empty, an empty keyset is produced, which results in a keyset validation
// error while handling the message.
func ByField(r Repository, f ...string) *Mapper {
	return &Mapper{
		r,
		byFieldResolver{f},
	}
}

type byFieldResolver struct {
	fields []string
}

func (r byFieldResolver) GenerateInstanceID(
	context.Context,
	ax.Envelope,
) (saga.InstanceID, error) {
	var id saga.InstanceID
	id.GenerateUUID()
	return id, nil
}

func (r byFieldResolver) MappingKeyForMessage(
	_ context.Context,
	env ax.Envelope,
) (string, bool, error) {
	k, ok := r.buildKey(env.Message)
	return k, ok, nil
}

func (r byFieldResolver) MappingKeysForInstance(
	_ context.Context,
	i saga.Instance,
) ([]string, error) {
	if k, ok := r.buildKey(i.Data); ok {
		return []string{k}, nil
	}

	// returning an empty key set will cause a key-set validation error
	return nil, nil
}

func (r byFieldResolver) buildKey(v interface{}) (string, bool) {
	f, err := reflectx.StringFields(v, r.fields)
	if err != nil {
		panic(err)
	}

	for _, v := range f {
		if v == "" {
			return "", false
		}
	}

	return strings.Join(f, "."), true
}
