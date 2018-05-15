package endpoint_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	. "github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithEnvelope / GetEnvelope", func() {
	It("transports a message envelope via the context", func() {
		expected := ax.NewEnvelope(&messagetest.Message{})
		ctx := WithEnvelope(context.Background(), expected)

		env, ok := GetEnvelope(ctx)

		Expect(ok).To(BeTrue())
		Expect(env).To(BeIdenticalTo(expected))
	})
})
