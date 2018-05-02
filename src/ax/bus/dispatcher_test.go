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
		sender     *bustest.MessageSenderMock
		dispatcher *Dispatcher
	)

	BeforeEach(func() {
		noOp := func(ax.MessageContext, ax.Message) error { return nil }

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

		sender = &bustest.MessageSenderMock{}
		dispatcher = &Dispatcher{Routes: t}
	})

	Describe("Initialize", func() {
		It("subscribes the transport to all handled event types", func() {
			t := &bustest.TransportMock{
				SubscribeFunc: func(_ context.Context, _ ax.MessageTypeSet) error {
					return nil
				},
			}

			err := dispatcher.Initialize(context.Background(), t)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(t.SubscribeCalls()).To(HaveLen(1))
			Expect(t.SubscribeCalls()[0].Mt.Members()).To(ConsistOf(
				ax.TypeOf(&messagetest.Event{}),
			))
		})
	})

	Describe("DeliverMessage", func() {
		message := InboundEnvelope{
			Envelope: ax.NewEnvelope(
				&messagetest.Message{},
			),
		}

		It("passes the message to each handler", func() {
			_ = dispatcher.DeliverMessage(
				context.Background(),
				sender,
				message,
			)

			Expect(h1.HandleMessageCalls()).To(HaveLen(1))
			Expect(h2.HandleMessageCalls()).To(HaveLen(1))
			Expect(h3.HandleMessageCalls()).To(HaveLen(0))
		})

		It("uses bus.MessageContext", func() {
			ctx := context.Background()

			_ = dispatcher.DeliverMessage(
				ctx,
				sender,
				message,
			)

			Expect(h1.HandleMessageCalls()).To(HaveLen(1))

			mc := h1.HandleMessageCalls()[0].Ctx.(*MessageContext)
			Expect(mc.Envelope).To(Equal(message.Envelope))
			Expect(mc.Sender).To(BeIdenticalTo(sender))
		})

		It("returns nil if all handlers succeed", func() {
			err := dispatcher.DeliverMessage(
				context.Background(),
				sender,
				message,
			)

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns an error if any handler fails", func() {
			h2.HandleMessageFunc = func(ax.MessageContext, ax.Message) error {
				return errors.New("<error>")
			}

			err := dispatcher.DeliverMessage(
				context.Background(),
				sender,
				message,
			)

			Expect(err).To(MatchError("<error>"))
		})
	})
})

var _ = Describe("DispatchTable", func() {
	var table DispatchTable

	Describe("NewDispatchTable", func() {
		It("returns an error when multiple handlers accept the same command", func() {
			h1 := MessageHandlerFunc(
				ax.TypesOf(&messagetest.Command{}),
				nil,
			)

			h2 := MessageHandlerFunc(
				ax.TypesOf(&messagetest.Command{}),
				nil,
			)

			_, err := NewDispatchTable(h1, h2)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Lookup", func() {
		h1 := MessageHandlerFunc(
			ax.TypesOf(&messagetest.Command{}),
			nil,
		)

		h2 := MessageHandlerFunc(
			ax.TypesOf(&messagetest.Event{}),
			nil,
		)

		h3 := MessageHandlerFunc(
			ax.TypesOf(&messagetest.Event{}),
			nil,
		)

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
