package routing_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	. "github.com/jmalloc/ax/src/ax/routing"
	"github.com/jmalloc/ax/src/internal/bustest"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Router", func() {
	var (
		next   *bustest.OutboundPipelineMock
		router *Router
	)

	BeforeEach(func() {
		next = &bustest.OutboundPipelineMock{
			InitializeFunc: func(context.Context, bus.Transport) error { return nil },
			AcceptFunc:     func(context.Context, bus.OutboundEnvelope) error { return nil },
		}
		router = &Router{
			Next: next,
		}
	})

	Describe("Initialize", func() {
		It("initializes the next stage", func() {
			t := &bustest.TransportMock{}

			err := router.Initialize(context.Background(), t)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.InitializeCalls()).To(HaveLen(1))
			Expect(next.InitializeCalls()[0].T).To(Equal(t))
		})
	})

	Describe("Accept", func() {
		Context("when there is a routing table", func() {
			BeforeEach(func() {
				t, err := NewEndpointTable(
					"ax.internal", "route-from-table",
				)
				Expect(err).ShouldNot(HaveOccurred())

				router.Routes = t
			})

			It("routes the message to the endpoint according to the routing table", func() {
				env := bus.OutboundEnvelope{
					Operation: bus.OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &messagetest.Message{},
					},
				}

				err := router.Accept(context.Background(), env)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(next.AcceptCalls()).To(HaveLen(1))
				Expect(next.AcceptCalls()[0].Env.DestinationEndpoint).To(Equal("route-from-table"))
			})

			It("returns the same result for subsequent messages (coverage of cache hit)", func() {
				env := bus.OutboundEnvelope{
					Operation: bus.OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &messagetest.Message{},
					},
				}

				err := router.Accept(context.Background(), env)
				Expect(err).ShouldNot(HaveOccurred())

				err = router.Accept(context.Background(), env)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(next.AcceptCalls()).To(HaveLen(2))

				a := next.AcceptCalls()[0].Env.DestinationEndpoint
				b := next.AcceptCalls()[1].Env.DestinationEndpoint

				Expect(a).To(Equal(b))
			})
		})

		Context("when there is no routing table", func() {
			It("routes the message to the endpoint named after the protocol buffers package name", func() {
				env := bus.OutboundEnvelope{
					Operation: bus.OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &messagetest.Message{},
					},
				}

				err := router.Accept(context.Background(), env)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(next.AcceptCalls()).To(HaveLen(1))
				Expect(next.AcceptCalls()[0].Env.DestinationEndpoint).To(Equal("ax.internal.messagetest"))
			})

			It("returns an error if the message does not have a protocol buffers package name", func() {
				env := bus.OutboundEnvelope{
					Operation: bus.OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &messagetest.NoPackage{},
					},
				}

				err := router.Accept(context.Background(), env)
				Expect(err).Should(HaveOccurred())
			})
		})

		It("does not replace the destination endpoint if it is already set", func() {
			env := bus.OutboundEnvelope{
				Operation:           bus.OpSendUnicast,
				DestinationEndpoint: "<endpoint>",
			}

			err := router.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.AcceptCalls()).To(HaveLen(1))
			Expect(next.AcceptCalls()[0].Env).To(Equal(env))
		})

		It("does not set the destination endpoint for multicast messages", func() {
			env := bus.OutboundEnvelope{
				Operation: bus.OpSendMulticast,
			}

			err := router.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.AcceptCalls()).To(HaveLen(1))
			Expect(next.AcceptCalls()[0].Env).To(Equal(env))
		})
	})
})
