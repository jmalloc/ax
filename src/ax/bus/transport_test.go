package bus_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	. "github.com/jmalloc/ax/src/ax/bus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type transport struct {
	sentMessage OutboundEnvelope
}

func (*transport) Initialize(context.Context, string) error                { panic("not implemented") }
func (*transport) Subscribe(context.Context, ax.MessageTypeSet) error      { panic("not implemented") }
func (*transport) ReceiveMessage(context.Context) (InboundEnvelope, error) { panic("not implemented") }
func (t *transport) SendMessage(_ context.Context, m OutboundEnvelope) error {
	t.sentMessage = m
	return nil
}

var _ = Describe("TransportStage", func() {
	Describe("SendMessage", func() {
		It("sends the message via the transport provided at initialization time", func() {
			stage := &TransportStage{}
			tr := &transport{}
			ctx := context.Background()

			stage.Initialize(ctx, tr)

			m := OutboundEnvelope{}
			m.MessageID.GenerateUUID()
			stage.SendMessage(ctx, m)

			Expect(tr.sentMessage).To(Equal(m))
		})
	})
})
