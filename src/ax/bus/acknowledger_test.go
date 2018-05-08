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

	Describe("Accept", func() {
		It("calls the next pipeline", func() {
			next.AcceptFunc = func(ctx context.Context, _ MessageSink, _ InboundEnvelope) error {
				return nil
			}

			inEnv := InboundEnvelope{
				Done: func(context.Context, InboundOperation) error { return nil },
			}

			ack.Accept(context.Background(), nil /* sender */, inEnv)

			Expect(next.AcceptCalls()).To(HaveLen(1))
		})

		Context("when there is a custom retry policy", func() {
			It("marks message as acknowledged if no error occurred", func() {
				next.AcceptFunc = func(ctx context.Context, _ MessageSink, _ InboundEnvelope) error {
					return nil
				}

				var doneOp InboundOperation
				inEnv := InboundEnvelope{
					Done: func(ctx context.Context, op InboundOperation) error {
						doneOp = op
						return nil
					},
				}

				ack.Accept(context.Background(), nil /* sender */, inEnv)

				Expect(next.AcceptCalls()).To(HaveLen(1))
				Expect(doneOp).To(Equal(OpAck))
			})

			It("marks message for retry if an error occurred and retry policy approves retry", func() {
				next.AcceptFunc = func(ctx context.Context, _ MessageSink, _ InboundEnvelope) error {
					return errors.New("something went wrong")
				}

				var doneOp InboundOperation
				inEnv := InboundEnvelope{
					Done: func(ctx context.Context, op InboundOperation) error {
						doneOp = op
						return nil
					},
				}

				ack.RetryPolicy = func(InboundEnvelope) bool { return true }

				ack.Accept(context.Background(), nil /* sender */, inEnv)

				Expect(next.AcceptCalls()).To(HaveLen(1))
				Expect(doneOp).To(Equal(OpRetry))
			})

			It("marks message as rejected if an error occurred and retry policy denies retry", func() {
				next.AcceptFunc = func(ctx context.Context, _ MessageSink, _ InboundEnvelope) error {
					return errors.New("something went wrong")
				}

				var doneOp InboundOperation
				inEnv := InboundEnvelope{
					Done: func(ctx context.Context, op InboundOperation) error {
						doneOp = op
						return nil
					},
				}

				ack.RetryPolicy = func(InboundEnvelope) bool { return false }

				ack.Accept(context.Background(), nil /* sender */, inEnv)

				Expect(next.AcceptCalls()).To(HaveLen(1))
				Expect(doneOp).To(Equal(OpReject))
			})
		})

		Context("when using the default retry policy", func() {
			It("marks message for retry if an error occurred and delivery count is less than 3", func() {
				next.AcceptFunc = func(ctx context.Context, _ MessageSink, _ InboundEnvelope) error {
					return errors.New("something went wrong")
				}

				var doneOp InboundOperation
				inEnv := InboundEnvelope{
					DeliveryCount: 1,
					Done: func(ctx context.Context, op InboundOperation) error {
						doneOp = op
						return nil
					},
				}

				ack.RetryPolicy = DefaultRetryPolicy

				ack.Accept(context.Background(), nil /* sender */, inEnv)

				Expect(next.AcceptCalls()).To(HaveLen(1))
				Expect(doneOp).To(Equal(OpRetry))
			})

			It("marks message as rejected if an error occurred and delivery count is equal to 3", func() {
				next.AcceptFunc = func(ctx context.Context, _ MessageSink, _ InboundEnvelope) error {
					return errors.New("something went wrong")
				}

				var doneOp InboundOperation
				inEnv := InboundEnvelope{
					DeliveryCount: 3,
					Done: func(ctx context.Context, op InboundOperation) error {
						doneOp = op
						return nil
					},
				}

				ack.RetryPolicy = DefaultRetryPolicy

				ack.Accept(context.Background(), nil /* sender */, inEnv)

				Expect(next.AcceptCalls()).To(HaveLen(1))
				Expect(doneOp).To(Equal(OpReject))
			})
		})
	})
})
