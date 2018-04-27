package marshaling_test

import (
	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax/internal/messagetest"
	. "github.com/jmalloc/ax/src/ax/marshaling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MarshalProtobuf", func() {
	message := &messagetest.NonAxMessage{
		Value: "<value>",
	}

	It("marshals the message using protocol buffers", func() {
		_, data, err := MarshalProtobuf(message)
		Expect(err).ShouldNot(HaveOccurred())

		var m messagetest.NonAxMessage
		err = proto.Unmarshal(data, &m)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(proto.Equal(&m, message)).To(BeTrue())
	})

	It("includes the protocol information in the content-type", func() {
		ct, _, err := MarshalProtobuf(message)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ct).To(Equal("application/vnd.google.protobuf; proto=ax.internal.messagetest.NonAxMessage"))
	})

	It("returns an error if the protocol name is not available", func() {
		var m proto.Message // no concrete value

		_, _, err := MarshalProtobuf(m)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("UnmarshalProtobuf", func() {
	message := &messagetest.NonAxMessage{
		Value: "<value>",
	}
	data, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	It("unmarshals the message using the protocol specified in the content-type", func() {
		m, err := UnmarshalProtobuf(
			"application/vnd.google.protobuf; proto=ax.internal.messagetest.NonAxMessage",
			data,
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(proto.Equal(m, message)).To(BeTrue())
	})

	It("returns an error if the content-type is invalid", func() {
		_, err := UnmarshalProtobuf("", data)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("UnmarshalProtobufParams", func() {
	message := &messagetest.NonAxMessage{
		Value: "<value>",
	}
	data, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	It("unmarshals the message using the protocol specified in the content-type parameters", func() {
		p := map[string]string{
			"proto": "ax.internal.messagetest.NonAxMessage",
		}
		m, err := UnmarshalProtobufParams("application/vnd.google.protobuf", p, data)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(proto.Equal(m, message)).To(BeTrue())
	})

	It("returns an error if the content-type is not the protocol buffers content type", func() {
		_, err := UnmarshalProtobufParams("application/json", nil, data)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the content-type parameters do not contain the protocol name", func() {
		_, err := UnmarshalProtobufParams("application/vnd.google.protobuf", nil, data)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the protocol name is not registered", func() {
		p := map[string]string{
			"proto": "ax.internal.messagetest.Unknown",
		}
		_, err := UnmarshalProtobufParams("application/vnd.google.protobuf", p, data)
		Expect(err).Should(HaveOccurred())
	})
})
