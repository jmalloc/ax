package outboxtest

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/outbox"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/internal/messagetest"
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

			causationID = ax.MessageID{}
			causationID.GenerateUUID()

			correlationID = ax.MessageID{}
			correlationID.GenerateUUID()
		})

		g.AfterEach(func() {
			cancel()
		})

		g.Describe("LoadOutbox", func() {
			g.Context("when the outbox exists", func() {
				var m1, m2 bus.OutboundEnvelope
				var t1, t2 time.Time

				g.BeforeEach(func() {
					t1 = time.Now()
					t2 = time.Now()

					m1 = bus.OutboundEnvelope{
						Envelope: ax.Envelope{
							CausationID:   causationID,
							CorrelationID: correlationID,
							Time:          t1,
							Message: &messagetest.Command{
								Value: "<foo>",
							},
						},
						Operation:           bus.OpSendUnicast,
						DestinationEndpoint: "<dest>",
					}
					m1.MessageID.GenerateUUID()

					m2 = bus.OutboundEnvelope{
						Envelope: ax.Envelope{
							CausationID:   causationID,
							CorrelationID: correlationID,
							Time:          t2,
							Message: &messagetest.Event{
								Value: "<bar>",
							},
						},
						Operation: bus.OpSendMulticast,
					}
					m2.MessageID.GenerateUUID()

					tx, com, err := store.BeginTx(ctx)
					if err != nil {
						panic(err)
					}
					defer com.Rollback()

					err = repo.SaveOutbox(
						ctx,
						tx,
						causationID,
						[]bus.OutboundEnvelope{m1, m2},
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

					for i, env := range envs {
						if env.MessageID == m1.MessageID {
							m.Expect(env.Time).To(m.BeTemporally("==", t1))
						} else if env.MessageID == m2.MessageID {
							m.Expect(env.Time).To(m.BeTemporally("==", t2))
						}

						// zero times for easy comparison
						envs[i].Time = time.Time{}
					}

					// zero times for easy comparison
					m1.Time = time.Time{}
					m2.Time = time.Time{}

					m.Expect(envs).To(m.ConsistOf(m1, m2))
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

					m.Expect(envs).To(m.HaveLen(1))
					m.Expect(envs[0].MessageID).To(m.Equal(m2.MessageID))
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

				env := bus.OutboundEnvelope{
					Envelope: ax.Envelope{
						CausationID:   causationID,
						CorrelationID: correlationID,
						Time:          time.Now(),
						Message:       &messagetest.Message{},
					},
					Operation:           bus.OpSendUnicast,
					DestinationEndpoint: "<dest>",
				}
				env.MessageID.GenerateUUID()

				err = repo.SaveOutbox(
					ctx,
					tx,
					causationID,
					[]bus.OutboundEnvelope{env},
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

				env := bus.OutboundEnvelope{
					Envelope: ax.Envelope{
						CausationID:   causationID,
						CorrelationID: correlationID,
						Time:          time.Now(),
						Message:       &messagetest.Message{},
					},
					Operation:           bus.OpSendUnicast,
					DestinationEndpoint: "<dest>",
				}
				env.MessageID.GenerateUUID()

				err = repo.SaveOutbox(
					ctx,
					tx,
					causationID,
					[]bus.OutboundEnvelope{env},
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
