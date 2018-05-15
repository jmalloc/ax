package marshaling

import (
	"fmt"
	"mime"

	"github.com/golang/protobuf/proto"
)

// Unmarshal unmarshals a protocol buffers message from some serialized
// representation. ct is the MIME content-type for the binary data.
func Unmarshal(ct string, data []byte) (proto.Message, error) {
	ctn, p, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	switch ctn {
	case ProtobufContentType:
		return UnmarshalProtobufParams(ctn, p, data)
	case JSONContentType:
		return UnmarshalJSONParams(ctn, p, data)
	default:
		return nil, fmt.Errorf(
			"can not unmarshal '%s', content-type is not supported",
			ct,
		)
	}
}
