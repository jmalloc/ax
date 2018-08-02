package messagestoretests

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/messagestore"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axtest/testmessages"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// MessageStoreSuite returns a test suite for implementations of messagestore.GloballyOrderedStore,
func MessageStoreSuite(
	getStore func() persistence.DataStore,
	getMessageStore func() messagestore.GloballyOrderedStore,
) func() {
	return func() {
		var (
			causationID, correlationID ax.MessageID
			store                      persistence.DataStore
			msgStore                   messagestore.GloballyOrderedStore
			ctx                        context.Context
			cancel                     func()
		)

		g.BeforeEach(func() {
			store = getStore()
			msgStore = getMessageStore()
			var fn func()
			ctx, fn = context.WithTimeout(context.Background(), 15*time.Second)
			cancel = fn // defeat go vet warning about unused cancel func

			causationID = ax.GenerateMessageID()
			correlationID = ax.GenerateMessageID()
		})

		g.AfterEach(func() {
			cancel()
		})

		g.Describe("AppendMessages", func() {
			var m1, m2 ax.Envelope
			var t1, t2, t3, t4 time.Time
			var s1, s2 string
			g.BeforeEach(func() {
				t1 = time.Now()
				t2 = time.Now()
				t3 = time.Now()
				t4 = time.Now()
				s1 = "<stream1>"
				s2 = "<stream4>"

				m1 = ax.Envelope{
					MessageID:     ax.GenerateMessageID(),
					CausationID:   causationID,
					CorrelationID: correlationID,
					CreatedAt:     t1,
					SendAt:        t2,
					Message: &testmessages.Command{
						Value: "<foo>",
					},
				}
				m2 = ax.Envelope{
					MessageID:     ax.GenerateMessageID(),
					CausationID:   causationID,
					CorrelationID: correlationID,
					CreatedAt:     t3,
					SendAt:        t4,
					Message: &testmessages.Event{
						Value: "<bar>",
					},
				}
			})
			g.Context("when the offset is the next unused offset in the stream", func() {
				g.It("returns no error", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					offset := uint64(0)

					err = msgStore.AppendMessages(
						ctx,
						tx,
						s1,
						offset,
						[]ax.Envelope{m1, m2},
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					err = msgStore.AppendMessages(
						ctx,
						tx,
						s2,
						offset,
						[]ax.Envelope{m1, m2},
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					err = com.Commit()
					m.Expect(err).ShouldNot(m.HaveOccurred())
				})
			})
			g.Context("when the offset is not the next unused offset in the stream", func() {
				g.It("returns an error", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					offset := uint64(100)
					err = msgStore.AppendMessages(
						ctx,
						tx,
						s1,
						offset,
						[]ax.Envelope{m1, m2},
					)
					m.Expect(err).Should(m.HaveOccurred())

					err = com.Commit()
					m.Expect(err).ShouldNot(m.HaveOccurred())
				})
			})
		})
	}
}
