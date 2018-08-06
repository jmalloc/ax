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

func (h *messageHandler) HandleWithEnvelope(msg *testmessages.MessageB, env ax.Envelope) {
	h.Letter = "B"
}

func (h *messageHandler) HandleWithContext(ctx context.Context, msg *testmessages.MessageC) error {
	h.Letter = "C"
	return errors.New("C")
}

func (h *messageHandler) HandleWithContextAndEnvelope(ctx context.Context, msg *testmessages.MessageD, env ax.Envelope) error {
	h.Letter = "D"
	return errors.New("D")
}

func (h *messageHandler) HandleWithSender(ctx context.Context, s ax.Sender, msg *testmessages.MessageE) error {
	h.Letter = "E"
	return errors.New("E")
}

func (h *messageHandler) HandleWithSenderAndEnvelope(ctx context.Context, s ax.Sender, msg *testmessages.MessageF, env ax.Envelope) error {
	h.Letter = "F"
	return errors.New("F")
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
				ax.TypeOf(&testmessages.MessageE{}),
				ax.TypeOf(&testmessages.MessageF{}),
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
					ax.NewEnvelope(m),
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
				"HandleWithEnvelope",
				&testmessages.MessageB{},
				"B",
				false,
			),
			Entry(
				"HandleWithContext",
				&testmessages.MessageC{},
				"C",
				true,
			),
			Entry(
				"HandleWithContextAndEnvelope",
				&testmessages.MessageD{},
				"D",
				true,
			),
			Entry(
				"HandleWithSender",
				&testmessages.MessageE{},
				"E",
				true,
			),
			Entry(
				"HandleWithSenderAndEnvelope",
				&testmessages.MessageF{},
				"F",
				true,
			),
		)
	})
})
