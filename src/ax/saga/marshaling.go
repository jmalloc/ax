package saga

import (
	"fmt"

	"github.com/jmalloc/ax/src/ax/marshaling"
)

// MarshalData marshals d to a binary representation.
func MarshalData(d Data) (contentType string, data []byte, err error) {
	return marshaling.MarshalJSON(d)
}

// UnmarshalData unmarshals a saga instance from some serialized
// representation. ct is the MIME content-type for the binary data.
func UnmarshalData(ct string, data []byte) (Data, error) {
	v, err := marshaling.Unmarshal(ct, data)
	if err != nil {
		return nil, err
	}

	if m, ok := v.(Data); ok {
		return m, nil
	}

	return nil, fmt.Errorf(
		"can not unmarshal '%s', content-type is not a saga instance",
		ct,
	)
}
