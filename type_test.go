package ax_test

import (
	"reflect"

	. "github.com/jmalloc/ax"
	"github.com/jmalloc/ax/axtest/testmessages"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageType", func() {
	message := TypeOf(&testmessages.Message{})

	Describe("TypeOf", func() {
		It("returns a message type with the correct name", func() {
			Expect(message.Name).To(Equal("axtest.testmessages.Message"))
		})

		It("returns a message type with the correct struct type", func() {
			Expect(message.StructType).To(Equal(reflect.TypeOf(testmessages.Message{})))
		})
	})

	Describe("TypeByName", func() {
		It("returns a message type with the correct name", func() {
			mt, ok := TypeByName("axtest.testmessages.Message")
			Expect(ok).To(BeTrue())
			Expect(mt.Name).To(Equal("axtest.testmessages.Message"))
		})

		It("returns a message type with the correct struct type", func() {
			mt, ok := TypeByName("axtest.testmessages.Message")
			Expect(ok).To(BeTrue())
			Expect(mt.StructType).To(Equal(reflect.TypeOf(testmessages.Message{})))
		})

		It("returns false if the message name is not registered", func() {
			_, ok := TypeByName("axtest.testmessages.Unknown")
			Expect(ok).To(BeFalse())
		})

		It("returns false if the message name is registered, but the message type is not of ax.Message", func() {
			_, ok := TypeByName("axtest.testmessages.NonAxMessage")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("TypeByGoType", func() {
		It("returns the correct message type", func() {
			mt := TypeByGoType(reflect.TypeOf(&testmessages.Message{}))
			Expect(mt).To(Equal(TypeOf(&testmessages.Message{})))
		})

		It("panics if the type is not a message", func() {
			Expect(func() {
				TypeByGoType(reflect.TypeOf(&testmessages.NonAxMessage{}))
			}).To(Panic())
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
		command := TypeOf(&testmessages.Command{})

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
		event := TypeOf(&testmessages.Event{})

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
			Expect(message.ToSet()).To(Equal(TypesOf(&testmessages.Message{})))
		})
	})

	Describe("New", func() {
		It("returns a pointer to a new instance of the message struct", func() {
			Expect(message.New()).To(Equal(&testmessages.Message{}))
		})
	})

	Describe("PackageName", func() {
		It("returns the protocol buffers message name", func() {
			Expect(message.MessageName()).To(Equal("Message"))
		})

		It("returns the message name if the message is not in a package", func() {
			mt := TypeOf(&testmessages.NoPackage{})
			Expect(mt.MessageName()).To(Equal("NoPackage"))
		})
	})

	Describe("PackageName", func() {
		It("returns the protocol buffers package name", func() {
			Expect(message.PackageName()).To(Equal("axtest.testmessages"))
		})

		It("returns an empty string if the message is not in a package", func() {
			mt := TypeOf(&testmessages.NoPackage{})
			Expect(mt.PackageName()).To(Equal(""))
		})
	})

	Describe("String", func() {
		It("suffixes a question mark on commands", func() {
			mt := TypeOf(&testmessages.Command{})
			Expect(mt.String()).To(Equal("axtest.testmessages.Command?"))
		})

		It("suffixes an exclamation mark on events", func() {
			mt := TypeOf(&testmessages.Event{})
			Expect(mt.String()).To(Equal("axtest.testmessages.Event!"))
		})

		It("does not add a suffix to generic messages", func() {
			mt := TypeOf(&testmessages.Message{})
			Expect(mt.String()).To(Equal("axtest.testmessages.Message"))
		})
	})
})

var _ = Describe("MessageTypeSet", func() {
	message := TypeOf(&testmessages.Message{})
	command := TypeOf(&testmessages.Command{})
	event := TypeOf(&testmessages.Event{})

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
					&testmessages.Message{},
					&testmessages.Command{},
				).Members(),
			).To(ConsistOf(
				message,
				command,
			))
		})

		It("deduplicates repeated types", func() {
			Expect(
				TypesOf(
					&testmessages.Message{},
					&testmessages.Message{},
				).Len(),
			).To(Equal(1))
		})

		It("returns an empty set when called with no arguments", func() {
			Expect(TypesOf().Len()).To(Equal(0))
		})
	})

	Describe("TypesByGoType", func() {
		It("returns a set containing the message types of the arguments", func() {
			Expect(
				TypesByGoType(
					reflect.TypeOf(&testmessages.Message{}),
					reflect.TypeOf(&testmessages.Command{}),
				).Members(),
			).To(ConsistOf(
				message,
				command,
			))
		})

		It("deduplicates repeated types", func() {
			Expect(
				TypesByGoType(
					reflect.TypeOf(&testmessages.Message{}),
					reflect.TypeOf(&testmessages.Message{}),
				).Len(),
			).To(Equal(1))
		})

		It("returns an empty set when called with no arguments", func() {
			Expect(TypesByGoType().Len()).To(Equal(0))
		})

		It("panics if any of the types if not a message", func() {
			Expect(func() {
				TypesByGoType(
					reflect.TypeOf(&testmessages.Message{}),
					reflect.TypeOf(&testmessages.NonAxMessage{}),
				)
			}).To(Panic())
		})
	})

	Describe("Has", func() {
		set := TypesOf(&testmessages.Message{})

		It("returns true if the message type is a member of the set", func() {
			Expect(set.Has(message)).To(BeTrue())
		})

		It("returns false if the message type is a member of the set", func() {
			Expect(set.Has(command)).To(BeFalse())
		})
	})

	Describe("Add", func() {
		set := TypesOf(&testmessages.Message{})

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

	Describe("Union", func() {
		setA := TypesOf(
			&testmessages.Message{},
			&testmessages.Command{},
		)

		setB := TypesOf(
			&testmessages.Command{},
			&testmessages.Event{},
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

	Describe("Intersection", func() {
		setA := TypesOf(
			&testmessages.Message{},
			&testmessages.Command{},
		)

		setB := TypesOf(
			&testmessages.Command{},
			&testmessages.Event{},
		)

		It("returns the intersection of two sets", func() {
			set := setA.Intersection(setB)

			Expect(set.Members()).To(ConsistOf(
				command,
			))
		})

		It("does not modify the original sets", func() {
			setA.Intersection(setB)

			Expect(setA.Members()).To(ConsistOf(message, command))
			Expect(setB.Members()).To(ConsistOf(command, event))
		})

		It("returns the LHS if it is empty", func() {
			lhs := NewMessageTypeSet()
			Expect(lhs.Intersection(setA)).To(Equal(lhs))
		})

		It("returns the RHS if it is empty", func() {
			rhs := NewMessageTypeSet()
			Expect(setA.Intersection(rhs)).To(Equal(rhs))
		})
	})
})
