package endpoint_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
	. "github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/axtest/mocks"
	"github.com/jmalloc/ax/src/axtest/testmessages"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SinkSender", func() {
	var (
		sink                               *BufferedSink
		validator1, validator2, validator3 *mocks.ValidatorMock
		sender                             SinkSender
	)

	BeforeEach(func() {
		sink = &BufferedSink{}
		validator1 = &mocks.ValidatorMock{
			ValidateFunc: func(ctx context.Context, m ax.Message) error {
				return nil
			},
		}
		validator2 = &mocks.ValidatorMock{
			ValidateFunc: func(ctx context.Context, m ax.Message) error {
				return nil
			},
		}
		validator3 = &mocks.ValidatorMock{
			ValidateFunc: func(ctx context.Context, m ax.Message) error {
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
			_, err := sender.ExecuteCommand(context.Background(), &testmessages.Command{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(sink.Envelopes()).To(HaveLen(1))
			env := sink.Envelopes()[0]
			Expect(env.Operation).To(Equal(OpSendUnicast))
			Expect(env.Message).To(Equal(&testmessages.Command{}))
		})

		It("configures the outbound message as a child of the envelope in ctx", func() {
			env := ax.NewEnvelope(&testmessages.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_, _ = sender.ExecuteCommand(ctx, &testmessages.Command{})

			Expect(sink.Envelopes()).To(HaveLen(1))
			Expect(sink.Envelopes()[0].CausationID).To(Equal(env.MessageID))
		})

		It("returns the sent envelope", func() {
			env, _ := sender.ExecuteCommand(context.Background(), &testmessages.Command{})

			Expect(env).To(Equal(sink.Envelopes()[0].Envelope))
		})

		It("returns a validation error if one of the validators fails", func() {
			expected := errors.New("test validation error")
			validator2.ValidateFunc = func(ctx context.Context, m ax.Message) error {
				return expected
			}

			env := ax.NewEnvelope(&testmessages.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_, err := sender.ExecuteCommand(ctx, &testmessages.Command{})
			Expect(err).Should(MatchError(expected))
		})

		It("uses default message validators if the validators slice is empty", func() {
			sink = &BufferedSink{}
			sender = SinkSender{
				Sink: sink,
			}

			env := ax.NewEnvelope(&testmessages.Message{})
			ctx := WithEnvelope(context.Background(), env)
			_, err := sender.ExecuteCommand(
				ctx,
				&testmessages.SelfValidatingCommand{},
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns an error if default command validation fails", func() {
			sink = &BufferedSink{}
			sender = SinkSender{
				Sink: sink,
			}
			env := ax.NewEnvelope(&testmessages.Message{})
			ctx := WithEnvelope(context.Background(), env)
			_, err := sender.ExecuteCommand(
				ctx,
				&testmessages.FailedSelfValidatingCommand{},
			)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("PublishEvent", func() {
		It("sends a multicast message to the sink", func() {
			_, err := sender.PublishEvent(context.Background(), &testmessages.Event{})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(sink.Envelopes()).To(HaveLen(1))
			env := sink.Envelopes()[0]
			Expect(env.Operation).To(Equal(OpSendMulticast))
			Expect(env.Message).To(Equal(&testmessages.Event{}))
		})

		It("configures the outbound message as a child of the envelope in ctx", func() {
			env := ax.NewEnvelope(&testmessages.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_, _ = sender.PublishEvent(ctx, &testmessages.Event{})

			Expect(sink.Envelopes()).To(HaveLen(1))
			Expect(sink.Envelopes()[0].CausationID).To(Equal(env.MessageID))
		})

		It("returns the sent envelope", func() {
			env, _ := sender.PublishEvent(context.Background(), &testmessages.Event{})

			Expect(env).To(Equal(sink.Envelopes()[0].Envelope))
		})

		It("returns a validation error if one of the validators fails", func() {
			expected := errors.New("test validation error")
			validator2.ValidateFunc = func(ctx context.Context, m ax.Message) error {
				return expected
			}

			env := ax.NewEnvelope(&testmessages.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_, err := sender.PublishEvent(ctx, &testmessages.Event{})
			Expect(err).Should(MatchError(expected))
		})

		It("uses default message validators if the validators slice is empty", func() {
			sink = &BufferedSink{}
			sender = SinkSender{
				Sink: sink,
			}

			env := ax.NewEnvelope(&testmessages.Message{})
			ctx := WithEnvelope(context.Background(), env)
			_, err := sender.PublishEvent(
				ctx,
				&testmessages.SelfValidatingEvent{},
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns an error if default event validation fails", func() {
			sink = &BufferedSink{}
			sender = SinkSender{
				Sink: sink,
			}
			env := ax.NewEnvelope(&testmessages.Message{})
			ctx := WithEnvelope(context.Background(), env)
			_, err := sender.PublishEvent(
				ctx,
				&testmessages.FailedSelfValidatingEvent{},
			)
			Expect(err).Should(HaveOccurred())
		})
	})
})
