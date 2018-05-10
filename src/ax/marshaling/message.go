package marshaling

import (
	"fmt"
	"mime"

	"github.com/golang/protobuf/proto"

	"github.com/jmalloc/ax/src/ax"
)

// MarshalMessage marshals m to a binary representation.
func MarshalMessage(m ax.Message) (contentType string, data []byte, err error) {
	return MarshalProtobuf(m)
}

// UnmarshalMessage unmarshals a message from a binary representation.
// ct is the MIME content-type for the binary data.
func UnmarshalMessage(ct string, data []byte) (ax.Message, error) {

	ctn, p, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	var pb proto.Message

	switch ctn {
	case ProtobufContentType:
		pb, err = UnmarshalProtobufParams(ctn, p, data)
		if err != nil {
			return nil, err
		}
	case JSONContentType:
		pb, err = UnmarshalJSONParams(ctn, p, data)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(
			"can not unmarshal '%s', content-type is not supported",
			ct,
		)
	}

	if m, ok := pb.(ax.Message); ok {
		return m, nil
	}

	return nil, fmt.Errorf(
		"can not unmarshal '%s', content-type is not a message",
		ct,
	)
}
