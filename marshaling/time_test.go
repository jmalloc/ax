package marshaling_test

import (
	"time"

	. "github.com/jmalloc/ax/marshaling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MarshalTime", func() {
	It("returns the time in RFC3339Nano format", func() {
		t := time.Now()

		Expect(MarshalTime(t)).To(Equal(
			t.Format(time.RFC3339Nano),
		))
	})
})

var _ = Describe("UnmarshalTime", func() {
	It("unmarshals from RFC3339Nano format", func() {
		t := time.Now()
		s := t.Format(time.RFC3339Nano)

		var v time.Time
		err := UnmarshalTime(s, &v)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(v.Equal(t)).To(BeTrue())
	})

	It("returns an error if the string is not a valid time", func() {
		t := time.Now()
		s := t.Format(time.RFC1123) // note wrong format

		var v time.Time
		err := UnmarshalTime(s, &v)

		Expect(err).Should(HaveOccurred())
	})
})
