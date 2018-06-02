package endpoint_test

import (
	"context"
	"errors"

	. "github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/internal/endpointtest"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Default Validators", func() {

	Describe("Basic Validator", func() {
		It("does not return an error if the message is valid", func() {
			v := BasicValidator{}
			err := v.Validate(context.Background(), &messagetest.Message{})
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns an instance of ValidatonError if the message is nil", func() {
			v := BasicValidator{}
			err := v.Validate(context.Background(), nil)
			Expect(err).Should(HaveOccurred())
			_, ok := err.(*ValidationError)
			Expect(ok).Should(BeTrue())
		})

		It("invokes Validate method on a message if it implements SelfValidatingMessage interface", func() {
			v := BasicValidator{}
			svm := &endpointtest.SelfValidatingMessageMock{
				ValidateFunc: func() error {
					return nil
				},
			}
			err := v.Validate(context.Background(), svm)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(svm.ValidateCalls()).Should(HaveLen(1))
		})

		It("returns an error if SelfValidatingMessage.Validate method fails", func() {
			v := BasicValidator{}
			expected := errors.New("self-validating message test error")
			svm := &endpointtest.SelfValidatingMessageMock{
				ValidateFunc: func() error {
					return expected
				},
			}
			err := v.Validate(context.Background(), svm)
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(MatchError(expected))
		})
	})
})
