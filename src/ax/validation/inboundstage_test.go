package validation_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/internal/validationtest"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	. "github.com/jmalloc/ax/src/ax/validation"

	"github.com/jmalloc/ax/src/internal/messagetest"

	"github.com/jmalloc/ax/src/internal/endpointtest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InboundStage", func() {
	var (
		next                               *endpointtest.InboundPipelineMock
		validator1, validator2, validator3 *validationtest.ValidatorMock
	)

	BeforeEach(func() {
		next = &endpointtest.InboundPipelineMock{
			InitializeFunc: func(
				context.Context,
				*endpoint.Endpoint,
			) error {
				return nil
			},
			AcceptFunc: func(
				context.Context,
				endpoint.MessageSink,
				endpoint.InboundEnvelope,
			) error {
				return nil
			},
		}
		validator1 = &validationtest.ValidatorMock{
			ValidateFunc: func(ctx context.Context, msg ax.Message) error {
				return nil
			},
		}
		validator2 = &validationtest.ValidatorMock{
			ValidateFunc: func(ctx context.Context, msg ax.Message) error {
				return nil
			},
		}
		validator3 = &validationtest.ValidatorMock{
			ValidateFunc: func(ctx context.Context, msg ax.Message) error {
				return nil
			},
		}
	})

	Describe("Initialize", func() {
		ep := &endpoint.Endpoint{}

		It("initializes the following stage", func() {
			os := &InboundStage{
				Next: next,
				Validators: []Validator{
					validator1, validator2, validator3,
				},
			}

			err := os.Initialize(context.Background(), ep)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(next.InitializeCalls()).To(HaveLen(1))

			// don't compare ctx
			args := next.InitializeCalls()[0]
			Expect(args.Ep).To(Equal(ep))
		})

		It("returns an error if the initialization of the next stage fails", func() {
			expected := errors.New("<error>")
			next.InitializeFunc = func(context.Context, *endpoint.Endpoint) error { return expected }

			os := &InboundStage{
				Next:       next,
				Validators: []Validator{},
			}

			err := os.Initialize(context.Background(), &endpoint.Endpoint{})
			Expect(err).To(Equal(expected))
		})
	})

	Describe("Accept", func() {
		sink := &endpoint.BufferedSink{}
		env := endpoint.InboundEnvelope{
			Envelope: ax.NewEnvelope(
				&messagetest.Message{},
			),
		}

		It("passes message to next stage's Accept method if none of validators fails", func() {
			os := &InboundStage{
				Next: next,
				Validators: []Validator{
					validator1, validator2, validator3,
				},
			}

			err := os.Accept(context.Background(), sink, env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(next.AcceptCalls())).Should(BeNumerically("==", 1))
		})

		It("returns an error if any of validator fails", func() {
			expected := errors.New("<error>")

			os := &InboundStage{
				Next: next,
				Validators: []Validator{
					validator1, validator2, validator3,
				},
			}

			validator2.ValidateFunc = func(ctx context.Context, msg ax.Message) error {
				return expected
			}

			err := os.Accept(context.Background(), sink, env)
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(Equal(expected))
			Expect(len(next.AcceptCalls())).Should(BeNumerically("==", 0))
		})
	})
})
