package marshaling

import (
	"bytes"
	"fmt"
	"mime"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

const (
	// JSONContentType is the content-type name for JSON-encoded messages.
	JSONContentType = "application/vnd+ax.message+json"
)

// MarshalJSON JSON-ensodes a message
func MarshalJSON(msg proto.Message) ([]byte, error) {

	b := new(bytes.Buffer)

	m := jsonpb.Marshaler{
		EmitDefaults: false,
		EnumsAsInts:  false,
		Indent:       "  ",
		OrigName:     false,
	}

	err := m.Marshal(b, msg)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// UnmarshalJSON decodes JSON-encoded data into `proto.Message`
func UnmarshalJSON(ct string, data []byte) (proto.Message, error) {
	ctn, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	if ctn != JSONContentType {
		return nil, fmt.Errorf(
			"can not unmarshal '%s' using JSON, expected content-type is '%s'",
			ctn,
			JSONContentType,
		)
	}

	um := jsonpb.Unmarshaler{
		AllowUnknownFields: false,
	}

	var msg proto.Message

	if err := um.Unmarshal(bytes.NewReader(data), msg); err != nil {
		return nil, err
	}

	return msg, nil
}
