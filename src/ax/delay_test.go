package ax_test

import (
	"time"

	. "github.com/jmalloc/ax/src/ax"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delay", func() {
	It("returns an option that delays sending", func() {
		env := Envelope{}
		t := time.Now()
		d := 10 * time.Second
		opt := Delay(d)

		err := opt.ApplyExecuteOption(&env)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(env.DelayedUntil).To(BeTemporally("~", t.Add(d)))
	})
})

var _ = Describe("DelayUntil", func() {
	It("returns an option that delays sending", func() {
		env := Envelope{}
		t := time.Now().Add(10 * time.Second)
		opt := DelayUntil(t)

		err := opt.ApplyExecuteOption(&env)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(env.DelayedUntil).To(BeTemporally("==", t))
	})
})
