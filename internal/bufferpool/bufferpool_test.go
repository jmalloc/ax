package bufferpool_test

import (
	"bytes"

	. "github.com/jmalloc/ax/internal/bufferpool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get", func() {
	It("returns a bytes.Buffer pointer", func() {
		buffer := Get()
		Expect(buffer).ShouldNot(BeNil())
	})

	It("recycles buffers", func() {
		buffer := Get()
		Put(buffer)

		Expect(Get()).To(Equal(buffer))
	})
})

var _ = Describe("Put", func() {
	It("accepts a buffer pointer", func() {
		var buffer bytes.Buffer
		Put(&buffer)
	})

	It("accepts a nil pointer", func() {
		var buffer *bytes.Buffer
		Put(buffer)

		Expect(Get()).ShouldNot(BeNil())
	})
})
