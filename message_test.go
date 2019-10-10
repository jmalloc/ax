package ax_test

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	. "github.com/jmalloc/ax"
	"github.com/jmalloc/ax/axtest/testmessages"
	"github.com/jmalloc/ax/ident"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("GenerateMessageID", func() {
	It("generates a unique identifier for a message", func() {
		id := GenerateMessageID()
		u, err := uuid.FromString(id.Get())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(u.String()).To(Equal(id.Get()))
	})
})

var _ = Describe("ParseMessageID", func() {
	It("sets the ID to the parsed value", func() {
		id, err := ParseMessageID("<id>")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(id.Get()).To(Equal("<id>"))
	})

	It("returns an error if the value is empty", func() {
		_, err := ParseMessageID("")

		Expect(err).To(Equal(ident.ErrEmptyID))
	})

	It("returns an error if the ID is already set", func() {
		id := MustParseMessageID("<id>")
		err := id.Parse("<id>")

		Expect(err).To(Equal(ident.ErrIDNotEmpty))
	})
})

var _ = Describe("MustParseMessageID", func() {
	It("sets the ID to the parsed value", func() {
		id := MustParseMessageID("<id>")

		Expect(id.Get()).To(Equal("<id>"))
	})

	It("panics if the value is empty", func() {
		Expect(func() {
			MustParseMessageID("")
		}).To(Panic())
	})

	It("panics if the ID is already set", func() {
		id := MustParseMessageID("<id>")

		Expect(func() {
			id.MustParse("<id>")
		}).To(Panic())
	})
})

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
