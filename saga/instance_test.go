package saga_test

import (
	"github.com/jmalloc/ax/ident"
	. "github.com/jmalloc/ax/saga"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("GenerateInstanceID", func() {
	It("generates a unique identifier for a message", func() {
		id := GenerateInstanceID()
		u, err := uuid.FromString(id.Get())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(u.String()).To(Equal(id.Get()))
	})
})

var _ = Describe("ParseInstanceID", func() {
	It("sets the ID to the parsed value", func() {
		id, err := ParseInstanceID("<id>")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(id.Get()).To(Equal("<id>"))
	})

	It("returns an error if the value is empty", func() {
		_, err := ParseInstanceID("")

		Expect(err).To(Equal(ident.ErrEmptyID))
	})

	It("returns an error if the ID is already set", func() {
		id := MustParseInstanceID("<id>")
		err := id.Parse("<id>")

		Expect(err).To(Equal(ident.ErrIDNotEmpty))
	})
})

var _ = Describe("MustParseInstanceID", func() {
	It("sets the ID to the parsed value", func() {
		id := MustParseInstanceID("<id>")

		Expect(id.Get()).To(Equal("<id>"))
	})

	It("panics if the value is empty", func() {
		Expect(func() {
			MustParseInstanceID("")
		}).To(Panic())
	})

	It("panics if the ID is already set", func() {
		id := MustParseInstanceID("<id>")

		Expect(func() {
			id.MustParse("<id>")
		}).To(Panic())
	})
})
