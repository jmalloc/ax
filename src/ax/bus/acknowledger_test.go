package bus_test

import (
	"context"

	. "github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/internal/bustest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Acknowledger", func() {
	var (
		next *bustest.InboundPipelineMock
		ack  *Acknowledger
	)

	BeforeEach(func() {
		next = &bustest.InboundPipelineMock{}
		ack = &Acknowledger{
			Next: next,
		}
	})

	Describe("Initialize", func() {
		It("calls the next pipeline", func() {
			next.InitializeFunc = func(ctx context.Context, _ Transport) error {
				return nil
			}

			ack.Initialize(context.Background(), nil /* transport */)

			Expect(next.InitializeCalls()).To(HaveLen(1))
			Expect(ack.RetryPolicy).To(Not(BeNil()))
		})

		It("sets the default retry policy if none is set", func() {
			// TODO
		})

		It("does not set the retry policy if one is already set", func() {
			// TODO
		})
	})

	Describe("DeliverMessage", func() {
		It("calls the next pipeline", func() {
			next.DeliverMessageFunc = func(ctx context.Context, _ MessageSender, _ InboundEnvelope) error {
				return nil
			}

			ack.DeliverMessage(context.Background(), nil /* sender */, InboundEnvelope{})

			Expect(next.DeliverMessageCalls()).To(HaveLen(1))
		})

		It("marks message as acknowledged if no error occured", func() {
			// TODO

			// next.DeliverMessageFunc = func(ctx context.Context, _ MessageSender, _ InboundEnvelope) error {
			// 	return nil
			// }
			//
			// ack.DeliverMessage(context.Background(), nil /* sender */, InboundEnvelope{})
			//
			// Expect(next.DeliverMessageCalls()).To(HaveLen(1))
		})

		It("marks message for retry if an error occured and retry policy approves retry", func() {
			// TODO
		})

		It("marks message as rejected if an error occured and retry policy denies retry", func() {
			// TODO
		})
	})
})
