package bus_test

import (
	"context"

	. "github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/internal/bustest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TransportStage", func() {
	Describe("Accept", func() {
		It("sends the message via the transport provided at initialization time", func() {
			stage := &TransportStage{}

			ctx := context.Background()
			tr := &bustest.TransportMock{
				SendFunc: func(context.Context, OutboundEnvelope) error {
					return nil
				},
			}

			stage.Initialize(ctx, tr)

			env := OutboundEnvelope{}
			env.MessageID.GenerateUUID()
			stage.Accept(ctx, env)

			Expect(tr.SendCalls()).To(HaveLen(1))
			Expect(tr.SendCalls()[0].Env).To(Equal(env))
		})
	})
})
