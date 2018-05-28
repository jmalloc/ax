package endpoint_test

import (
	"context"
	"time"

	. "github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/internal/endpointtest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimeLimiter", func() {
	var (
		next *endpointtest.InboundPipelineMock
		tl   *TimeLimiter
	)

	BeforeEach(func() {
		next = &endpointtest.InboundPipelineMock{}
		tl = &TimeLimiter{
			Timeout: 1 * time.Second,
			Next:    next,
		}
	})

	Describe("Initialize", func() {
		It("calls the next pipeline with the timeout set on the context", func() {
			ctxOriginal := context.Background()

			next.InitializeFunc = func(ctx context.Context, _ *Endpoint) error {
				Expect(ctx).To(Equal(ctxOriginal))
				return nil
			}

			tl.Initialize(ctxOriginal, &Endpoint{})

			Expect(next.InitializeCalls()).To(HaveLen(1))
		})
	})

	Describe("Accept", func() {
		It("calls the next pipeline with the timeout set on the context", func() {
			ctxOriginal := context.Background()

			next.AcceptFunc = func(ctx context.Context, _ MessageSink, _ InboundEnvelope) error {
				Expect(ctx).To(Not(Equal(ctxOriginal)))
				return nil
			}

			tl.Accept(ctxOriginal, nil /* sink */, InboundEnvelope{})

			Expect(next.AcceptCalls()).To(HaveLen(1))
		})

		It("calls the next pipeline with the default timeout set on the context", func() {
			tl = &TimeLimiter{
				Next: next,
			}

			ctxOriginal := context.Background()

			next.AcceptFunc = func(ctx context.Context, _ MessageSink, _ InboundEnvelope) error {
				Expect(ctx).To(Not(Equal(ctxOriginal)))
				return nil
			}

			tl.Accept(ctxOriginal, nil /* sink */, InboundEnvelope{})

			Expect(next.AcceptCalls()).To(HaveLen(1))
		})
	})
})
