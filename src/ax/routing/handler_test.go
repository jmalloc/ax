package routing_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	. "github.com/jmalloc/ax/src/ax/routing"
	"github.com/jmalloc/ax/src/axtest/testmessages"
	"github.com/jmalloc/twelf/src/twelf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type messageHandler struct {
	Letter string
}

func (h *messageHandler) Handle(msg *testmessages.MessageA) {
	h.Letter = "A"
}

func (h *messageHandler) HandleWithMessageContext(*testmessages.MessageB, ax.MessageContext) {
	h.Letter = "B"
}

func (h *messageHandler) HandleWithSender(context.Context, ax.Sender, *testmessages.MessageC) error {
	h.Letter = "C"
	return errors.New("C")
}

func (h *messageHandler) HandleWithSenderAndMessageContext(
	context.Context,
	ax.Sender,
	*testmessages.MessageD,
	ax.MessageContext,
) error {
	h.Letter = "D"
	return errors.New("D")
}

var _ = Describe("MessageHandler", func() {
	var handler MessageHandler
	var value *messageHandler

	BeforeEach(func() {
		value = &messageHandler{}
		handler = NewMessageHandler(value)
	})

	Describe("MessageTypes", func() {
		It("returns the types of the handled messages", func() {
			Expect(handler.MessageTypes().Members()).To(ConsistOf(
				ax.TypeOf(&testmessages.MessageA{}),
				ax.TypeOf(&testmessages.MessageB{}),
				ax.TypeOf(&testmessages.MessageC{}),
				ax.TypeOf(&testmessages.MessageD{}),
			))
		})
	})

	Describe("HandleMessage", func() {
		DescribeTable(
			"calls the expected method",
			func(
				m ax.Message,
				letter string,
				expectErr bool,
			) {
				err := handler.HandleMessage(
					context.Background(),
					&endpoint.SinkSender{},
					ax.NewMessageContext(
						ax.NewEnvelope(m),
						nil, // span
						twelf.SilentLogger,
					),
				)

				Expect(value.Letter).To(Equal(letter))

				if expectErr {
					Expect(err).To(MatchError(letter))
				} else {
					Expect(err).ShouldNot(HaveOccurred())
				}
			},
			Entry(
				"Handle",
				&testmessages.MessageA{},
				"A",
				false,
			),
			Entry(
				"HandleWithMessageContext",
				&testmessages.MessageB{},
				"B",
				false,
			),
			Entry(
				"HandleWithSender",
				&testmessages.MessageC{},
				"C",
				true,
			),
			Entry(
				"HandleWithSenderAndMessageContext",
				&testmessages.MessageD{},
				"D",
				true,
			),
		)
	})
})
