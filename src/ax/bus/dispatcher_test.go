package bus_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
	. "github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/internal/bustest"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dispatcher", func() {
	var (
		h1, h2, h3 *bustest.MessageHandlerMock
		sink       *BufferedSink
		dispatcher *Dispatcher
	)

	BeforeEach(func() {
		noOp := func(context.Context, ax.Sender, ax.Envelope) error { return nil }

		h1 = &bustest.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(
					&messagetest.Command{},
					&messagetest.Message{},
				)
			},
			HandleMessageFunc: noOp,
		}
		h2 = &bustest.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(
					&messagetest.Message{},
					&messagetest.Event{},
				)
			},
			HandleMessageFunc: noOp,
		}
		h3 = &bustest.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(
					&messagetest.Event{},
				)
			},
			HandleMessageFunc: noOp,
		}

		t, err := NewDispatchTable(h1, h2, h3)
		Expect(err).ShouldNot(HaveOccurred())

		sink = &BufferedSink{}
		dispatcher = &Dispatcher{Routes: t}
	})

	Describe("Initialize", func() {
		It("subscribes the transport to all handled message types", func() {
			t := &bustest.TransportMock{
				SubscribeFunc: func(_ context.Context, _ Operation, _ ax.MessageTypeSet) error {
					return nil
				},
			}

			ctx := context.Background()
			err := dispatcher.Initialize(ctx, t)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(t.SubscribeCalls()).To(ConsistOf(
				struct {
					Ctx context.Context
					Op  Operation
					Mt  ax.MessageTypeSet
				}{
					ctx,
					OpSendUnicast,
					ax.TypesOf(
						&messagetest.Command{},
						&messagetest.Message{},
					),
				},
				struct {
					Ctx context.Context
					Op  Operation
					Mt  ax.MessageTypeSet
				}{
					ctx,
					OpSendMulticast,
					ax.TypesOf(
						&messagetest.Event{},
						&messagetest.Message{},
					),
				},
			))
		})
	})

	Describe("Accept", func() {
		ctx := context.Background()
		env := InboundEnvelope{
			Envelope: ax.NewEnvelope(
				&messagetest.Message{},
			),
		}

		It("passes the message to each handler", func() {
			_ = dispatcher.Accept(ctx, sink, env)

			Expect(h1.HandleMessageCalls()).To(HaveLen(1))
			Expect(h2.HandleMessageCalls()).To(HaveLen(1))
			Expect(h3.HandleMessageCalls()).To(HaveLen(0))
		})

		It("uses a context that contains the message envelope", func() {
			h1.HandleMessageFunc = func(ctx context.Context, _ ax.Sender, _ ax.Envelope) error {
				defer GinkgoRecover()

				e, ok := GetEnvelope(ctx)

				Expect(ok).To(BeTrue())
				Expect(e).To(BeIdenticalTo(env.Envelope))

				return nil
			}

			_ = dispatcher.Accept(ctx, sink, env)

			Expect(h1.HandleMessageCalls()).To(HaveLen(1))
		})

		It("passes a sender that sends messages via the message sink", func() {
			h1.HandleMessageFunc = func(ctx context.Context, s ax.Sender, _ ax.Envelope) error {
				return s.ExecuteCommand(ctx, &messagetest.Command{})
			}

			_ = dispatcher.Accept(ctx, sink, env)

			Expect(h1.HandleMessageCalls()).To(HaveLen(1))
			Expect(sink.Envelopes).To(HaveLen(1))
		})

		It("returns nil if all handlers succeed", func() {
			err := dispatcher.Accept(ctx, sink, env)

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns an error if any handler fails", func() {
			h2.HandleMessageFunc = func(context.Context, ax.Sender, ax.Envelope) error {
				return errors.New("<error>")
			}

			err := dispatcher.Accept(ctx, sink, env)

			Expect(err).To(MatchError("<error>"))
		})
	})
})

var _ = Describe("DispatchTable", func() {
	var table DispatchTable

	Describe("NewDispatchTable", func() {
		It("returns an error when multiple handlers accept the same command", func() {
			h1 := &bustest.MessageHandlerMock{
				MessageTypesFunc: func() ax.MessageTypeSet {
					return ax.TypesOf(&messagetest.Command{})
				},
			}

			h2 := &bustest.MessageHandlerMock{
				MessageTypesFunc: func() ax.MessageTypeSet {
					return ax.TypesOf(&messagetest.Command{})
				},
			}

			_, err := NewDispatchTable(h1, h2)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Lookup", func() {
		h1 := &bustest.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(&messagetest.Command{})
			},
		}

		h2 := &bustest.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(&messagetest.Event{})
			},
		}

		h3 := &bustest.MessageHandlerMock{
			MessageTypesFunc: func() ax.MessageTypeSet {
				return ax.TypesOf(&messagetest.Event{})
			},
		}

		BeforeEach(func() {
			t, err := NewDispatchTable(h1, h2, h3)
			Expect(err).ShouldNot(HaveOccurred())
			table = t
		})

		It("returns all of the handlers that handle the given message type", func() {
			mt := ax.TypeOf(&messagetest.Event{})
			h := table.Lookup(mt)

			Expect(h).To(ConsistOf(h2, h3))
		})
	})
})
