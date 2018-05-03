package bus_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	. "github.com/jmalloc/ax/src/ax/bus"
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
			InitializeFunc:  func(context.Context, Transport) error { return nil },
			SendMessageFunc: func(context.Context, OutboundEnvelope) error { return nil },
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

	Describe("SendMessage", func() {
		Context("when there is a routing table", func() {
			BeforeEach(func() {
				t, err := NewRoutingTable(
					"ax.internal", "route-from-table",
				)
				Expect(err).ShouldNot(HaveOccurred())

				router.Routes = t
			})

			It("routes the message to the endpoint according to the routing table", func() {
				m := OutboundEnvelope{
					Operation: OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &messagetest.Message{},
					},
				}

				err := router.SendMessage(context.Background(), m)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(next.SendMessageCalls()).To(HaveLen(1))
				Expect(next.SendMessageCalls()[0].M.DestinationEndpoint).To(Equal("route-from-table"))
			})

			It("returns the same result for subsequent messages (coverage of cache hit)", func() {
				m := OutboundEnvelope{
					Operation: OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &messagetest.Message{},
					},
				}

				err := router.SendMessage(context.Background(), m)
				Expect(err).ShouldNot(HaveOccurred())

				err = router.SendMessage(context.Background(), m)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(next.SendMessageCalls()).To(HaveLen(2))

				a := next.SendMessageCalls()[0].M.DestinationEndpoint
				b := next.SendMessageCalls()[1].M.DestinationEndpoint

				Expect(a).To(Equal(b))
			})
		})

		Context("when there is no routing table", func() {
			It("routes the message to the endpoint named after the protocol buffers package name", func() {
				m := OutboundEnvelope{
					Operation: OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &messagetest.Message{},
					},
				}

				err := router.SendMessage(context.Background(), m)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(next.SendMessageCalls()).To(HaveLen(1))
				Expect(next.SendMessageCalls()[0].M.DestinationEndpoint).To(Equal("ax.internal.messagetest"))
			})

			It("returns an error if the message does not have a protocol buffers package name", func() {
				m := OutboundEnvelope{
					Operation: OpSendUnicast,
					Envelope: ax.Envelope{
						Message: &messagetest.NoPackage{},
					},
				}

				err := router.SendMessage(context.Background(), m)
				Expect(err).Should(HaveOccurred())
			})
		})

		It("does not replace the destination endpoint if it is already set", func() {
			m := OutboundEnvelope{
				Operation:           OpSendUnicast,
				DestinationEndpoint: "<endpoint>",
			}

			err := router.SendMessage(context.Background(), m)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.SendMessageCalls()).To(HaveLen(1))
			Expect(next.SendMessageCalls()[0].M).To(Equal(m))
		})

		It("does not set the destination endpoint for multicast messages", func() {
			m := OutboundEnvelope{Operation: OpSendMulticast}

			err := router.SendMessage(context.Background(), m)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.SendMessageCalls()).To(HaveLen(1))
			Expect(next.SendMessageCalls()[0].M).To(Equal(m))
		})
	})
})

var _ = Describe("RoutingTable", func() {
	var table RoutingTable

	Describe("NewRoutingTable", func() {
		It("returns an error when passed an odd number of arguments", func() {
			_, err := NewRoutingTable("foo")
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Lookup", func() {
		BeforeEach(func() {
			t, err := NewRoutingTable(
				"foo", "route:foo",
				"foo.qux", "route:foo.qux",
				"foo.bar.ExactMatch", "route:foo.bar.ExactMatch",
			)
			Expect(err).ShouldNot(HaveOccurred())
			table = t
		})

		It("favors an exact match", func() {
			ep, ok := table.Lookup(ax.MessageType{Name: "foo.bar.ExactMatch"})
			Expect(ok).To(BeTrue())
			Expect(ep).To(Equal("route:foo.bar.ExactMatch"))
		})

		It("returns the longest match when there is no exact match", func() {
			ep, ok := table.Lookup(ax.MessageType{Name: "foo.qux.Message"})
			Expect(ok).To(BeTrue())
			Expect(ep).To(Equal("route:foo.qux"))
		})

		Context("when there is no default route", func() {
			It("returns false for a message with no matching routes", func() {
				_, ok := table.Lookup(ax.MessageType{Name: "baz.qux.Message"})
				Expect(ok).To(BeFalse())
			})
		})

		Context("when there is a default route", func() {
			BeforeEach(func() {
				t, err := NewRoutingTable(
					"foo", "route:foo",
					"foo.qux", "route:foo.qux",
					"foo.bar.ExactMatch", "route:foo.bar.ExactMatch",
					"", "route:default",
				)
				Expect(err).ShouldNot(HaveOccurred())
				table = t
			})

			It("returns the default route for a message with no better matching routes", func() {
				ep, ok := table.Lookup(ax.MessageType{Name: "baz.qux.Message"})
				Expect(ok).To(BeTrue())
				Expect(ep).To(Equal("route:default"))
			})
		})
	})
})
