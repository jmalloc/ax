package routing_test

import (
	"github.com/jmalloc/ax"
	. "github.com/jmalloc/ax/routing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EndpointTable", func() {
	var table EndpointTable

	Describe("NewEndpointTable", func() {
		It("returns an error when passed an odd number of arguments", func() {
			_, err := NewEndpointTable("foo")
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Lookup", func() {
		BeforeEach(func() {
			t, err := NewEndpointTable(
				"foo", "route:foo",
				"foo.qux", "route:foo.qux",
				"foo.bar.ExactMatch", "route:foo.bar.ExactMatch",
			)
			Expect(err).ShouldNot(HaveOccurred())
			table = t
		})

		It("favors an exact match", func() {
			ep, ok := table.Lookup(ax.MessageType{Name: "foo.bar.ExactMatch"})
			Expect(ok).To(BeTrue())
			Expect(ep).To(Equal("route:foo.bar.ExactMatch"))
		})

		It("returns the longest match when there is no exact match", func() {
			ep, ok := table.Lookup(ax.MessageType{Name: "foo.qux.Message"})
			Expect(ok).To(BeTrue())
			Expect(ep).To(Equal("route:foo.qux"))
		})

		Context("when there is no default route", func() {
			It("returns false for a message with no matching routes", func() {
				_, ok := table.Lookup(ax.MessageType{Name: "baz.qux.Message"})
				Expect(ok).To(BeFalse())
			})
		})

		Context("when there is a default route", func() {
			BeforeEach(func() {
				t, err := NewEndpointTable(
					"foo", "route:foo",
					"foo.qux", "route:foo.qux",
					"foo.bar.ExactMatch", "route:foo.bar.ExactMatch",
					"", "route:default",
				)
				Expect(err).ShouldNot(HaveOccurred())
				table = t
			})

			It("returns the default route for a message with no better matching routes", func() {
				ep, ok := table.Lookup(ax.MessageType{Name: "baz.qux.Message"})
				Expect(ok).To(BeTrue())
				Expect(ep).To(Equal("route:default"))
			})
		})
	})
})
