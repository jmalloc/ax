package endpoint_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
	. "github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/axtest/mocks"
	"github.com/jmalloc/ax/src/axtest/testmessages"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fanout", func() {
	var (
		next1, next2, next3 *mocks.InboundPipelineMock
		fanout              Fanout
	)

	BeforeEach(func() {
		next1 = &mocks.InboundPipelineMock{
			InitializeFunc: func(context.Context, *Endpoint) error { return nil },
			AcceptFunc:     func(context.Context, MessageSink, InboundEnvelope) error { return nil },
		}
		next2 = &mocks.InboundPipelineMock{
			InitializeFunc: func(context.Context, *Endpoint) error { return nil },
			AcceptFunc:     func(context.Context, MessageSink, InboundEnvelope) error { return nil },
		}
		next3 = &mocks.InboundPipelineMock{
			InitializeFunc: func(context.Context, *Endpoint) error { return nil },
			AcceptFunc:     func(context.Context, MessageSink, InboundEnvelope) error { return nil },
		}

		fanout = Fanout{next1, next2, next3}
	})

	Describe("Initialize", func() {
		ep := &Endpoint{}

		It("initializes all of the following stages", func() {
			err := fanout.Initialize(context.Background(), ep)
			Expect(err).ShouldNot(HaveOccurred())

			for _, p := range []*mocks.InboundPipelineMock{next1, next2, next3} {
				Expect(p.InitializeCalls()).To(HaveLen(1))

				// don't compare ctx
				args := p.InitializeCalls()[0]
				Expect(args.Ep).To(Equal(ep))
			}
		})

		It("returns an error if any of the next stages fails", func() {
			expected := errors.New("<error>")
			next2.InitializeFunc = func(context.Context, *Endpoint) error { return expected }

			err := fanout.Initialize(context.Background(), &Endpoint{})
			Expect(err).To(Equal(expected))
		})
	})

	Describe("Accept", func() {
		sink := &BufferedSink{}
		env := InboundEnvelope{
			Envelope: ax.NewEnvelope(
				&testmessages.Message{},
			),
		}

		It("forwards the message to all of the following stages", func() {
			err := fanout.Accept(context.Background(), sink, env)
			Expect(err).ShouldNot(HaveOccurred())

			for _, p := range []*mocks.InboundPipelineMock{next1, next2, next3} {
				Expect(p.AcceptCalls()).To(HaveLen(1))

				// don't compare ctx
				args := p.AcceptCalls()[0]
				Expect(args.Sink).To(Equal(sink))
				Expect(args.Env).To(Equal(env))
			}
		})

		It("returns an error if any of the next stages fail", func() {
			expected := errors.New("<error>")
			next2.AcceptFunc = func(context.Context, MessageSink, InboundEnvelope) error { return expected }

			err := fanout.Accept(context.Background(), sink, env)
			Expect(err).To(Equal(expected))
		})
	})
})
