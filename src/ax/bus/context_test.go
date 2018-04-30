package bus_test

import (
	"context"

	. "github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/internal/bustest"
	"github.com/jmalloc/ax/src/ax/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageContext", func() {
	var (
		sender *bustest.MessageSenderMock
		ctx    *MessageContext
	)

	BeforeEach(func() {
		sender = &bustest.MessageSenderMock{
			SendMessageFunc: func(context.Context, OutboundEnvelope) error { return nil },
		}
		ctx = &MessageContext{
			Sender: sender,
		}

		ctx.Envelope.MessageID.GenerateUUID()
	})

	Describe("MessageEnvelope", func() {
		It("returns the envelope containing the inbound message", func() {
			Expect(ctx.MessageEnvelope()).To(Equal(ctx.Envelope))
		})
	})

	Describe("ExecuteCommand", func() {
		It("sends a unicast message via the sender", func() {
			err := ctx.ExecuteCommand(&messagetest.Command{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(sender.SendMessageCalls()).To(HaveLen(1))

			m := sender.SendMessageCalls()[0].M
			Expect(m.Operation).To(Equal(OpSendUnicast))
			Expect(m.Message).To(Equal(&messagetest.Command{}))
		})

		It("configures the outbound message as a child of the inbound message", func() {
			_ = ctx.ExecuteCommand(&messagetest.Command{})
			Expect(sender.SendMessageCalls()[0].M.CausationID).To(Equal(ctx.Envelope.MessageID))
		})
	})

	Describe("PublishEvent", func() {
		It("sends a multicast message via the sender", func() {
			err := ctx.PublishEvent(&messagetest.Event{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(sender.SendMessageCalls()).To(HaveLen(1))

			m := sender.SendMessageCalls()[0].M
			Expect(m.Operation).To(Equal(OpSendMulticast))
			Expect(m.Message).To(Equal(&messagetest.Event{}))
		})

		It("configures the outbound message as a child of the inbound message", func() {
			_ = ctx.PublishEvent(&messagetest.Event{})
			Expect(sender.SendMessageCalls()[0].M.CausationID).To(Equal(ctx.Envelope.MessageID))
		})
	})
})
