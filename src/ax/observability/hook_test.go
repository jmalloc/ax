package observability_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax/bus"
	. "github.com/jmalloc/ax/src/ax/observability"
	"github.com/jmalloc/ax/src/internal/bustest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InboundHook", func() {
	var (
		next *bustest.InboundPipelineMock
		hook *InboundHook
	)

	BeforeEach(func() {
		next = &bustest.InboundPipelineMock{
			InitializeFunc: func(context.Context, bus.Transport) error { return nil },
			AcceptFunc:     func(context.Context, bus.MessageSink, bus.InboundEnvelope) error { return nil },
		}
		hook = &InboundHook{
			Next: next,
		}
	})

	Describe("Initialize", func() {
		It("initializes the next stage", func() {
			t := &bustest.TransportMock{}

			err := hook.Initialize(context.Background(), t)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.InitializeCalls()).To(HaveLen(1))
			Expect(next.InitializeCalls()[0].T).To(Equal(t))
		})
	})

	XDescribe("Accept", func() {
		Context("when the observer implements BeforeInboundObserver", func() {
			It("invokes the observer before invoking the next pipeline stage", func() {
			})

			It("passes the context returned by the observer to the next pipeline stage", func() {
			})

			It("passes the context returned by the observer to subsequent observers", func() {
			})
		})

		Context("when the observer implements AfterInboundObserver", func() {
			It("invokes the observer after invoking the next pipeline stage", func() {
			})

			It("passes the error returned by the next pipeline stage to the observer", func() {
			})
		})
	})
})

var _ = Describe("OutboundHook", func() {
	var (
		next *bustest.OutboundPipelineMock
		hook *OutboundHook
	)

	BeforeEach(func() {
		next = &bustest.OutboundPipelineMock{
			InitializeFunc: func(context.Context, bus.Transport) error { return nil },
			AcceptFunc:     func(context.Context, bus.OutboundEnvelope) error { return nil },
		}
		hook = &OutboundHook{
			Next: next,
		}
	})

	Describe("Initialize", func() {
		It("initializes the next stage", func() {
			t := &bustest.TransportMock{}

			err := hook.Initialize(context.Background(), t)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.InitializeCalls()).To(HaveLen(1))
			Expect(next.InitializeCalls()[0].T).To(Equal(t))
		})
	})

	XDescribe("Accept", func() {
		Context("when the observer implements BeforeOutboundObserver", func() {
			It("invokes the observer before invoking the next pipeline stage", func() {
			})

			It("passes the context returned by the observer to the next pipeline stage", func() {
			})

			It("passes the context returned by the observer to subsequent observers", func() {
			})
		})

		Context("when the observer implements AfterOutboundObserver", func() {
			It("invokes the observer after invoking the next pipeline stage", func() {
			})

			It("passes the error returned by the next pipeline stage to the observer", func() {
			})
		})
	})
})
