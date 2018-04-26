package ident_test

import (
	. "github.com/jmalloc/ax/src/ax/ident"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Format", func() {
	Context("when the identifier is formatted like a UUID", func() {
		It("returns the portion up to the first hyphen", func() {
			s := FormatID("7fd0fd54-28bc-49f2-8916-8128d8cfc69e")
			Expect(s).To(Equal("7fd0fd54"))
		})
	})

	Context("when the identifier is not formatted like a UUID", func() {
		It("returns the entire identifier", func() {
			s := FormatID("<identifier>")
			Expect(s).To(Equal("<identifier>"))
		})
	})
})
