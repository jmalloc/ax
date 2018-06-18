package direct

import (
	"strings"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/internal/reflectx"
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
	f, err := reflectx.StringFields(env.Message, r.fields)
	if err != nil {
		panic(err)
	}

	for _, v := range f {
		if v == "" {
			return saga.InstanceID{}, false
		}
	}

	var id saga.InstanceID
	id.MustParse(
		strings.Join(f, "."),
	)

	return id, true
}
