package ax_test

import (
	"reflect"

	. "github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageType", func() {
	message := TypeOf(&messagetest.Message{})

	Describe("TypeOf", func() {
		It("returns a message type with the correct name", func() {
			Expect(message.Name).To(Equal("ax.internal.messagetest.Message"))
		})

		It("returns a message type with the correct struct type", func() {
			Expect(message.StructType).To(Equal(reflect.TypeOf(messagetest.Message{})))
		})
	})

	Describe("TypeByName", func() {
		It("returns a message type with the correct name", func() {
			mt, ok := TypeByName("ax.internal.messagetest.Message")
			Expect(ok).To(BeTrue())
			Expect(mt.Name).To(Equal("ax.internal.messagetest.Message"))
		})

		It("returns a message type with the correct struct type", func() {
			mt, ok := TypeByName("ax.internal.messagetest.Message")
			Expect(ok).To(BeTrue())
			Expect(mt.StructType).To(Equal(reflect.TypeOf(messagetest.Message{})))
		})

		It("returns false if the message name is not registered", func() {
			_, ok := TypeByName("ax.internal.messagetest.Unknown")
			Expect(ok).To(BeFalse())
		})
	})

	Context("when the message is generic", func() {
		Describe("IsCommand", func() {
			It("returns false", func() {
				Expect(message.IsCommand()).To(BeFalse())
			})
		})

		Describe("IsEvent", func() {
			It("returns false", func() {
				Expect(message.IsEvent()).To(BeFalse())
			})
		})
	})

	Context("when the message is a command", func() {
		command := TypeOf(&messagetest.Command{})

		Describe("IsCommand", func() {
			It("returns true", func() {
				Expect(command.IsCommand()).To(BeTrue())
			})
		})

		Describe("IsEvent", func() {
			It("returns false", func() {
				Expect(command.IsEvent()).To(BeFalse())
			})
		})
	})

	Context("when the message is an event", func() {
		event := TypeOf(&messagetest.Event{})

		Describe("IsCommand", func() {
			It("returns false", func() {
				Expect(event.IsCommand()).To(BeFalse())
			})
		})

		Describe("IsEvent", func() {
			It("returns true", func() {
				Expect(event.IsEvent()).To(BeTrue())
			})
		})
	})

	Describe("ToSet", func() {
		It("returns a set containing only this message type", func() {
			Expect(message.ToSet()).To(Equal(TypesOf(&messagetest.Message{})))
		})
	})

	Describe("New", func() {
		It("returns a pointer to a new instance of the message struct", func() {
			Expect(message.New()).To(Equal(&messagetest.Message{}))
		})
	})

	Describe("PackageName", func() {
		It("returns the protocol buffers package name", func() {
			Expect(message.PackageName()).To(Equal("ax.internal.messagetest"))
		})

		It("returns an empty string if the message is not in a package", func() {
			mt := TypeOf(&messagetest.NoPackage{})
			Expect(mt.PackageName()).To(Equal(""))
		})
	})
})

var _ = Describe("MessageTypeSet", func() {
	message := TypeOf(&messagetest.Message{})
	command := TypeOf(&messagetest.Command{})
	event := TypeOf(&messagetest.Event{})

	Describe("NewMessageTypeSet", func() {
		It("returns a set containing the the arguments", func() {
			Expect(
				NewMessageTypeSet(
					message,
					command,
				).Members(),
			).To(ConsistOf(
				message,
				command,
			))
		})

		It("deduplicates repeated types", func() {
			Expect(
				NewMessageTypeSet(
					message,
					message,
				).Len(),
			).To(Equal(1))
		})

		It("returns an empty set when called with no arguments", func() {
			Expect(NewMessageTypeSet().Len()).To(Equal(0))
		})
	})

	Describe("TypesOf", func() {
		It("returns a set containing the message types of the arguments", func() {
			Expect(
				TypesOf(
					&messagetest.Message{},
					&messagetest.Command{},
				).Members(),
			).To(ConsistOf(
				message,
				command,
			))
		})

		It("deduplicates repeated types", func() {
			Expect(
				TypesOf(
					&messagetest.Message{},
					&messagetest.Message{},
				).Len(),
			).To(Equal(1))
		})

		It("returns an empty set when called with no arguments", func() {
			Expect(TypesOf().Len()).To(Equal(0))
		})
	})

	Describe("Has", func() {
		set := TypesOf(&messagetest.Message{})

		It("returns true if the message type is a member of the set", func() {
			Expect(set.Has(message)).To(BeTrue())
		})

		It("returns false if the message type is a member of the set", func() {
			Expect(set.Has(command)).To(BeFalse())
		})
	})

	Describe("Add", func() {
		set := TypesOf(&messagetest.Message{})

		It("returns a set containing the message type", func() {
			Expect(set.Add(command)).To(Equal(
				NewMessageTypeSet(
					message,
					command,
				),
			))
		})

		It("does not modify the original set", func() {
			set.Add(command)

			Expect(set.Members()).To(ConsistOf(message))
		})

		It("returns the original set if the message type is already a member", func() {
			Expect(set.Add(message)).To(Equal(set))
		})
	})

	Describe("Add", func() {
		setA := TypesOf(
			&messagetest.Message{},
			&messagetest.Command{},
		)

		setB := TypesOf(
			&messagetest.Command{},
			&messagetest.Event{},
		)

		It("returns the union of two sets", func() {
			set := setA.Union(setB)

			Expect(set.Members()).To(ConsistOf(
				message,
				command,
				event,
			))
		})

		It("does not modify the original sets", func() {
			setA.Union(setB)

			Expect(setA.Members()).To(ConsistOf(message, command))
			Expect(setB.Members()).To(ConsistOf(command, event))
		})

		It("returns the LHS if the RHS is empty", func() {
			Expect(setA.Union(NewMessageTypeSet())).To(Equal(setA))
		})

		It("returns the RHS if the LHS is empty", func() {
			Expect(NewMessageTypeSet().Union(setA)).To(Equal(setA))
		})
	})
})
