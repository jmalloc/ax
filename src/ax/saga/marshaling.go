package saga

import (
	"fmt"

	"github.com/jmalloc/ax/src/ax/marshaling"
)

// MarshalInstance marshals i to a binary representation.
func MarshalInstance(i Instance) (contentType string, data []byte, err error) {
	return marshaling.MarshalJSON(i)
}

// UnmarshalInstance unmarshals a saga instance from some serialized
// representation. ct is the MIME content-type for the binary data.
func UnmarshalInstance(ct string, data []byte) (Instance, error) {
	v, err := marshaling.Unmarshal(ct, data)
	if err != nil {
		return nil, err
	}

	if m, ok := v.(Instance); ok {
		return m, nil
	}

	return nil, fmt.Errorf(
		"can not unmarshal '%s', content-type is not a saga instance",
		ct,
	)
}
