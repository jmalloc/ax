package marshaling_test

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax/internal/messagetest"
	. "github.com/jmalloc/ax/src/ax/marshaling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MarshalJSON", func() {

	message := &messagetest.NonAxMessage{
		Value: "<value>",
	}

	It("marshals the message using JSON", func() {

		_, data, err := MarshalJSON(message)
		Expect(err).ShouldNot(HaveOccurred())

		var m messagetest.NonAxMessage
		err = json.Unmarshal(data, &m)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(m).Should(Equal(*message))
	})

	It("includes the protocol information in the content-type", func() {
		ct, _, err := MarshalJSON(message)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ct).To(Equal("application/vnd+ax.message+json; proto=ax.internal.messagetest.NonAxMessage"))
	})

	It("returns an error if the protocol name is not available", func() {
		var m proto.Message // no concrete value

		_, _, err := MarshalJSON(m)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("UnmarshalJSON", func() {

	message := &messagetest.NonAxMessage{
		Value: "<value>",
	}
	_, data, err := MarshalJSON(message)
	if err != nil {
		panic(err)
	}

	It("unmarshals the message using the JSON specified in the content-type", func() {

		m, err := UnmarshalJSON(
			fmt.Sprintf(
				"%s; proto=%s",
				JSONContentType,
				"ax.internal.messagetest.NonAxMessage",
			),
			data,
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(proto.Equal(m, message)).To(BeTrue())
	})

	It("returns an error if the content-type is invalid", func() {
		_, err := UnmarshalJSON("", data)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the content-type is not specific to message json encoding", func() {
		_, err := UnmarshalJSON("text/plain	", data)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the content-type does not specify protocol as a content type parameter", func() {
		_, err := UnmarshalJSON(JSONContentType, data)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if message type is unregistered", func() {
		_, err := UnmarshalJSON(
			fmt.Sprintf(
				"%s; proto=%s",
				JSONContentType,
				"ax.internal.messagetest.NonExistingType",
			),
			data,
		)
		Expect(err).Should(HaveOccurred())
	})
})
