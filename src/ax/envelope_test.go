package ax_test

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	. "github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/axtest/testmessages"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("Envelope", func() {
	Describe("NewEnvelope", func() {
		message := &testmessages.Message{}
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

	Describe("NewEnvelopeFromProto", func() {
		It("returns an equivalent Envelope", func() {
			sendAtPB, err := ptypes.TimestampProto(
				time.Now().Add(1 * time.Minute),
			)
			Expect(err).ShouldNot(HaveOccurred())

			pb := &EnvelopeProto{
				MessageId:     "<message>",
				CausationId:   "<causation>",
				CorrelationId: "<correlation>",
				CreatedAt:     ptypes.TimestampNow(),
				SendAt:        sendAtPB,
			}

			m := &testmessages.Message{
				Value: "<message>",
			}

			pb.Message, err = ptypes.MarshalAny(m)
			Expect(err).ShouldNot(HaveOccurred())

			env, err := NewEnvelopeFromProto(pb)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(env.MessageID.String()).To(Equal("<message>"))
			Expect(env.CorrelationID.String()).To(Equal("<correlation>"))
			Expect(env.CausationID.String()).To(Equal("<causation>"))

			createdAt, err := ptypes.TimestampProto(env.CreatedAt)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createdAt).To(Equal(pb.CreatedAt))

			sendAt, err := ptypes.TimestampProto(env.SendAt)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sendAt).To(Equal(pb.SendAt))
		})
	})

	Describe("NewChild", func() {
		rootMessage := &testmessages.Message{}
		root := NewEnvelope(rootMessage)
		branchMessage := &testmessages.Message{}
		branch := root.NewChild(branchMessage)
		leafMessage := &testmessages.Message{}
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

	Describe("Type", func() {
		message := &testmessages.Message{}
		env := NewEnvelope(message)

		It("returns the type of the message in the envelope", func() {
			Expect(env.Type()).To(Equal(TypeOf(message)))
		})
	})

	Describe("Delay", func() {
		It("returns 0 when the timestamps are equal", func() {
			t := time.Now()

			env := Envelope{
				CreatedAt: t,
				SendAt:    t,
			}

			Expect(env.Delay()).To(Equal(time.Duration(0)))
		})

		It("returns 0 when the send-at timestamp is before the created-at timestamp", func() {
			t := time.Now()

			env := Envelope{
				CreatedAt: t,
				SendAt:    t.Add(-10),
			}

			Expect(env.Delay()).To(Equal(time.Duration(0)))
		})

		It("returns the delay send-at timestamp is after the created-at timestamp", func() {
			d := time.Duration(10)
			t := time.Now()

			env := Envelope{
				CreatedAt: t,
				SendAt:    t.Add(d),
			}

			Expect(env.Delay()).To(Equal(d))
		})
	})

	Describe("Equal", func() {
		message := &testmessages.Message{
			Value: "<message>",
		}

		var env1, env2 Envelope

		BeforeEach(func() {
			env1 = NewEnvelope(message)
			env2 = env1
			env2.Message = proto.Clone(env1.Message).(Message)
		})

		It("returns true when the envelopes are equal", func() {
			Expect(env1.Equal(env2)).To(BeTrue())
		})

		It("returns false when the message ID is different", func() {
			env2.MessageID = MustParseMessageID("different")
			Expect(env1.Equal(env2)).To(BeFalse())
		})

		It("returns false when the causation ID is different", func() {
			env2.CausationID = MustParseMessageID("different")
			Expect(env1.Equal(env2)).To(BeFalse())
		})

		It("returns false when the correlation ID is different", func() {
			env2.CorrelationID = MustParseMessageID("different")
			Expect(env1.Equal(env2)).To(BeFalse())
		})

		It("returns false when the created-at time is different", func() {
			env2.CreatedAt = env2.CreatedAt.Add(1 * time.Minute)
			Expect(env1.Equal(env2)).To(BeFalse())
		})

		It("returns false when the send-at time is different", func() {
			env2.SendAt = env2.SendAt.Add(1 * time.Minute)
			Expect(env1.Equal(env2)).To(BeFalse())
		})

		It("returns false when the message is different", func() {
			env2.Message = &testmessages.Message{
				Value: "<different>",
			}
			Expect(env1.Equal(env2)).To(BeFalse())
		})
	})

	Describe("AsProto", func() {
		It("returns an equivalent EnvelopeProto", func() {
			env := Envelope{
				MessageID:     MustParseMessageID("<message>"),
				CausationID:   MustParseMessageID("<causation>"),
				CorrelationID: MustParseMessageID("<correlation>"),
				CreatedAt:     time.Now(),
				SendAt:        time.Now().Add(1 * time.Minute),
				Message: &testmessages.Message{
					Value: "<message>",
				},
			}

			pb, err := env.AsProto()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(pb.MessageId).To(Equal("<message>"))
			Expect(pb.CorrelationId).To(Equal("<correlation>"))
			Expect(pb.CausationId).To(Equal("<causation>"))

			createdAt, err := ptypes.Timestamp(pb.CreatedAt)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createdAt).To(BeTemporally("==", env.CreatedAt))

			sendAt, err := ptypes.Timestamp(pb.SendAt)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sendAt).To(BeTemporally("==", env.SendAt))
		})
	})
})

var _ = Describe("MarshalEnvelope and UnmarshalEnvelope", func() {
	It("faithfully encodes and decodes the message", func() {
		env := Envelope{
			MessageID:     MustParseMessageID("<message>"),
			CausationID:   MustParseMessageID("<causation>"),
			CorrelationID: MustParseMessageID("<correlation>"),
			CreatedAt:     time.Now(),
			SendAt:        time.Now().Add(1 * time.Minute),
			Message: &testmessages.Message{
				Value: "<message>",
			},
		}

		buf, err := MarshalEnvelope(env)
		Expect(err).ShouldNot(HaveOccurred())

		res, err := UnmarshalEnvelope(buf)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(env.Equal(res)).To(BeTrue())
	})
})
