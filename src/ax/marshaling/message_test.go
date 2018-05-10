package marshaling_test

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	. "github.com/jmalloc/ax/src/ax/marshaling"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("MarshalMessage", func() {
	message := &messagetest.Message{
		Value: "<value>",
	}

	It("marshals the message using protocol buffers", func() {
		_, data, err := MarshalMessage(message)
		Expect(err).ShouldNot(HaveOccurred())

		var m messagetest.Message
		err = proto.Unmarshal(data, &m)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(proto.Equal(&m, message)).To(BeTrue())
	})

	It("includes the protocol information in the content-type", func() {
		ct, _, err := MarshalMessage(message)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ct).To(Equal("application/vnd.google.protobuf; proto=ax.internal.messagetest.Message"))
	})
})

var _ = Describe("UnmarshalMessage", func() {

	message := &messagetest.Message{
		Value: "<value>",
	}

	nonAxMessage := &messagetest.NonAxMessage{}

	axMessagePB, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	jsonMarshaler := jsonpb.Marshaler{
		EmitDefaults: false,
		EnumsAsInts:  false,
		Indent:       "  ",
		OrigName:     false,
	}

	buf := new(bytes.Buffer)

	err = jsonMarshaler.Marshal(buf, message)
	if err != nil {
		panic(err)
	}

	axMessageJSON := buf.Bytes()

	nonAxMessagePB, err := proto.Marshal(nonAxMessage)
	if err != nil {
		panic(err)
	}

	DescribeTable(
		"unmarshals the message using the protocol specified in the content-type",
		func(ct string, data []byte, expected *messagetest.Message) {
			m, err := UnmarshalMessage(
				ct,
				data,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(proto.Equal(m, message)).To(BeTrue())
		},
		Entry(
			"protobuf",
			"application/vnd.google.protobuf; proto=ax.internal.messagetest.Message",
			axMessagePB,
			message,
		),
		Entry(
			"JSON",
			"application/json; proto=ax.internal.messagetest.Message",
			axMessageJSON,
			message,
		),
	)

	DescribeTable(
		"execution of UnmarshalMessage errors out",
		func(ct string, data []byte) {
			_, err := UnmarshalMessage(ct, data)
			Expect(err).Should(HaveOccurred())
		},
		Entry(
			"returns an error if the content-type is invalid",
			"",
			axMessagePB,
		),
		Entry(
			"returns an error if the content-type is not supported",
			"application/x-unknown",
			axMessagePB,
		),
		Entry(
			"returns an error if an error occurs unmarshaling the protocol buffers message",
			"application/vnd.google.protobuf; proto=ax.internal.messagetest.Unknown", // note unknown message type
			axMessagePB,
		),
		Entry(
			"returns an error if an error occurs unmarshaling the JSON message",
			"application/json; proto=ax.internal.messagetest.Unknown", // note unknown message type
			axMessageJSON,
		),
		Entry(
			"returns an error if the buffer contains a protocol buffer message that is not an ax.Message",
			"application/vnd.google.protobuf; proto=ax.internal.messagetest.NonAxMessage",
			nonAxMessagePB,
		),
	)
})
