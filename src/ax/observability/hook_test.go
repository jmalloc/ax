package observability_test

import (
	"context"
	"errors"

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
		before    *observabilitytest.BeforeInboundObserverMock
		after     *observabilitytest.AfterInboundObserverMock
		transport *bustest.TransportMock
		next      *bustest.InboundPipelineMock
		env       bus.InboundEnvelope
		hook      *InboundHook
	)

	BeforeEach(func() {
		transport = &bustest.TransportMock{}
		before = &observabilitytest.BeforeInboundObserverMock{
			BeforeInboundFunc: func(context.Context, bus.InboundEnvelope) {},
		}
		after = &observabilitytest.AfterInboundObserverMock{
			AfterInboundFunc: func(context.Context, bus.InboundEnvelope, error) {},
		}
		next = &bustest.InboundPipelineMock{
			InitializeFunc: func(context.Context, bus.Transport) error { return nil },
			AcceptFunc:     func(context.Context, bus.MessageSink, bus.InboundEnvelope) error { return nil },
		}
		env = bus.InboundEnvelope{
			Envelope: ax.NewEnvelope(
				&messagetest.Message{},
			),
		}
		hook = &InboundHook{
			Next:      next,
			Observers: []Observer{before, after},
		}
	})

	Describe("Initialize", func() {
		It("initializes the next stage", func() {
			err := hook.Initialize(context.Background(), transport)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.InitializeCalls()).To(HaveLen(1))
			Expect(next.InitializeCalls()[0].T).To(Equal(transport))
		})

		It("panics if an observer does not implement either of the inbound observer interfaces", func() {
			// outbound observer instead of inbound
			hook.Observers = append(hook.Observers, &observabilitytest.BeforeOutboundObserverMock{})

			Expect(func() {
				hook.Initialize(context.Background(), transport)
			}).To(Panic())
		})
	})

	Describe("Accept", func() {
		BeforeEach(func() {
			if err := hook.Initialize(context.Background(), transport); err != nil {
				panic(err)
			}
		})

		It("invokes the before-observer before processing the message", func() {
			before.BeforeInboundFunc = func(_ context.Context, e bus.InboundEnvelope) {
				Expect(e).To(Equal(env))
				Expect(next.AcceptCalls()).To(BeEmpty()) // ensure message has not been processed yet
			}

			err := hook.Accept(context.Background(), nil /*sink*/, env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(before.BeforeInboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		It("invokes the after-observer after processing the message", func() {
			after.AfterInboundFunc = func(_ context.Context, e bus.InboundEnvelope, err error) {
				Expect(e).To(Equal(env))
				Expect(err).To(BeNil())
				Expect(next.AcceptCalls()).To(HaveLen(1)) // ensure message has already been processed
			}

			err := hook.Accept(context.Background(), nil /*sink*/, env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(after.AfterInboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		XIt("provides the after-observer with the message processing error", func() {
			expected := errors.New("<error>")
			next.AcceptFunc = func(context.Context, bus.MessageSink, bus.InboundEnvelope) error {
				return expected
			}

			after.AfterInboundFunc = func(_ context.Context, e bus.InboundEnvelope, err error) {
				Expect(err).To(Equal(expected))
			}

			err := hook.Accept(context.Background(), nil /*sink*/, env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(after.AfterInboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})
	})
})

var _ = Describe("OutboundHook", func() {
	var (
		before    *observabilitytest.BeforeOutboundObserverMock
		after     *observabilitytest.AfterOutboundObserverMock
		transport *bustest.TransportMock
		next      *bustest.OutboundPipelineMock
		env       bus.OutboundEnvelope
		hook      *OutboundHook
	)

	BeforeEach(func() {
		transport = &bustest.TransportMock{}
		before = &observabilitytest.BeforeOutboundObserverMock{
			BeforeOutboundFunc: func(context.Context, bus.OutboundEnvelope) {},
		}
		after = &observabilitytest.AfterOutboundObserverMock{
			AfterOutboundFunc: func(context.Context, bus.OutboundEnvelope, error) {},
		}
		next = &bustest.OutboundPipelineMock{
			InitializeFunc: func(context.Context, bus.Transport) error { return nil },
			AcceptFunc:     func(context.Context, bus.OutboundEnvelope) error { return nil },
		}
		env = bus.OutboundEnvelope{
			Envelope: ax.NewEnvelope(
				&messagetest.Message{},
			),
		}
		hook = &OutboundHook{
			Next:      next,
			Observers: []Observer{before, after},
		}
	})

	Describe("Initialize", func() {
		It("initializes the next stage", func() {
			err := hook.Initialize(context.Background(), transport)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.InitializeCalls()).To(HaveLen(1))
			Expect(next.InitializeCalls()[0].T).To(Equal(transport))
		})

		It("panics if an observer does not implement either of the outbound observer interfaces", func() {
			// inbound observer instead of outbound
			hook.Observers = append(hook.Observers, &observabilitytest.BeforeInboundObserverMock{})

			Expect(func() {
				hook.Initialize(context.Background(), transport)
			}).To(Panic())
		})
	})

	Describe("Accept", func() {
		BeforeEach(func() {
			if err := hook.Initialize(context.Background(), transport); err != nil {
				panic(err)
			}
		})

		It("invokes the before-observer before processing the message", func() {
			before.BeforeOutboundFunc = func(_ context.Context, e bus.OutboundEnvelope) {
				Expect(e).To(Equal(env))
				Expect(next.AcceptCalls()).To(BeEmpty()) // ensure message has not been processed yet
			}

			err := hook.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(before.BeforeOutboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		It("invokes the after-observer after processing the message", func() {
			after.AfterOutboundFunc = func(_ context.Context, e bus.OutboundEnvelope, err error) {
				Expect(e).To(Equal(env))
				Expect(err).To(BeNil())
				Expect(next.AcceptCalls()).To(HaveLen(1)) // ensure message has already been processed
			}

			err := hook.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(after.AfterOutboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		XIt("provides the after-observer with the message processing error", func() {
			expected := errors.New("<error>")
			next.AcceptFunc = func(context.Context, bus.OutboundEnvelope) error {
				return expected
			}

			after.AfterOutboundFunc = func(_ context.Context, e bus.OutboundEnvelope, err error) {
				Expect(err).To(Equal(expected))
			}

			err := hook.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(after.AfterOutboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})
	})
})
