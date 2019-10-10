package direct

import (
	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/saga"
	"github.com/jmalloc/ax/saga/mapping/internal/byfieldx"
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
