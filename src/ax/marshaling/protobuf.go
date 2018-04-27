package marshaling

import (
	"errors"
	"fmt"
	"mime"
	"reflect"

	"github.com/golang/protobuf/proto"
)

const (
	// ProtobufContentType is the content-type name for protocol buffer messages.
	ProtobufContentType = "application/vnd.google.protobuf"
)

// MarshalProtobuf marshals a ProtocolBuffers messages to its binary
// representation and returns a MIME content-type that identifies the particular
// message protocol.
func MarshalProtobuf(m proto.Message) (ct string, data []byte, err error) {
	n := proto.MessageName(m)
	if n == "" {
		err = fmt.Errorf(
			"can not marshal '%s', protocol is not registered",
			reflect.TypeOf(m),
		)

		return
	}

	ct = mime.FormatMediaType(
		ProtobufContentType,
		map[string]string{"proto": n},
	)

	data, err = proto.Marshal(m)
	return
}

// UnmarshalProtobuf unmarshals a ProtocolBuffers messages from its binary
// representation using an unparsed MIME content-type to identify the particlar
// message protocol.
func UnmarshalProtobuf(ct string, data []byte) (proto.Message, error) {
	ctn, p, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	return UnmarshalProtobufParams(ctn, p, data)
}

// UnmarshalProtobufParams unmarshals a ProtocolBuffers messages from its binary
// representation using a pre-parsed MIME content-type to identify the particlar
// message protocol.
//
// ctn is the MIME content-type name, p is the set of pre-parsed content-type
// parameters, as returned by mime.ParseMediaType().
func UnmarshalProtobufParams(ctn string, p map[string]string, data []byte) (proto.Message, error) {
	if ctn != ProtobufContentType {
		return nil, fmt.Errorf(
			"can not unmarshal '%s' using protocol buffers, expected content-type is '%s'",
			ctn,
			ProtobufContentType,
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

	m := reflect.New(
		t.Elem(),
	).Interface().(proto.Message)

	return m, proto.Unmarshal(data, m)
}
