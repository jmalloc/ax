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

		It("returns a validation error if one of validators fails", func() {
			expected := errors.New("test validation error")
			validator2.ValidateFunc = func(ctx context.Context, msg ax.Message) error {
				return expected
			}

			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_, err := sender.ExecuteCommand(ctx, &messagetest.Command{})
			Expect(err).Should(MatchError(expected))
		})

		It("uses default message validator if none of validators are assigned", func() {
			sink = &BufferedSink{}
			sender = SinkSender{
				Sink: sink,
			}
			// execute a command with nil passed instead of valid message
			// to verify that default validator will emit an error
			_, err := sender.ExecuteCommand(context.Background(), nil)
			Expect(err).Should(HaveOccurred())

			// assert that error emitted is a validation error
			_, ok := err.(*ValidationError)
			Expect(ok).Should(BeTrue())
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

		It("returns a validation error if one of validators fails", func() {
			expected := errors.New("test validation error")
			validator2.ValidateFunc = func(ctx context.Context, msg ax.Message) error {
				return expected
			}

			env := ax.NewEnvelope(&messagetest.Message{})
			ctx := WithEnvelope(context.Background(), env)

			_, err := sender.PublishEvent(ctx, &messagetest.Event{})
			Expect(err).Should(MatchError(expected))
		})

		It("uses default message validator if none of validators are assigned", func() {
			sink = &BufferedSink{}
			sender = SinkSender{
				Sink: sink,
			}
			// publish an event with nil passed instead of valid message
			// to verify that default validator will emit an error
			_, err := sender.PublishEvent(context.Background(), nil)
			Expect(err).Should(HaveOccurred())

			// assert that error emitted is a validation error
			_, ok := err.(*ValidationError)
			Expect(ok).Should(BeTrue())
		})
	})
})
