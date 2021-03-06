package routing_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/axtest/mocks"
	"github.com/jmalloc/ax/axtest/testmessages"
	"github.com/jmalloc/ax/endpoint"
	. "github.com/jmalloc/ax/routing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ endpoint.InboundPipeline = (*Dispatcher)(nil) // ensure Dispatcher implements InboundPipeline

var _ = Describe("Dispatcher", func() {
	var (
		h1, h2, h3 *mocks.MessageHandlerMock
		sink       *endpoint.BufferedSink
		dispatcher *Dispatcher
	)

	BeforeEach(func() {
		noOp := func(context.Context, ax.Sender, ax.MessageContext) error { return nil }

		h1 = &mocks.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(
					&testmessages.Command{},
					&testmessages.Message{},
				)
			},
			HandleMessageFunc: noOp,
		}
		h2 = &mocks.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(
					&testmessages.Message{},
					&testmessages.Event{},
				)
			},
			HandleMessageFunc: noOp,
		}
		h3 = &mocks.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(
					&testmessages.Event{},
				)
			},
			HandleMessageFunc: noOp,
		}

		t, err := NewHandlerTable(h1, h2, h3)
		Expect(err).ShouldNot(HaveOccurred())

		sink = &endpoint.BufferedSink{}
		dispatcher = &Dispatcher{Routes: t}
	})

	Describe("Initialize", func() {
		It("subscribes the transport to all handled message types", func() {
			t := &mocks.InboundTransportMock{
				SubscribeFunc: func(context.Context, endpoint.Operation, ax.MessageTypeSet) error {
					return nil
				},
			}

			ctx := context.Background()
			err := dispatcher.Initialize(ctx, &endpoint.Endpoint{InboundTransport: t})
			Expect(err).ShouldNot(HaveOccurred())

			Expect(t.SubscribeCalls()).To(ConsistOf(
				struct {
					Ctx context.Context
					Op  endpoint.Operation
					Mt  ax.MessageTypeSet
				}{
					ctx,
					endpoint.OpSendUnicast,
					ax.TypesOf(
						&testmessages.Command{},
						&testmessages.Message{},
					),
				},
				struct {
					Ctx context.Context
					Op  endpoint.Operation
					Mt  ax.MessageTypeSet
				}{
					ctx,
					endpoint.OpSendMulticast,
					ax.TypesOf(
						&testmessages.Event{},
						&testmessages.Message{},
					),
				},
			))
		})
	})

	Describe("Accept", func() {
		ctx := context.Background()
		env := endpoint.InboundEnvelope{
			Envelope: ax.NewEnvelope(
				&testmessages.Message{},
			),
		}

		It("passes the message to each handler", func() {
			_ = dispatcher.Accept(ctx, sink, env)

			Expect(h1.HandleMessageCalls()).To(HaveLen(1))
			Expect(h2.HandleMessageCalls()).To(HaveLen(1))
			Expect(h3.HandleMessageCalls()).To(HaveLen(0))
		})

		It("passes a sender that sends messages via the message sink", func() {
			h1.HandleMessageFunc = func(ctx context.Context, s ax.Sender, _ ax.MessageContext) error {
				_, err := s.ExecuteCommand(ctx, &testmessages.Command{})
				return err
			}

			_ = dispatcher.Accept(ctx, sink, env)

			Expect(h1.HandleMessageCalls()).To(HaveLen(1))
			Expect(sink.Envelopes()).To(HaveLen(1))
		})

		It("returns nil if all handlers succeed", func() {
			err := dispatcher.Accept(ctx, sink, env)

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns an error if any handler fails", func() {
			h2.HandleMessageFunc = func(context.Context, ax.Sender, ax.MessageContext) error {
				return errors.New("<error>")
			}

			err := dispatcher.Accept(ctx, sink, env)

			Expect(err).To(MatchError("<error>"))
		})
	})
})
