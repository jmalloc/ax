package ax_test

import (
	. "github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("Envelope", func() {
	Describe("NewEnvelope", func() {
		message := &messagetest.Message{}
		env := NewEnvelope(message)

		It("returns an envelope containing the message", func() {
			Expect(env.Message).To(Equal(message))
		})

		It("generates a UUID message ID", func() {
			u, err := uuid.FromString(env.MessageID.Get())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(u.String()).To(Equal(env.MessageID.Get()))
		})

		It("sets the causation ID to the message ID", func() {
			Expect(env.CausationID).To(Equal(env.MessageID))
		})

		It("sets the correlation ID to the message ID", func() {
			Expect(env.CorrelationID).To(Equal(env.MessageID))
		})
	})

	Describe("NewChild", func() {
		rootMessage := &messagetest.Message{}
		root := NewEnvelope(rootMessage)
		branchMessage := &messagetest.Message{}
		branch := root.NewChild(branchMessage)
		leafMessage := &messagetest.Message{}
		leaf := branch.NewChild(leafMessage)

		It("returns an envelope containing the message", func() {
			Expect(leaf.Message).To(Equal(leafMessage))
		})

		It("generates a UUID message ID", func() {
			u, err := uuid.FromString(leaf.MessageID.Get())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(u.String()).To(Equal(leaf.MessageID.Get()))
		})

		It("sets the causation ID to the parent's message ID", func() {
			Expect(leaf.CausationID).To(Equal(branch.MessageID))
		})

		It("sets the correlation ID to the root's message ID", func() {
			Expect(leaf.CorrelationID).To(Equal(root.MessageID))
		})
	})
})
