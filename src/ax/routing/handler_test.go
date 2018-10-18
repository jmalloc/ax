package routing_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax/endpoint"

	"github.com/jmalloc/ax/src/ax"
	. "github.com/jmalloc/ax/src/ax/routing"
	"github.com/jmalloc/ax/src/axtest/testmessages"
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

func (h *messageHandler) HandleWithMessageContext(mctx ax.MessageContext, msg *testmessages.MessageB) {
	h.Letter = "B"
}

func (h *messageHandler) HandleWithError(ctx context.Context, msg *testmessages.MessageC) error {
	h.Letter = "C"
	return errors.New("C")
}

func (h *messageHandler) HandleWithErrorAndMessageContext(ctx context.Context, mctx ax.MessageContext, msg *testmessages.MessageD) error {
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
					ax.MessageContext{
						Envelope: ax.NewEnvelope(m),
						Sender:   &endpoint.SinkSender{},
					},
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
				"HandleWithError",
				&testmessages.MessageC{},
				"C",
				true,
			),
			Entry(
				"HandleWithErrorAndMessageContext",
				&testmessages.MessageD{},
				"D",
				true,
			),
		)
	})
})
