package bus_test

import (
	"context"

	. "github.com/jmalloc/ax/src/ax/bus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageBuffer", func() {
	b := &MessageBuffer{}

	Describe("SendMessage", func() {
		It("adds the message to the internal slice", func() {
			m := OutboundEnvelope{}
			m.MessageID.GenerateUUID()

			err := b.SendMessage(context.Background(), m)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(b.Messages).To(ConsistOf(m))
		})
	})
})
