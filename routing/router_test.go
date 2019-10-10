package routing_test

import (
	"context"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/axtest/mocks"
	"github.com/jmalloc/ax/axtest/testmessages"
	"github.com/jmalloc/ax/endpoint"
	. "github.com/jmalloc/ax/routing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ endpoint.OutboundPipeline = (*Router)(nil) // ensure Router implements OutboundPipeline

var _ = Describe("Router", func() {
	var (
		next   *mocks.OutboundPipelineMock
		router *Router
	)

	BeforeEach(func() {
		next = &mocks.OutboundPipelineMock{
			InitializeFunc: func(context.Context, *endpoint.Endpoint) error { return nil },
			AcceptFunc:     func(context.Context, endpoint.OutboundEnvelope) error { return nil },
		}
		router = &Router{
			Next: next,
		}
	})

	Describe("Initialize", func() {
		It("initializes the next stage", func() {
			ep := &endpoint.Endpoint{}

			err := router.Initialize(context.Background(), ep)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.InitializeCalls()).To(HaveLen(1))
			Expect(next.InitializeCalls()[0].Ep).To(Equal(ep))
		})
	})

	Describe("Accept", func() {
		Context("when there is a routing table", func() {
			BeforeEach(func() {
				t, err := NewEndpointTable(
					"axtest", "route-from-table",
				)
				Expect(err).ShouldNot(HaveOccurred())

				router.Routes = t
			})

			It("routes the message to the endpoint according to the routing table", func() {
				env := endpoint.OutboundEnvelope{
					Operation: endpoint.OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &testmessages.Message{},
					},
				}

				err := router.Accept(context.Background(), env)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(next.AcceptCalls()).To(HaveLen(1))
				Expect(next.AcceptCalls()[0].Env.DestinationEndpoint).To(Equal("route-from-table"))
			})

			It("returns the same result for subsequent messages (coverage of cache hit)", func() {
				env := endpoint.OutboundEnvelope{
					Operation: endpoint.OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &testmessages.Message{},
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
				env := endpoint.OutboundEnvelope{
					Operation: endpoint.OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &testmessages.Message{},
					},
				}

				err := router.Accept(context.Background(), env)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(next.AcceptCalls()).To(HaveLen(1))
				Expect(next.AcceptCalls()[0].Env.DestinationEndpoint).To(Equal("axtest.testmessages"))
			})

			It("returns an error if the message does not have a protocol buffers package name", func() {
				env := endpoint.OutboundEnvelope{
					Operation: endpoint.OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &testmessages.NoPackage{},
					},
				}

				err := router.Accept(context.Background(), env)
				Expect(err).Should(HaveOccurred())
			})
		})

		It("does not replace the destination endpoint if it is already set", func() {
			env := endpoint.OutboundEnvelope{
				Operation:           endpoint.OpSendUnicast,
				DestinationEndpoint: "<endpoint>",
			}

			err := router.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.AcceptCalls()).To(HaveLen(1))
			Expect(next.AcceptCalls()[0].Env).To(Equal(env))
		})

		It("does not set the destination endpoint for multicast messages", func() {
			env := endpoint.OutboundEnvelope{
				Operation: endpoint.OpSendMulticast,
			}

			err := router.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.AcceptCalls()).To(HaveLen(1))
			Expect(next.AcceptCalls()[0].Env).To(Equal(env))
		})
	})
})
