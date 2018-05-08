package persistence_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax/bus"
	. "github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/internal/bustest"
	"github.com/jmalloc/ax/src/internal/persistencetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Injector", func() {
	var (
		next *bustest.InboundPipelineMock
		inj  *Injector
	)

	BeforeEach(func() {
		next = &bustest.InboundPipelineMock{}
		inj = &Injector{
			DataStore: &persistencetest.DataStoreMock{},
			Next:      next,
		}
	})

	Describe("Initialize", func() {
		It("calls the next pipeline with the data store in the context", func() {
			next.InitializeFunc = func(ctx context.Context, _ bus.Transport) error {
				ds, ok := GetDataStore(ctx)
				Expect(ok).To(BeTrue())
				Expect(ds).To(Equal(inj.DataStore))
				return nil
			}

			inj.Initialize(context.Background(), nil /* transport */)

			Expect(next.InitializeCalls()).To(HaveLen(1))
		})
	})

	Describe("Accept", func() {
		It("calls the next pipeline with the data store in the context", func() {
			next.AcceptFunc = func(ctx context.Context, _ bus.MessageSink, _ bus.InboundEnvelope) error {
				ds, ok := GetDataStore(ctx)
				Expect(ok).To(BeTrue())
				Expect(ds).To(Equal(inj.DataStore))
				return nil
			}

			inj.Accept(context.Background(), nil /* sink */, bus.InboundEnvelope{})

			Expect(next.AcceptCalls()).To(HaveLen(1))
		})
	})
})
