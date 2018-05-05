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
	JSONContentType = "application/json"
)

// MarshalJSON marshals a ProtocolBuffers messages to its JSON
// representation and returns a MIME content-type that identifies the particular
// message protocol.
func MarshalJSON(msg proto.Message) (ct string, data []byte, err error) {

	n := proto.MessageName(msg)
	if n == "" {
		err = fmt.Errorf(
			"can not marshal '%s', protocol is not registered",
			reflect.TypeOf(msg),
		)
		return
	}

	b := new(bytes.Buffer)

	m := jsonpb.Marshaler{
		EmitDefaults: false,
		EnumsAsInts:  false,
		Indent:       "  ",
		OrigName:     false,
	}

	ct = mime.FormatMediaType(
		JSONContentType,
		map[string]string{"proto": n},
	)

	err = m.Marshal(b, msg)
	data = b.Bytes()

	return
}

// UnmarshalJSON unmarshals a ProtocolBuffers message from its JSON
// representation using an unparsed MIME content-type to identify the particular
// message protocol.
func UnmarshalJSON(ct string, data []byte) (proto.Message, error) {

	ctn, p, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	return UnmarshalJSONParams(ctn, p, data)

}

// UnmarshalJSONParams unmarshals a JSON-encoded message using
// a pre-parsed MIME content-type to identify the particlar
// message protocol.
//
// ctn is the MIME content-type name, p is the set of pre-parsed content-type
// parameters, as returned by mime.ParseMediaType().
func UnmarshalJSONParams(ctn string, p map[string]string, data []byte) (proto.Message, error) {

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
		AllowUnknownFields: true,
	}

	m := reflect.New(
		t.Elem(),
	).Interface().(proto.Message)

	return m, um.Unmarshal(
		bytes.NewReader(data),
		m,
	)
}
