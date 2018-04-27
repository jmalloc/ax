package persistence_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax/bus"
	. "github.com/jmalloc/ax/src/ax/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type nextPipelineStage struct {
	Context context.Context
}

func (n *nextPipelineStage) Initialize(ctx context.Context, t bus.Transport) error {
	n.Context = ctx
	return nil
}

func (n *nextPipelineStage) DeliverMessage(ctx context.Context, s bus.MessageSender, m bus.InboundEnvelope) error {
	n.Context = ctx
	return nil
}

var _ = Describe("Injector", func() {
	next := &nextPipelineStage{}
	inj := &Injector{
		DataStore: &dataStore{},
		Next:      next,
	}

	BeforeEach(func() {
		next.Context = nil
	})

	Describe("Initialize", func() {
		It("calls the next pipeline with the data store in the context", func() {
			inj.Initialize(context.Background(), nil /* transport */)

			ds, ok := GetDataStore(next.Context)
			Expect(ok).To(BeTrue())
			Expect(ds).To(Equal(inj.DataStore))
		})
	})

	Describe("DeliverMessage", func() {
		It("calls the next pipeline with the data store in the context", func() {
			inj.DeliverMessage(context.Background(), nil /* sender */, bus.InboundEnvelope{})

			ds, ok := GetDataStore(next.Context)
			Expect(ok).To(BeTrue())
			Expect(ds).To(Equal(inj.DataStore))
		})
	})
})
