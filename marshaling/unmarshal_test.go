package marshaling_test

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/axtest/testmessages"
	. "github.com/jmalloc/ax/marshaling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unmarshal", func() {
	message := &testmessages.NonAxMessage{
		Value: "<value>",
	}

	messagePB, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	var messageJSON bytes.Buffer
	marshaller := jsonpb.Marshaler{}
	err = marshaller.Marshal(&messageJSON, message)
	if err != nil {
		panic(err)
	}

	DescribeTable(
		"unmarshals the message using the protocol specified in the content-type",
		func(ct string, data []byte) {
			m, err := Unmarshal(ct, data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(proto.Equal(m, message)).To(BeTrue())
		},
		Entry(
			"protobuf",
			"application/vnd.google.protobuf; proto=axtest.testmessages.NonAxMessage",
			messagePB,
		),
		Entry(
			"JSON",
			"application/json; proto=axtest.testmessages.NonAxMessage",
			messageJSON.Bytes(),
		),
	)

	It("returns an error if the content type is invalid", func() {
		_, err := Unmarshal("", messagePB)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the content type is not supported", func() {
		_, err := Unmarshal("application/x-unknown", messagePB)
		Expect(err).Should(HaveOccurred())
	})

	DescribeTable(
		"returns an error if an error occurs in the underlying unmarshaler",
		func(ct string, data []byte) {
			_, err := Unmarshal(ct, data)
			Expect(err).Should(HaveOccurred())
		},
		Entry(
			"returns an error if an error occurs unmarshaling the protocol buffers message",
			"application/vnd.google.protobuf; proto=axtest.testmessages.Unknown", // note unknown message type
			messagePB,
		),
		Entry(
			"returns an error if an error occurs unmarshaling the JSON message",
			"application/json; proto=axtest.testmessages.Unknown", // note unknown message type
			messageJSON.Bytes(),
		),
	)
})
