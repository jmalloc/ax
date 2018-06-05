package endpoint_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
	. "github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/internal/endpointtest"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SinkSender", func() {
	var (
		sink                               *BufferedSink
		validator1, validator2, validator3 *endpointtest.ValidatorMock
		sender                             SinkSender
	)

	BeforeEach(func() {
		sink = &BufferedSink{}
		validator1 = &endpointtest.ValidatorMock{
			ValidateFunc: func(ctx context.Context, msg ax.Message) error {
				return nil
			},
		}
		validator2 = &endpointtest.ValidatorMock{
			ValidateFunc: func(ctx context.Context, msg ax.Message) error {
				return nil
			},
		}
		validator3 = &endpointtest.ValidatorMock{
			ValidateFunc: func(ctx context.Context, msg ax.Message) error {
				return nil
			},
		}
		sender = SinkSender{
			Sink: sink,
			Validators: []Validator{
				validator1,
				validator2,
				validator3,
			},
		}
	})

	Describe("ExecuteCommand", func() {
		It("sends a unicast message to the sink", func() {
			_, err := sender.ExecuteCommand(context.Background(), &messagetest.Command{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(sink.Envelopes()).To(HaveLen(1))
			env := sink.Envelopes()[0]
			Expect(env.Operation).To(Equal(OpSendUnicast))
			Expect(env.Message).To(Equal(&messagetest.Command{}))
		})

		It("configures the outbound message as a child of the envelope in ctx", func() {
			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_, _ = sender.ExecuteCommand(ctx, &messagetest.Command{})

			Expect(sink.Envelopes()).To(HaveLen(1))
			Expect(sink.Envelopes()[0].CausationID).To(Equal(env.MessageID))
		})

		It("returns the sent envelope", func() {
			env, _ := sender.ExecuteCommand(context.Background(), &messagetest.Command{})

			Expect(env).To(Equal(sink.Envelopes()[0].Envelope))
		})

		It("returns a validation error if one of the validators fails", func() {
			expected := errors.New("test validation error")
			validator2.ValidateFunc = func(ctx context.Context, msg ax.Message) error {
				return expected
			}

			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_, err := sender.ExecuteCommand(ctx, &messagetest.Command{})
			Expect(err).Should(MatchError(expected))
		})

		It("uses default message validators if the validators slice is empty", func() {
			sink = &BufferedSink{}
			sender = SinkSender{
				Sink: sink,
			}

			cmd := &endpointtest.SelfValidatingCommandMock{
				ValidateFunc: func() error {
					return nil
				},
			}
			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)
			_, err := sender.ExecuteCommand(ctx, cmd)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns an error if default command validation fails", func() {
			sink = &BufferedSink{}
			sender = SinkSender{
				Sink: sink,
			}

			expected := errors.New("test validation error")
			cmd := &endpointtest.SelfValidatingCommandMock{
				ValidateFunc: func() error {
					return expected
				},
			}
			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)
			_, err := sender.ExecuteCommand(ctx, cmd)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("PublishEvent", func() {
		It("sends a multicast message to the sink", func() {
			_, err := sender.PublishEvent(context.Background(), &messagetest.Event{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(sink.Envelopes()).To(HaveLen(1))
			env := sink.Envelopes()[0]
			Expect(env.Operation).To(Equal(OpSendMulticast))
			Expect(env.Message).To(Equal(&messagetest.Event{}))
		})

		It("configures the outbound message as a child of the envelope in ctx", func() {
			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_, _ = sender.PublishEvent(ctx, &messagetest.Event{})

			Expect(sink.Envelopes()).To(HaveLen(1))
			Expect(sink.Envelopes()[0].CausationID).To(Equal(env.MessageID))
		})

		It("returns the sent envelope", func() {
			env, _ := sender.PublishEvent(context.Background(), &messagetest.Event{})

			Expect(env).To(Equal(sink.Envelopes()[0].Envelope))
		})

		It("returns a validation error if one of the validators fails", func() {
			expected := errors.New("test validation error")
			validator2.ValidateFunc = func(ctx context.Context, msg ax.Message) error {
				return expected
			}

			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_, err := sender.PublishEvent(ctx, &messagetest.Event{})
			Expect(err).Should(MatchError(expected))
		})

		It("uses default message validators if the validators slice is empty", func() {
			sink = &BufferedSink{}
			sender = SinkSender{
				Sink: sink,
			}

			ev := &endpointtest.SelfValidatingEventMock{
				ValidateFunc: func() error {
					return nil
				},
			}
			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)
			_, err := sender.PublishEvent(ctx, ev)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns an error if default event validation fails", func() {
			sink = &BufferedSink{}
			sender = SinkSender{
				Sink: sink,
			}

			expected := errors.New("test validation error")
			ev := &endpointtest.SelfValidatingEventMock{
				ValidateFunc: func() error {
					return expected
				},
			}
			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)
			_, err := sender.PublishEvent(ctx, ev)
			Expect(err).Should(HaveOccurred())
		})
	})
})
