package endpoint_test

import (
	. "github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidationError", func() {
	It("returns a valid pointer when initializing a new validation error", func() {
		err := NewValidationError(
			"test validation error",
			&messagetest.Message{},
		)
		Expect(err).ShouldNot(BeNil())
	})
	It("returns a correct error message containing original error message", func() {
		exected := "test validation error"
		err := NewValidationError(
			exected,
			&messagetest.Message{},
		)
		Expect(err.Error()).Should(ContainSubstring(exected))
	})
})
