package persistence_test

import (
	"context"

	"github.com/jmalloc/ax/axtest/mocks"
	"github.com/jmalloc/ax/endpoint"
	. "github.com/jmalloc/ax/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InboundInjector", func() {
	var (
		next *mocks.InboundPipelineMock
		inj  *InboundInjector
	)

	BeforeEach(func() {
		next = &mocks.InboundPipelineMock{}
		inj = &InboundInjector{
			DataStore: &mocks.DataStoreMock{},
			Next:      next,
		}
	})

	Describe("Initialize", func() {
		It("calls the next pipeline with the data store in the context", func() {
			next.InitializeFunc = func(ctx context.Context, _ *endpoint.Endpoint) error {
				ds, ok := GetDataStore(ctx)
				Expect(ok).To(BeTrue())
				Expect(ds).To(Equal(inj.DataStore))
				return nil
			}

			inj.Initialize(context.Background(), &endpoint.Endpoint{})

			Expect(next.InitializeCalls()).To(HaveLen(1))
		})
	})

	Describe("Accept", func() {
		It("calls the next pipeline with the data store in the context", func() {
			next.AcceptFunc = func(ctx context.Context, _ endpoint.MessageSink, _ endpoint.InboundEnvelope) error {
				ds, ok := GetDataStore(ctx)
				Expect(ok).To(BeTrue())
				Expect(ds).To(Equal(inj.DataStore))
				return nil
			}

			inj.Accept(context.Background(), nil /* sink */, endpoint.InboundEnvelope{})

			Expect(next.AcceptCalls()).To(HaveLen(1))
		})
	})
})

var _ = Describe("OutboundInjector", func() {
	var (
		next *mocks.OutboundPipelineMock
		inj  *OutboundInjector
	)

	BeforeEach(func() {
		next = &mocks.OutboundPipelineMock{}
		inj = &OutboundInjector{
			DataStore: &mocks.DataStoreMock{},
			Next:      next,
		}
	})

	Describe("Initialize", func() {
		It("calls the next pipeline with the data store in the context", func() {
			next.InitializeFunc = func(ctx context.Context, _ *endpoint.Endpoint) error {
				ds, ok := GetDataStore(ctx)
				Expect(ok).To(BeTrue())
				Expect(ds).To(Equal(inj.DataStore))
				return nil
			}

			inj.Initialize(context.Background(), &endpoint.Endpoint{})

			Expect(next.InitializeCalls()).To(HaveLen(1))
		})
	})

	Describe("Accept", func() {
		It("calls the next pipeline with the data store in the context", func() {
			next.AcceptFunc = func(ctx context.Context, _ endpoint.OutboundEnvelope) error {
				ds, ok := GetDataStore(ctx)
				Expect(ok).To(BeTrue())
				Expect(ds).To(Equal(inj.DataStore))
				return nil
			}

			inj.Accept(context.Background(), endpoint.OutboundEnvelope{})

			Expect(next.AcceptCalls()).To(HaveLen(1))
		})
	})
})
