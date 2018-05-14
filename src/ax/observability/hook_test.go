package observability_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	. "github.com/jmalloc/ax/src/ax/observability"
	"github.com/jmalloc/ax/src/internal/bustest"
	"github.com/jmalloc/ax/src/internal/messagetest"
	"github.com/jmalloc/ax/src/internal/observabilitytest"
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

	Describe("Accept", func() {
		Context("when the observer implements BeforeInboundObserver", func() {
			var observer *observabilitytest.BeforeInboundObserverMock

			BeforeEach(func() {
				observer = &observabilitytest.BeforeInboundObserverMock{}

				hook.Observers = []Observer{observer}

				if err := hook.Initialize(
					context.Background(),
					&bustest.TransportMock{},
				); err != nil {
					panic(err)
				}
			})

			It("invokes the observer before invoking the next pipeline stage", func() {
				env := bus.InboundEnvelope{
					Envelope: ax.NewEnvelope(
						&messagetest.Message{},
					),
				}

				observer.BeforeInboundFunc = func(
					_ context.Context,
					e bus.InboundEnvelope,
				) {
					Expect(e).To(Equal(env))
					Expect(next.AcceptCalls()).To(BeEmpty())
				}

				err := hook.Accept(context.Background(), nil /*sink*/, env)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(observer.BeforeInboundCalls()).To(HaveLen(1))
			})
		})

		Context("when the observer implements AfterInboundObserver", func() {
			XIt("invokes the observer after invoking the next pipeline stage", func() {
			})

			XIt("passes the error returned by the next pipeline stage to the observer", func() {
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

	Describe("Accept", func() {
		var observer *observabilitytest.BeforeOutboundObserverMock

		BeforeEach(func() {
			observer = &observabilitytest.BeforeOutboundObserverMock{}

			hook.Observers = []Observer{observer}

			if err := hook.Initialize(
				context.Background(),
				&bustest.TransportMock{},
			); err != nil {
				panic(err)
			}
		})

		Context("when the observer implements BeforeOutboundObserver", func() {
			It("invokes the observer before invoking the next pipeline stage", func() {
				env := bus.OutboundEnvelope{
					Envelope: ax.NewEnvelope(
						&messagetest.Message{},
					),
				}

				observer.BeforeOutboundFunc = func(
					_ context.Context,
					e bus.OutboundEnvelope,
				) {
					Expect(e).To(Equal(env))
					Expect(next.AcceptCalls()).To(BeEmpty())
				}

				err := hook.Accept(context.Background(), env)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(observer.BeforeOutboundCalls()).To(HaveLen(1))
			})
		})

		Context("when the observer implements AfterOutboundObserver", func() {
			XIt("invokes the observer after invoking the next pipeline stage", func() {
			})

			XIt("passes the error returned by the next pipeline stage to the observer", func() {
			})
		})
	})
})
