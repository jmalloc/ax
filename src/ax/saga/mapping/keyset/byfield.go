package keyset

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/ax/saga/mapping/internal/byfieldx"
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
	if len(f) == 0 {
		panic("at least one field must be specified")
	}

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
	k, err := byfieldx.FieldsToKey(env.Message, r.fields)
	if err != nil {
		// specifying incorrect fields is a programmer error
		// TODO: this could be verified eagerly once
		// https://github.com/jmalloc/ax/issues/81 is done
		panic(err)
	}

	return k, k != "", nil
}

func (r byFieldResolver) MappingKeysForInstance(
	_ context.Context,
	i saga.Instance,
) ([]string, error) {
	k, err := byfieldx.FieldsToKey(i.Data, r.fields)
	if err != nil {
		// specifying incorrect fields is a programmer error
		// TODO: this could be verified eagerly once
		// https://github.com/jmalloc/ax/issues/81 is done
		panic(err)
	}

	if k != "" {
		return []string{k}, nil
	}

	// returning an empty key set will cause a key-set validation error
	return nil, nil
}
