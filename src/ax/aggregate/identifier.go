package aggregate

import (
	"fmt"
	"reflect"

	"github.com/jmalloc/ax/src/ax"
)

// Identifier is an interface for determining the aggregate ID that a command
// targets.
type Identifier interface {
	// AggregateID returns the ID of the aggregate that m targets.
	AggregateID(m ax.Command) (ID, error)
}

// ByFieldIdentifier is an Identifier instance that treats a specific field of
// the message as the aggregate ID.
type ByFieldIdentifier struct {
	FieldName string
}

// AggregateID returns the ID of the aggregate that m targets.
func (i *ByFieldIdentifier) AggregateID(m ax.Command) (ID, error) {
	mt := ax.TypeOf(m)

	f, ok := mt.StructType.FieldByName(i.FieldName)
	if !ok {
		return ID{}, fmt.Errorf(
			"%s does not contain a field named %s",
			reflect.TypeOf(m),
			i.FieldName,
		)
	}

	if f.Type.Kind() != reflect.String {
		return ID{}, fmt.Errorf(
			"%s.%s is not a string",
			reflect.TypeOf(m),
			i.FieldName,
		)
	}

	v := reflect.
		ValueOf(m).
		Elem().
		FieldByIndex(f.Index).
		Interface().(string)

	var id ID
	err := id.Parse(v)

	return id, err
}
