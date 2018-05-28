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
		It("calls the next pipeline", func() {
			next.InitializeFunc = func(ctx context.Context, _ *Endpoint) error {
				return nil
			}

			tl.Initialize(context.Background(), &Endpoint{})

			Expect(next.InitializeCalls()).To(HaveLen(1))
		})
	})

	Describe("Accept", func() {
		It("calls the next pipeline with the timeout set on the context", func() {
			next.AcceptFunc = func(ctx context.Context, _ MessageSink, _ InboundEnvelope) error {
				_, ok := ctx.Deadline()
				Expect(ok).To(BeTrue())
				return nil
			}

			tl.Accept(context.Background(), nil /* sink */, InboundEnvelope{})

			Expect(next.AcceptCalls()).To(HaveLen(1))
		})

		It("calls the next pipeline with the default timeout set on the context if none given", func() {
			tl = &TimeLimiter{
				Next: next,
			}

			next.AcceptFunc = func(ctx context.Context, _ MessageSink, _ InboundEnvelope) error {
				_, ok := ctx.Deadline()
				Expect(ok).To(BeTrue())
				return nil
			}

			tl.Accept(context.Background(), nil /* sink */, InboundEnvelope{})

			Expect(next.AcceptCalls()).To(HaveLen(1))
		})
	})
})
