package delayedmessagetests

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/axtest"
	"github.com/jmalloc/ax/src/axtest/testmessages"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/delayedmessage"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/persistence"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// RepositorySuite returns a test suite for implementations of
// delayedmessage.Repository.
func RepositorySuite(
	getStore func() persistence.DataStore,
	getRepo func() delayedmessage.Repository,
) func() {
	return func() {
		var (
			causationID, correlationID ax.MessageID
			store                      persistence.DataStore
			repo                       delayedmessage.Repository
			ctx                        context.Context
			cancel                     func()
		)

		g.BeforeEach(func() {
			store = getStore()
			repo = getRepo()

			var fn func()
			ctx, fn = context.WithTimeout(context.Background(), 15*time.Second)
			cancel = fn // defeat go vet warning about unused cancel func

			causationID = ax.GenerateMessageID()
			correlationID = ax.GenerateMessageID()
		})

		g.AfterEach(func() {
			cancel()
		})

		g.Describe("LoadNextMessage", func() {
			g.Context("when there are messages stored", func() {
				var m1, m2 endpoint.OutboundEnvelope
				var t1, t2, t3 time.Time

				g.BeforeEach(func() {
					t1 = time.Now()
					t2 = t1.Add(1 * time.Second)
					t3 = t2.Add(1 * time.Second)

					m1 = endpoint.OutboundEnvelope{
						Envelope: ax.Envelope{
							MessageID:     ax.GenerateMessageID(),
							CausationID:   causationID,
							CorrelationID: correlationID,
							CreatedAt:     t1,
							SendAt:        t2,
							Message: &testmessages.Command{
								Value: "<foo>",
							},
						},
						Operation:           endpoint.OpSendUnicast,
						DestinationEndpoint: "<dest>",
					}

					m2 = endpoint.OutboundEnvelope{
						Envelope: ax.Envelope{
							MessageID:     ax.GenerateMessageID(),
							CausationID:   causationID,
							CorrelationID: correlationID,
							CreatedAt:     t1,
							SendAt:        t3,
							Message: &testmessages.Event{
								Value: "<bar>",
							},
						},
						Operation: endpoint.OpSendMulticast,
					}

					tx, com, err := store.BeginTx(ctx)
					if err != nil {
						panic(err)
					}
					defer com.Rollback()

					err = repo.SaveMessage(ctx, tx, m1)
					if err != nil {
						panic(err)
					}

					err = repo.SaveMessage(ctx, tx, m2)
					if err != nil {
						panic(err)
					}

					err = com.Commit()
					if err != nil {
						panic(err)
					}
				})

				g.It("returns true", func() {
					_, ok, err := repo.LoadNextMessage(ctx, store)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).To(m.BeTrue())
				})

				g.It("returns the next message", func() {
					env, _, err := repo.LoadNextMessage(ctx, store)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(
						axtest.OutboundEnvelopesEqual(env, m1),
					).To(m.BeTrue())
				})

				g.It("does not return messages that are marked as sent", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					err = repo.MarkAsSent(ctx, tx, m1)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					err = com.Commit()
					m.Expect(err).ShouldNot(m.HaveOccurred())

					env, _, err := repo.LoadNextMessage(ctx, store)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(
						axtest.OutboundEnvelopesEqual(env, m2),
					).To(m.BeTrue())
				})
			})

			g.Context("when there are no messages", func() {
				g.It("returns false", func() {
					_, ok, err := repo.LoadNextMessage(ctx, store)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).To(m.BeFalse())
				})
			})
		})

		g.Describe("SaveMessage", func() {
			g.It("does not return an error if the message already exists", func() {
				tx, com, err := store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				defer com.Rollback()

				env := endpoint.OutboundEnvelope{
					Envelope: ax.Envelope{
						MessageID:     ax.GenerateMessageID(),
						CausationID:   causationID,
						CorrelationID: correlationID,
						CreatedAt:     time.Now(),
						SendAt:        time.Now(),
						Message:       &testmessages.Message{},
					},
					Operation:           endpoint.OpSendUnicast,
					DestinationEndpoint: "<dest>",
				}

				err = repo.SaveMessage(ctx, tx, env)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				err = repo.SaveMessage(ctx, tx, env)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})
		})

		g.Describe("MarkAsSent", func() {
			g.It("does not return an error if the message has already been marked as sent", func() {
				tx, com, err := store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				defer com.Rollback()

				env := endpoint.OutboundEnvelope{
					Envelope: ax.Envelope{
						MessageID:     ax.GenerateMessageID(),
						CausationID:   causationID,
						CorrelationID: correlationID,
						CreatedAt:     time.Now(),
						SendAt:        time.Now(),
						Message:       &testmessages.Message{},
					},
					Operation:           endpoint.OpSendUnicast,
					DestinationEndpoint: "<dest>",
				}

				err = repo.SaveMessage(ctx, tx, env)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				err = com.Commit()
				m.Expect(err).ShouldNot(m.HaveOccurred())

				tx, com, err = store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				defer com.Rollback()

				err = repo.MarkAsSent(ctx, tx, env)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				err = repo.MarkAsSent(ctx, tx, env)
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})
		})
	}
}
