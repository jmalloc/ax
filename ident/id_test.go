package ident_test

import (
	. "github.com/jmalloc/ax/ident"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("ID", func() {
	var id ID

	BeforeEach(func() {
		id = ID{}
	})

	Describe("GenerateUUID", func() {
		It("sets the ID to a new random UUID", func() {
			id.GenerateUUID()

			u, err := uuid.FromString(id.Get())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(u.String()).To(Equal(id.Get()))
		})

		It("panics if the ID is already set", func() {
			id.MustParse("<id>")

			Expect(func() {
				id.GenerateUUID()
			}).To(Panic())
		})
	})

	Describe("Parse", func() {
		It("sets the ID to the parsed value", func() {
			err := id.Parse("<id>")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(id.Get()).To(Equal("<id>"))
		})

		It("returns an error if the value is empty", func() {
			err := id.Parse("")

			Expect(err).To(Equal(ErrEmptyID))
		})

		It("returns an error if the ID is already set", func() {
			id.MustParse("<id>")
			err := id.Parse("<id>")

			Expect(err).To(Equal(ErrIDNotEmpty))
		})
	})

	Describe("MustParse", func() {
		It("sets the ID to the parsed value", func() {
			id.MustParse("<id>")

			Expect(id.Get()).To(Equal("<id>"))
		})

		It("panics if the value is empty", func() {
			Expect(func() {
				id.MustParse("")
			}).To(Panic())
		})

		It("panics if the ID is already set", func() {
			id.MustParse("<id>")

			Expect(func() {
				id.MustParse("<id>")
			}).To(Panic())
		})
	})

	Describe("Get", func() {
		It("returns the ID as string", func() {
			id.MustParse("<id>")

			Expect(id.Get()).To(Equal("<id>"))
		})

		It("panics if the ID has not been set", func() {
			Expect(func() {
				id.Get()
			}).To(Panic())
		})
	})

	Describe("String", func() {
		It("returns the ID as string", func() {
			id.MustParse("<id>")

			Expect(id.String()).To(Equal("<id>"))
		})

		It("uses short formatting for UUIDs", func() {
			id.MustParse("7fd0fd54-28bc-49f2-8916-8128d8cfc69e")

			Expect(id.String()).To(Equal("7fd0fd54"))
		})

		It("returns a fixed string if the ID has not been set", func() {
			Expect(id.String()).To(Equal("<unidentified>"))
		})
	})

	Describe("Validate", func() {
		It("returns nil if the ID is set", func() {
			id.MustParse("<id>")
			err := id.Validate()

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns an error if the ID is not set", func() {
			err := id.Validate()

			Expect(err).To(Equal(ErrEmptyID))
		})
	})

	Describe("MustValidate", func() {
		It("does not panic if the ID is set", func() {
			id.MustParse("<id>")
			id.MustValidate()
		})

		It("panics if the ID is not set", func() {
			Expect(func() {
				id.MustValidate()
			}).To(Panic())
		})
	})

	Context("JSON marshaling", func() {
		Describe("MarshalJSON", func() {
			It("returns the ID as a JSON string", func() {
				id.MustParse("<id>")
				buf, err := id.MarshalJSON()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buf).To(Equal([]byte(`"\u003cid\u003e"`)))
			})

			It("returns an empty JSON string if the ID is not set", func() {
				buf, err := id.MarshalJSON()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buf).To(Equal([]byte(`""`)))
			})
		})

		Describe("UnmarshalJSON", func() {
			It("sets the ID to the value of a JSON string", func() {
				err := id.UnmarshalJSON([]byte(`"\u003cid\u003e"`))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(id.Get()).To(Equal("<id>"))
			})

			It("allows an empty JSON string to be unmarshaled", func() {
				err := id.UnmarshalJSON([]byte(`""`))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(id.Validate()).To(Equal(ErrEmptyID))
			})
		})
	})

	Context("text marshaling", func() {
		Describe("MarshalText", func() {
			It("returns the ID as a string", func() {
				id.MustParse("<id>")
				buf, err := id.MarshalText()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buf).To(Equal([]byte(`<id>`)))
			})

			It("returns an empty string if the ID is not set", func() {
				buf, err := id.MarshalText()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(buf).To(Equal([]byte{}))
			})
		})

		Describe("UnmarshalText", func() {
			It("sets the ID to the value of the buffer", func() {
				err := id.UnmarshalText([]byte(`<id>`))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(id.Get()).To(Equal("<id>"))
			})

			It("allows an empty string to be unmarshaled", func() {
				err := id.UnmarshalText([]byte{})

				Expect(err).ShouldNot(HaveOccurred())
				Expect(id.Validate()).To(Equal(ErrEmptyID))
			})
		})
	})

	Context("SQL marshaling", func() {
		Describe("Value", func() {
			It("returns the ID as a string", func() {
				id.MustParse("<id>")
				v, err := id.Value()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(v).To(Equal("<id>"))
			})

			It("returns nil if the ID is not set", func() {
				v, err := id.Value()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(v).To(BeNil())
			})
		})

		Describe("Scan", func() {
			DescribeTable(
				"sets the ID to the scanned value",
				func(v interface{}, expected string) {
					err := id.Scan(v)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(id.Get()).To(Equal(expected))
				},
				Entry("string", "<id>", "<id>"),
				Entry("byte-slice", []byte("<id>"), "<id>"),
				Entry("integer", int64(123), "123"),
			)

			DescribeTable(
				"allows unmarshaling from empty values",
				func(v interface{}) {
					err := id.Scan(v)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(id.Validate()).To(Equal(ErrEmptyID))
				},
				Entry("empty string", ""),
				Entry("empty byte-slice", []byte{}),
				Entry("nil", nil),
			)

			It("returns an error if the type is not supported", func() {
				err := id.Scan(3.14)

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
