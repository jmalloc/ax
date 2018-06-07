package observability_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/axtest/mocks"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	. "github.com/jmalloc/ax/src/ax/observability"
	"github.com/jmalloc/ax/src/axtest/testmessages"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InboundHook", func() {
	var (
		before *mocks.BeforeInboundObserverMock
		after  *mocks.AfterInboundObserverMock
		ep     *endpoint.Endpoint
		next   *mocks.InboundPipelineMock
		env    endpoint.InboundEnvelope
		hook   *InboundHook
	)

	BeforeEach(func() {
		ep = &endpoint.Endpoint{}
		before = &mocks.BeforeInboundObserverMock{
			BeforeInboundFunc: func(context.Context, endpoint.InboundEnvelope) {},
		}
		after = &mocks.AfterInboundObserverMock{
			AfterInboundFunc: func(context.Context, endpoint.InboundEnvelope, error) {},
		}
		next = &mocks.InboundPipelineMock{
			InitializeFunc: func(context.Context, *endpoint.Endpoint) error { return nil },
			AcceptFunc:     func(context.Context, endpoint.MessageSink, endpoint.InboundEnvelope) error { return nil },
		}
		env = endpoint.InboundEnvelope{
			Envelope: ax.NewEnvelope(
				&testmessages.Message{},
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
			hook.Observers = append(hook.Observers, &mocks.BeforeOutboundObserverMock{})

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
		before *mocks.BeforeOutboundObserverMock
		after  *mocks.AfterOutboundObserverMock
		ep     *endpoint.Endpoint
		next   *mocks.OutboundPipelineMock
		env    endpoint.OutboundEnvelope
		hook   *OutboundHook
	)

	BeforeEach(func() {
		ep = &endpoint.Endpoint{}
		before = &mocks.BeforeOutboundObserverMock{
			BeforeOutboundFunc: func(context.Context, endpoint.OutboundEnvelope) {},
		}
		after = &mocks.AfterOutboundObserverMock{
			AfterOutboundFunc: func(context.Context, endpoint.OutboundEnvelope, error) {},
		}
		next = &mocks.OutboundPipelineMock{
			InitializeFunc: func(context.Context, *endpoint.Endpoint) error { return nil },
			AcceptFunc:     func(context.Context, endpoint.OutboundEnvelope) error { return nil },
		}
		env = endpoint.OutboundEnvelope{
			Envelope: ax.NewEnvelope(
				&testmessages.Message{},
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
			hook.Observers = append(hook.Observers, &mocks.BeforeInboundObserverMock{})

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
