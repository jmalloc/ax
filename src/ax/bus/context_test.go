package bus_test

import (
	. "github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageContext", func() {
	var (
		sink *BufferedSink
		ctx  *MessageContext
	)

	BeforeEach(func() {
		sink = &BufferedSink{}
		ctx = &MessageContext{
			Sink: sink,
		}

		ctx.Envelope.MessageID.GenerateUUID()
	})

	Describe("MessageEnvelope", func() {
		It("returns the envelope containing the inbound message", func() {
			Expect(ctx.MessageEnvelope()).To(Equal(ctx.Envelope))
		})
	})

	Describe("ExecuteCommand", func() {
		It("sends a unicast message via the sink", func() {
			err := ctx.ExecuteCommand(&messagetest.Command{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sink.Envelopes).To(HaveLen(1))
			Expect(sink.Envelopes[0].Operation).To(Equal(OpSendUnicast))
			Expect(sink.Envelopes[0].Message).To(Equal(&messagetest.Command{}))
		})

		It("configures the outbound message as a child of the inbound message", func() {
			_ = ctx.ExecuteCommand(&messagetest.Command{})
			Expect(sink.Envelopes[0].CausationID).To(Equal(ctx.Envelope.MessageID))
		})
	})

	Describe("PublishEvent", func() {
		It("sends a multicast message via the sink", func() {
			err := ctx.PublishEvent(&messagetest.Event{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sink.Envelopes).To(HaveLen(1))
			Expect(sink.Envelopes[0].Operation).To(Equal(OpSendMulticast))
			Expect(sink.Envelopes[0].Message).To(Equal(&messagetest.Event{}))
		})

		It("configures the outbound message as a child of the inbound message", func() {
			_ = ctx.PublishEvent(&messagetest.Event{})
			Expect(sink.Envelopes[0].CausationID).To(Equal(ctx.Envelope.MessageID))
		})
	})
})
