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

var _ = Describe("SelfValidator", func() {

	Describe("Validate", func() {
		It("does not return an error if the message is valid", func() {
			v := SelfValidator{}
			err := v.Validate(context.Background(), &messagetest.Message{})
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("invokes Validate method on the message if it implements SelfValidatingMessage interface", func() {
			v := SelfValidator{}
			s := &endpointtest.SelfValidatingMessageMock{
				ValidateFunc: func() error {
					return nil
				},
			}
			err := v.Validate(context.Background(), s)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(s.ValidateCalls()).Should(HaveLen(1))
		})

		It("returns an error if SelfValidatingMessage.Validate method fails", func() {
			v := SelfValidator{}
			s := &endpointtest.SelfValidatingMessageMock{}
			expected := errors.New("self-validating message test error")
			s.ValidateFunc = func() error {
				return expected
			}
			err := v.Validate(context.Background(), s)
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(MatchError(expected))
		})
	})
})
