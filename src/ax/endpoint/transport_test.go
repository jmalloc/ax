package endpoint_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax"

	. "github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/axtest/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TransportStage", func() {
	Describe("Accept", func() {
		It("sends the message via the transport provided at initialization time", func() {
			stage := &TransportStage{}

			ctx := context.Background()
			tr := &mocks.OutboundTransportMock{
				SendFunc: func(context.Context, OutboundEnvelope) error {
					return nil
				},
			}

			stage.Initialize(ctx, &Endpoint{OutboundTransport: tr})

			env := OutboundEnvelope{
				Envelope: ax.Envelope{
					MessageID: ax.GenerateMessageID(),
				},
			}
			stage.Accept(ctx, env)

			Expect(tr.SendCalls()).To(HaveLen(1))
			Expect(tr.SendCalls()[0].Env).To(Equal(env))
		})
	})
})
