package ax_test

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	. "github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/axtest/testmessages"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("MarshalMessage", func() {
	message := &testmessages.Message{
		Value: "<value>",
	}

	It("marshals the message using protocol buffers", func() {
		_, data, err := MarshalMessage(message)
		Expect(err).ShouldNot(HaveOccurred())

		var m testmessages.Message
		err = proto.Unmarshal(data, &m)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(proto.Equal(&m, message)).To(BeTrue())
	})

	It("includes the protocol information in the content-type", func() {
		ct, _, err := MarshalMessage(message)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ct).To(Equal("application/vnd.google.protobuf; proto=axtest.testmessages.Message"))
	})
})

var _ = Describe("UnmarshalMessage", func() {
	message := &testmessages.Message{
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
			m, err := UnmarshalMessage(ct, data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(proto.Equal(m, message)).To(BeTrue())
		},
		Entry(
			"protobuf",
			"application/vnd.google.protobuf; proto=axtest.testmessages.Message",
			messagePB,
		),
		Entry(
			"JSON",
			"application/json; proto=axtest.testmessages.Message",
			messageJSON.Bytes(),
		),
	)

	It("returns an error if an error occurs in the underlying unmarshaler", func() {
		_, err := UnmarshalMessage(
			"application/vnd.google.protobuf; proto=axtest.testmessages.Unknown", // note unknown message type
			messagePB,
		)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the content is not an ax.Message", func() {
		pb, err := proto.Marshal(&testmessages.NonAxMessage{})
		Expect(err).ShouldNot(HaveOccurred())

		_, err = UnmarshalMessage(
			"application/vnd.google.protobuf; proto=axtest.testmessages.NonAxMessage",
			pb,
		)

		Expect(err).Should(HaveOccurred())
	})
})
