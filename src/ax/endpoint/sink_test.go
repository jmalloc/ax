package endpoint_test

import (
	"context"

	. "github.com/jmalloc/ax/src/ax/endpoint"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BufferedSink", func() {
	b := &BufferedSink{}

	Describe("Accept", func() {
		It("adds the message to the internal slice", func() {
			env := OutboundEnvelope{}
			env.MessageID.GenerateUUID()

			err := b.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(b.Envelopes).To(ConsistOf(env))
		})
	})
})
