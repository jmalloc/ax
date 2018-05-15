package observability_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	. "github.com/jmalloc/ax/src/ax/observability"
	"github.com/jmalloc/ax/src/internal/endpointtest"
	"github.com/jmalloc/ax/src/internal/messagetest"
	"github.com/jmalloc/ax/src/internal/observabilitytest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InboundHook", func() {
	var (
		before *observabilitytest.BeforeInboundObserverMock
		after  *observabilitytest.AfterInboundObserverMock
		ep     *endpoint.Endpoint
		next   *endpointtest.InboundPipelineMock
		env    endpoint.InboundEnvelope
		hook   *InboundHook
	)

	BeforeEach(func() {
		ep = &endpoint.Endpoint{}
		before = &observabilitytest.BeforeInboundObserverMock{
			BeforeInboundFunc: func(context.Context, endpoint.InboundEnvelope) {},
		}
		after = &observabilitytest.AfterInboundObserverMock{
			AfterInboundFunc: func(context.Context, endpoint.InboundEnvelope, error) {},
		}
		next = &endpointtest.InboundPipelineMock{
			InitializeFunc: func(context.Context, *endpoint.Endpoint) error { return nil },
			AcceptFunc:     func(context.Context, endpoint.MessageSink, endpoint.InboundEnvelope) error { return nil },
		}
		env = endpoint.InboundEnvelope{
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
			err := hook.Initialize(context.Background(), ep)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.InitializeCalls()).To(HaveLen(1))
			Expect(next.InitializeCalls()[0].Ep).To(Equal(ep))
		})

		It("panics if an observer does not implement either of the inbound observer interfaces", func() {
			// outbound observer instead of inbound
			hook.Observers = append(hook.Observers, &observabilitytest.BeforeOutboundObserverMock{})

			Expect(func() {
				hook.Initialize(context.Background(), ep)
			}).To(Panic())
		})
	})

	Describe("Accept", func() {
		BeforeEach(func() {
			if err := hook.Initialize(context.Background(), ep); err != nil {
				panic(err)
			}
		})

		It("invokes the before-observer before processing the message", func() {
			before.BeforeInboundFunc = func(_ context.Context, e endpoint.InboundEnvelope) {
				Expect(e).To(Equal(env))
				Expect(next.AcceptCalls()).To(BeEmpty()) // ensure message has not been processed yet
			}

			err := hook.Accept(context.Background(), nil /*sink*/, env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(before.BeforeInboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		It("invokes the after-observer after processing the message", func() {
			after.AfterInboundFunc = func(_ context.Context, e endpoint.InboundEnvelope, err error) {
				Expect(e).To(Equal(env))
				Expect(err).To(BeNil())
				Expect(next.AcceptCalls()).To(HaveLen(1)) // ensure message has already been processed
			}

			err := hook.Accept(context.Background(), nil /*sink*/, env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(after.AfterInboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		It("provides the after-observer with the message processing error", func() {
			expected := errors.New("<error>")
			next.AcceptFunc = func(context.Context, endpoint.MessageSink, endpoint.InboundEnvelope) error {
				return expected
			}

			after.AfterInboundFunc = func(_ context.Context, e endpoint.InboundEnvelope, err error) {
				Expect(err).To(Equal(expected))
			}

			err := hook.Accept(context.Background(), nil /*sink*/, env)
			Expect(err).To(Equal(err))
			Expect(after.AfterInboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})
	})
})

var _ = Describe("OutboundHook", func() {
	var (
		before *observabilitytest.BeforeOutboundObserverMock
		after  *observabilitytest.AfterOutboundObserverMock
		ep     *endpoint.Endpoint
		next   *endpointtest.OutboundPipelineMock
		env    endpoint.OutboundEnvelope
		hook   *OutboundHook
	)

	BeforeEach(func() {
		ep = &endpoint.Endpoint{}
		before = &observabilitytest.BeforeOutboundObserverMock{
			BeforeOutboundFunc: func(context.Context, endpoint.OutboundEnvelope) {},
		}
		after = &observabilitytest.AfterOutboundObserverMock{
			AfterOutboundFunc: func(context.Context, endpoint.OutboundEnvelope, error) {},
		}
		next = &endpointtest.OutboundPipelineMock{
			InitializeFunc: func(context.Context, *endpoint.Endpoint) error { return nil },
			AcceptFunc:     func(context.Context, endpoint.OutboundEnvelope) error { return nil },
		}
		env = endpoint.OutboundEnvelope{
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
			err := hook.Initialize(context.Background(), ep)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.InitializeCalls()).To(HaveLen(1))
			Expect(next.InitializeCalls()[0].Ep).To(Equal(ep))
		})

		It("panics if an observer does not implement either of the outbound observer interfaces", func() {
			// inbound observer instead of outbound
			hook.Observers = append(hook.Observers, &observabilitytest.BeforeInboundObserverMock{})

			Expect(func() {
				hook.Initialize(context.Background(), ep)
			}).To(Panic())
		})
	})

	Describe("Accept", func() {
		BeforeEach(func() {
			if err := hook.Initialize(context.Background(), ep); err != nil {
				panic(err)
			}
		})

		It("invokes the before-observer before processing the message", func() {
			before.BeforeOutboundFunc = func(_ context.Context, e endpoint.OutboundEnvelope) {
				Expect(e).To(Equal(env))
				Expect(next.AcceptCalls()).To(BeEmpty()) // ensure message has not been processed yet
			}

			err := hook.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(before.BeforeOutboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		It("invokes the after-observer after processing the message", func() {
			after.AfterOutboundFunc = func(_ context.Context, e endpoint.OutboundEnvelope, err error) {
				Expect(e).To(Equal(env))
				Expect(err).To(BeNil())
				Expect(next.AcceptCalls()).To(HaveLen(1)) // ensure message has already been processed
			}

			err := hook.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(after.AfterOutboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		It("provides the after-observer with the message processing error", func() {
			expected := errors.New("<error>")
			next.AcceptFunc = func(context.Context, endpoint.OutboundEnvelope) error {
				return expected
			}

			after.AfterOutboundFunc = func(_ context.Context, e endpoint.OutboundEnvelope, err error) {
				Expect(err).To(Equal(expected))
			}

			err := hook.Accept(context.Background(), env)
			Expect(err).To(Equal(err))
			Expect(after.AfterOutboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})
	})
})
