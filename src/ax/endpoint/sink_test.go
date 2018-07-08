package endpoint_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax"

	. "github.com/jmalloc/ax/src/ax/endpoint"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BufferedSink", func() {
	env := OutboundEnvelope{
		Envelope: ax.Envelope{
			MessageID: ax.GenerateMessageID(),
		},
	}

	sink := &BufferedSink{}

	Describe("Accept", func() {
		It("adds the message to the internal slice", func() {
			err := sink.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sink.Envelopes()).To(ConsistOf(env))
		})
	})

	Describe("Reset", func() {
		It("removes the buffered envelopes", func() {
			sink.Accept(context.Background(), env)
			sink.Reset()

			Expect(sink.Envelopes()).To(BeEmpty())
		})
	})

	Describe("TakeEnvelopes", func() {
		It("returns the buffered envelopes", func() {
			sink.Accept(context.Background(), env)
			Expect(sink.TakeEnvelopes()).To(ConsistOf(env))
		})

		It("clears the buffer", func() {
			sink.Accept(context.Background(), env)
			sink.TakeEnvelopes()

			Expect(sink.Envelopes()).To(BeEmpty())
		})
	})
})
