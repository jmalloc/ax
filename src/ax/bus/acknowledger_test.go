package bus_test

import (
	"context"
	"errors"

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
	})

	Describe("DeliverMessage", func() {
		It("calls the next pipeline", func() {
			next.DeliverMessageFunc = func(ctx context.Context, _ MessageSender, _ InboundEnvelope) error {
				return nil
			}

			inEnv := InboundEnvelope{
				Done: func(context.Context, InboundOperation) error { return nil },
			}

			ack.DeliverMessage(context.Background(), nil /* sender */, inEnv)

			Expect(next.DeliverMessageCalls()).To(HaveLen(1))
		})

		It("marks message as acknowledged if no error occured", func() {
			next.DeliverMessageFunc = func(ctx context.Context, _ MessageSender, _ InboundEnvelope) error {
				return nil
			}

			inEnv := InboundEnvelope{
				Done: func(ctx context.Context, op InboundOperation) error {
					Expect(op).To(Equal(OpAck))
					return nil
				},
			}

			ack.DeliverMessage(context.Background(), nil /* sender */, inEnv)

			Expect(next.DeliverMessageCalls()).To(HaveLen(1))
		})

		It("marks message for retry if an error occured and retry policy approves retry", func() {
			next.DeliverMessageFunc = func(ctx context.Context, _ MessageSender, _ InboundEnvelope) error {
				return errors.New("something went wrong")
			}

			inEnv := InboundEnvelope{
				Done: func(ctx context.Context, op InboundOperation) error {
					Expect(op).To(Equal(OpRetry))
					return nil
				},
			}

			ack.RetryPolicy = func(InboundEnvelope) bool { return true }

			ack.DeliverMessage(context.Background(), nil /* sender */, inEnv)

			Expect(next.DeliverMessageCalls()).To(HaveLen(1))
		})

		It("marks message as rejected if an error occured and retry policy denies retry", func() {
			next.DeliverMessageFunc = func(ctx context.Context, _ MessageSender, _ InboundEnvelope) error {
				return errors.New("something went wrong")
			}

			inEnv := InboundEnvelope{
				Done: func(ctx context.Context, op InboundOperation) error {
					Expect(op).To(Equal(OpReject))
					return nil
				},
			}

			ack.RetryPolicy = func(InboundEnvelope) bool { return false }

			ack.DeliverMessage(context.Background(), nil /* sender */, inEnv)

			Expect(next.DeliverMessageCalls()).To(HaveLen(1))
		})

		It("marks message for retry if an error occured and default retry policy approves retry", func() {
			next.DeliverMessageFunc = func(ctx context.Context, _ MessageSender, _ InboundEnvelope) error {
				return errors.New("something went wrong")
			}

			inEnv := InboundEnvelope{
				DeliveryCount: 1,
				Done: func(ctx context.Context, op InboundOperation) error {
					Expect(op).To(Equal(OpRetry))
					return nil
				},
			}

			ack.RetryPolicy = DefaultRetryPolicy

			ack.DeliverMessage(context.Background(), nil /* sender */, inEnv)

			Expect(next.DeliverMessageCalls()).To(HaveLen(1))
		})
	})
})
