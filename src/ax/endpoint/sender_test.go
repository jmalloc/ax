package endpoint_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	. "github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SinkSender", func() {
	var (
		sink   *BufferedSink
		sender SinkSender
	)

	BeforeEach(func() {
		sink = &BufferedSink{}
		sender = SinkSender{Sink: sink}
	})

	Describe("ExecuteCommand", func() {
		It("sends a unicast message to the sink", func() {
			err := sender.ExecuteCommand(context.Background(), &messagetest.Command{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(sink.Envelopes()).To(HaveLen(1))
			env := sink.Envelopes()[0]
			Expect(env.Operation).To(Equal(OpSendUnicast))
			Expect(env.Message).To(Equal(&messagetest.Command{}))
		})

		It("configures the outbound message as a child of the envelope in ctx", func() {
			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_ = sender.ExecuteCommand(ctx, &messagetest.Command{})

			Expect(sink.Envelopes()).To(HaveLen(1))
			Expect(sink.Envelopes()[0].CausationID).To(Equal(env.MessageID))
		})
	})

	Describe("PublishEvent", func() {
		It("sends a multicast message to the sink", func() {
			err := sender.PublishEvent(context.Background(), &messagetest.Event{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(sink.Envelopes()).To(HaveLen(1))
			env := sink.Envelopes()[0]
			Expect(env.Operation).To(Equal(OpSendMulticast))
			Expect(env.Message).To(Equal(&messagetest.Event{}))
		})

		It("configures the outbound message as a child of the envelope in ctx", func() {
			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_ = sender.PublishEvent(ctx, &messagetest.Event{})

			Expect(sink.Envelopes()).To(HaveLen(1))
			Expect(sink.Envelopes()[0].CausationID).To(Equal(env.MessageID))
		})
	})
})
