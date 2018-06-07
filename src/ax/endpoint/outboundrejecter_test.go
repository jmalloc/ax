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

var _ = Describe("OutboundRejecter", func() {
	var (
		next                               *mocks.OutboundPipelineMock
		validator1, validator2, validator3 *mocks.ValidatorMock
	)

	BeforeEach(func() {
		next = &mocks.OutboundPipelineMock{
			InitializeFunc: func(
				context.Context,
				*Endpoint,
			) error {
				return nil
			},
			AcceptFunc: func(
				ctx context.Context,
				env OutboundEnvelope,
			) error {
				return nil
			},
		}
		validator1 = &mocks.ValidatorMock{
			ValidateFunc: func(ctx context.Context, m ax.Message) error {
				return nil
			},
		}
		validator2 = &mocks.ValidatorMock{
			ValidateFunc: func(ctx context.Context, m ax.Message) error {
				return nil
			},
		}
		validator3 = &mocks.ValidatorMock{
			ValidateFunc: func(ctx context.Context, m ax.Message) error {
				return nil
			},
		}
	})

	Describe("Initialize", func() {
		ep := &Endpoint{}

		It("initializes the following stage", func() {
			os := &OutboundRejecter{
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
			next.InitializeFunc = func(context.Context, *Endpoint) error { return expected }

			os := &OutboundRejecter{
				Next:       next,
				Validators: []Validator{},
			}

			err := os.Initialize(context.Background(), &Endpoint{})
			Expect(err).To(Equal(expected))
		})
	})

	Describe("Accept", func() {
		env := OutboundEnvelope{
			Envelope: ax.NewEnvelope(
				&testmessages.Message{},
			),
		}

		It("passes message to next stage's Accept method if none of validators fails", func() {
			os := &OutboundRejecter{
				Next: next,
				Validators: []Validator{
					validator1, validator2, validator3,
				},
			}

			err := os.Accept(context.Background(), env)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(next.AcceptCalls())).Should(BeNumerically("==", 1))
		})

		It("returns an error if any of validator fails", func() {
			expected := errors.New("<error>")

			os := &OutboundRejecter{
				Next: next,
				Validators: []Validator{
					validator1, validator2, validator3,
				},
			}

			validator2.ValidateFunc = func(ctx context.Context, m ax.Message) error {
				return expected
			}

			err := os.Accept(context.Background(), env)
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(Equal(expected))
			Expect(len(next.AcceptCalls())).Should(BeNumerically("==", 0))
		})
	})
})
