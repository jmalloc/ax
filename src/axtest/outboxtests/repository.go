package outboxtests

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/axtest"
	"github.com/jmalloc/ax/src/axtest/testmessages"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/outbox"
	"github.com/jmalloc/ax/src/ax/persistence"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// RepositorySuite returns a test suite for implementations of outbox.Repository.
func RepositorySuite(
	getStore func() persistence.DataStore,
	getRepo func() outbox.Repository,
) func() {
	return func() {
		var (
			causationID, correlationID ax.MessageID
			store                      persistence.DataStore
			repo                       outbox.Repository
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

		g.Describe("LoadOutbox", func() {
			g.Context("when the outbox exists", func() {
				var m1, m2 endpoint.OutboundEnvelope
				var t1, t2, t3, t4 time.Time

				g.BeforeEach(func() {
					t1 = time.Now()
					t2 = time.Now()
					t3 = time.Now()
					t4 = time.Now()

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
							CreatedAt:     t3,
							SendAt:        t4,
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

					err = repo.SaveOutbox(
						ctx,
						tx,
						causationID,
						[]endpoint.OutboundEnvelope{m1, m2},
					)
					if err != nil {
						panic(err)
					}

					err = com.Commit()
					if err != nil {
						panic(err)
					}
				})

				g.It("returns true", func() {
					_, ok, err := repo.LoadOutbox(ctx, store, causationID)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).To(m.BeTrue())
				})

				g.It("returns the messages in the outbox", func() {
					envs, _, err := repo.LoadOutbox(ctx, store, causationID)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(
						axtest.ConsistsOfOutboundEnvelopes(envs, m1, m2),
					).To(m.BeTrue())
				})

				g.It("does not return messages that are marked as sent", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					err = repo.MarkAsSent(
						ctx,
						tx,
						m1,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					err = com.Commit()
					m.Expect(err).ShouldNot(m.HaveOccurred())

					envs, _, err := repo.LoadOutbox(ctx, store, causationID)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(
						axtest.ConsistsOfOutboundEnvelopes(envs, m2),
					).To(m.BeTrue())
				})
			})

			g.Context("when the outbox exists but contains no messages", func() {
				g.BeforeEach(func() {
					tx, com, err := store.BeginTx(ctx)
					if err != nil {
						panic(err)
					}
					defer com.Rollback()

					err = repo.SaveOutbox(
						ctx,
						tx,
						causationID,
						nil,
					)
					if err != nil {
						panic(err)
					}

					err = com.Commit()
					if err != nil {
						panic(err)
					}
				})

				g.It("returns true", func() {
					_, ok, err := repo.LoadOutbox(ctx, store, causationID)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).To(m.BeTrue())
				})
			})

			g.Context("when the outbox does not exist", func() {
				g.It("returns false", func() {
					_, ok, err := repo.LoadOutbox(ctx, store, causationID)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).To(m.BeFalse())
				})
			})
		})

		g.Describe("SaveOutbox", func() {
			g.It("returns an error if the outbox already exists", func() {
				tx, com, err := store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				defer com.Rollback()

				err = repo.SaveOutbox(
					ctx,
					tx,
					causationID,
					nil,
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				err = repo.SaveOutbox(
					ctx,
					tx,
					causationID,
					nil,
				)
				m.Expect(err).Should(m.HaveOccurred())
			})

			g.It("does not add new messages to an existing outbox", func() {
				tx, com, err := store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				defer com.Rollback()

				err = repo.SaveOutbox(
					ctx,
					tx,
					causationID,
					nil,
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())

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

				err = repo.SaveOutbox(
					ctx,
					tx,
					causationID,
					[]endpoint.OutboundEnvelope{env},
				)
				m.Expect(err).Should(m.HaveOccurred())

				err = com.Commit()
				m.Expect(err).ShouldNot(m.HaveOccurred())

				envs, ok, err := repo.LoadOutbox(ctx, store, causationID)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(envs).To(m.BeEmpty())
				m.Expect(ok).To(m.BeTrue())
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

				err = repo.SaveOutbox(
					ctx,
					tx,
					causationID,
					[]endpoint.OutboundEnvelope{env},
				)
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
