package persistence_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax/endpoint"
	. "github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axtest/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Injector", func() {
	var (
		next *mocks.InboundPipelineMock
		inj  *Injector
	)

	BeforeEach(func() {
		next = &mocks.InboundPipelineMock{}
		inj = &Injector{
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
