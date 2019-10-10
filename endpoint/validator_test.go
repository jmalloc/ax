package endpoint_test

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/axtest/mocks"
	"github.com/jmalloc/ax/axtest/testmessages"
	. "github.com/jmalloc/ax/endpoint"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SelfValidator", func() {

	Describe("Validate", func() {
		It("does not return an error if the message is valid", func() {
			v := SelfValidator{}
			err := v.Validate(context.Background(), &testmessages.Message{})
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("invokes Validate method on the message if it implements SelfValidatingMessage interface", func() {
			v := SelfValidator{}
			s := &mocks.SelfValidatingMessageMock{
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
			s := &mocks.SelfValidatingMessageMock{}
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
