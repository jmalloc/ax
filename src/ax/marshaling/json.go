package marshaling

import (
	"bytes"
	"errors"
	"fmt"
	"mime"
	"reflect"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

const (
	// JSONContentType is the content-type name for JSON-encoded messages.
	JSONContentType = "application/vnd+ax.message+json"
)

// MarshalJSON JSON-ensodes a `proto.Message`
func MarshalJSON(msg proto.Message) (string, []byte, error) {

	b := new(bytes.Buffer)

	m := jsonpb.Marshaler{
		EmitDefaults: false,
		EnumsAsInts:  false,
		Indent:       "  ",
		OrigName:     false,
	}

	err := m.Marshal(b, msg)
	if err != nil {
		return "", nil, err
	}

	return JSONContentType, b.Bytes(), nil
}

// UnmarshalJSON decodes JSON-encoded data into `proto.Message`
func UnmarshalJSON(ct string, data []byte) (proto.Message, error) {

	ctn, p, err := mime.ParseMediaType(ct)
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

	n, ok := p["proto"]
	if !ok {
		return nil, errors.New(
			"can not unmarshal message, protocol is not specified in content-type parameters",
		)
	}

	t := proto.MessageType(n)
	if t == nil {
		return nil, fmt.Errorf(
			"can not unmarshal '%s', protocol is not registered",
			n,
		)
	}

	um := jsonpb.Unmarshaler{
		AllowUnknownFields: false,
	}

	m := reflect.New(
		t.Elem(),
	).Interface().(proto.Message)

	if err := um.Unmarshal(
		bytes.NewReader(data),
		m,
	); err != nil {
		return nil, err
	}

	return m, nil
}
