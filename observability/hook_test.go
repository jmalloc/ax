package observability_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/axtest/mocks"
	"github.com/jmalloc/ax/axtest/testmessages"
	"github.com/jmalloc/ax/endpoint"
	. "github.com/jmalloc/ax/observability"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InboundHook", func() {
	var (
		observer *mocks.ObserverMock
		ep       *endpoint.Endpoint
		next     *mocks.InboundPipelineMock
		env      endpoint.InboundEnvelope
		hook     *InboundHook
	)

	BeforeEach(func() {
		ep = &endpoint.Endpoint{}
		observer = &mocks.ObserverMock{
			BeforeInboundFunc: func(context.Context, endpoint.InboundEnvelope) {},
			AfterInboundFunc:  func(context.Context, endpoint.InboundEnvelope, error) {},
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
			Observers: []Observer{observer},
		}
	})

	Describe("Initialize", func() {
		It("initializes the next stage", func() {
			err := hook.Initialize(context.Background(), ep)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.InitializeCalls()).To(HaveLen(1))
			Expect(next.InitializeCalls()[0].Ep).To(Equal(ep))
		})
	})

	Describe("Accept", func() {
		BeforeEach(func() {
			if err := hook.Initialize(context.Background(), ep); err != nil {
				panic(err)
			}
		})

		It("invokes the before-observer before processing the message", func() {
			observer.BeforeInboundFunc = func(_ context.Context, e endpoint.InboundEnvelope) {
				Expect(e).To(Equal(env))
				Expect(next.AcceptCalls()).To(BeEmpty()) // ensure message has not been processed yet
			}

			err := hook.Accept(context.Background(), nil /*sink*/, env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(observer.BeforeInboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		It("invokes the after-observer after processing the message", func() {
			observer.AfterInboundFunc = func(_ context.Context, e endpoint.InboundEnvelope, err error) {
				Expect(e).To(Equal(env))
				Expect(err).To(BeNil())
				Expect(next.AcceptCalls()).To(HaveLen(1)) // ensure message has already been processed
			}

			err := hook.Accept(context.Background(), nil /*sink*/, env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(observer.AfterInboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		It("provides the after-observer with the message processing error", func() {
			expected := errors.New("<error>")
			next.AcceptFunc = func(context.Context, endpoint.MessageSink, endpoint.InboundEnvelope) error {
				return expected
			}

			observer.AfterInboundFunc = func(_ context.Context, e endpoint.InboundEnvelope, err error) {
				Expect(err).To(Equal(expected))
			}

			err := hook.Accept(context.Background(), nil /*sink*/, env)
			Expect(err).To(Equal(err))
			Expect(observer.AfterInboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})
	})
})

var _ = Describe("OutboundHook", func() {
	var (
		observer *mocks.ObserverMock
		ep       *endpoint.Endpoint
		next     *mocks.OutboundPipelineMock
		env      endpoint.OutboundEnvelope
		hook     *OutboundHook
	)

	BeforeEach(func() {
		ep = &endpoint.Endpoint{}
		observer = &mocks.ObserverMock{
			BeforeOutboundFunc: func(context.Context, endpoint.OutboundEnvelope) {},
			AfterOutboundFunc:  func(context.Context, endpoint.OutboundEnvelope, error) {},
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
			Observers: []Observer{observer},
		}
	})

	Describe("Initialize", func() {
		It("initializes the next stage", func() {
			err := hook.Initialize(context.Background(), ep)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(next.InitializeCalls()).To(HaveLen(1))
			Expect(next.InitializeCalls()[0].Ep).To(Equal(ep))
		})
	})

	Describe("Accept", func() {
		BeforeEach(func() {
			if err := hook.Initialize(context.Background(), ep); err != nil {
				panic(err)
			}
		})

		It("invokes the before-observer before processing the message", func() {
			observer.BeforeOutboundFunc = func(_ context.Context, e endpoint.OutboundEnvelope) {
				Expect(e).To(Equal(env))
				Expect(next.AcceptCalls()).To(BeEmpty()) // ensure message has not been processed yet
			}

			err := hook.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(observer.BeforeOutboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		It("invokes the after-observer after processing the message", func() {
			observer.AfterOutboundFunc = func(_ context.Context, e endpoint.OutboundEnvelope, err error) {
				Expect(e).To(Equal(env))
				Expect(err).To(BeNil())
				Expect(next.AcceptCalls()).To(HaveLen(1)) // ensure message has already been processed
			}

			err := hook.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(observer.AfterOutboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})

		It("provides the after-observer with the message processing error", func() {
			expected := errors.New("<error>")
			next.AcceptFunc = func(context.Context, endpoint.OutboundEnvelope) error {
				return expected
			}

			observer.AfterOutboundFunc = func(_ context.Context, e endpoint.OutboundEnvelope, err error) {
				Expect(err).To(Equal(expected))
			}

			err := hook.Accept(context.Background(), env)
			Expect(err).To(Equal(err))
			Expect(observer.AfterOutboundCalls()).To(HaveLen(1)) // ensure observer is actually called
		})
	})
})
