package sagatests

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/ax/saga/persistence/crud"
	"github.com/jmalloc/ax/src/axtest/testmessages"

	"github.com/jmalloc/ax/src/ax/persistence"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// CRUDRepositorySuite returns a test suite for implementations of crud.Repository.
func CRUDRepositorySuite(
	getStore func() persistence.DataStore,
	getRepo func() crud.Repository,
) func() {
	return func() {
		var (
			store  persistence.DataStore
			repo   crud.Repository
			ctx    context.Context
			cancel func()
		)

		g.BeforeEach(func() {
			store = getStore()
			repo = getRepo()

			var fn func()
			ctx, fn = context.WithTimeout(context.Background(), 15*time.Second)
			cancel = fn // defeat go vet warning about unused cancel func
		})

		g.AfterEach(func() {
			cancel()
		})

		g.Describe("LoadSagaInstance", func() {
			g.Context("when the instance exists", func() {
				var (
					expected saga.Instance
					pk       string
					tx       persistence.Tx
					com      persistence.Committer
				)
				g.BeforeEach(func() {
					var err error
					pk = "<test>"
					expected = saga.Instance{
						InstanceID: saga.GenerateInstanceID(),
						Revision:   saga.Revision(0),
						Data: &testmessages.Data{
							Value: "<foo>",
						},
					}

					tx, com, err = store.BeginTx(ctx)
					if err != nil {
						panic(err)
					}

					err = repo.SaveSagaInstance(
						ctx,
						tx,
						pk,
						expected,
					)
					if err != nil {
						panic(err)
					}

					err = com.Commit()
					if err != nil {
						panic(err)
					}

					tx, com, err = store.BeginTx(ctx)
					if err != nil {
						panic(err)
					}
				})
				g.AfterEach(func() {
					com.Rollback()
				})

				g.It("returns true", func() {
					_, ok, err := repo.LoadSagaInstance(
						ctx,
						tx,
						pk,
						expected.InstanceID,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).To(m.BeTrue())
				})

				g.It("returns the saga instance from the repository", func() {
					i, _, err := repo.LoadSagaInstance(
						ctx,
						tx,
						pk,
						expected.InstanceID,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(i.InstanceID).Should(m.Equal(expected.InstanceID))
				})

				// g.It("does not return messages that are marked as sent", func() {
				// 	tx, com, err := store.BeginTx(ctx)
				// 	m.Expect(err).ShouldNot(m.HaveOccurred())
				// 	defer com.Rollback()

				// 	err = repo.MarkAsSent(
				// 		ctx,
				// 		tx,
				// 		m1,
				// 	)
				// 	m.Expect(err).ShouldNot(m.HaveOccurred())

				// 	err = com.Commit()
				// 	m.Expect(err).ShouldNot(m.HaveOccurred())

				// 	envs, _, err := repo.LoadOutbox(ctx, store, causationID)
				// 	m.Expect(err).ShouldNot(m.HaveOccurred())
				// 	m.Expect(
				// 		axtest.ConsistsOfOutboundEnvelopes(envs, m2),
				// 	).To(m.BeTrue())
				// })
			})

			// g.Context("when the outbox exists but contains no messages", func() {
			// 	g.BeforeEach(func() {
			// 		tx, com, err := store.BeginTx(ctx)
			// 		if err != nil {
			// 			panic(err)
			// 		}
			// 		defer com.Rollback()

			// 		err = repo.SaveOutbox(
			// 			ctx,
			// 			tx,
			// 			causationID,
			// 			nil,
			// 		)
			// 		if err != nil {
			// 			panic(err)
			// 		}

			// 		err = com.Commit()
			// 		if err != nil {
			// 			panic(err)
			// 		}
			// 	})

			// 	g.It("returns true", func() {
			// 		_, ok, err := repo.LoadOutbox(ctx, store, causationID)
			// 		m.Expect(err).ShouldNot(m.HaveOccurred())
			// 		m.Expect(ok).To(m.BeTrue())
			// 	})
			// })

			// g.Context("when the outbox does not exist", func() {
			// 	g.It("returns false", func() {
			// 		_, ok, err := repo.LoadOutbox(ctx, store, causationID)
			// 		m.Expect(err).ShouldNot(m.HaveOccurred())
			// 		m.Expect(ok).To(m.BeFalse())
			// 	})
			// })
		})

		// g.Describe("SaveOutbox", func() {
		// 	g.It("returns an error if the outbox already exists", func() {
		// 		tx, com, err := store.BeginTx(ctx)
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())
		// 		defer com.Rollback()

		// 		err = repo.SaveOutbox(
		// 			ctx,
		// 			tx,
		// 			causationID,
		// 			nil,
		// 		)
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())

		// 		err = repo.SaveOutbox(
		// 			ctx,
		// 			tx,
		// 			causationID,
		// 			nil,
		// 		)
		// 		m.Expect(err).Should(m.HaveOccurred())
		// 	})

		// 	g.It("does not add new messages to an existing outbox", func() {
		// 		tx, com, err := store.BeginTx(ctx)
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())
		// 		defer com.Rollback()

		// 		err = repo.SaveOutbox(
		// 			ctx,
		// 			tx,
		// 			causationID,
		// 			nil,
		// 		)
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())

		// 		env := endpoint.OutboundEnvelope{
		// 			Envelope: ax.Envelope{
		// 				MessageID:     ax.GenerateMessageID(),
		// 				CausationID:   causationID,
		// 				CorrelationID: correlationID,
		// 				CreatedAt:     time.Now(),
		// 				SendAt:        time.Now(),
		// 				Message:       &testmessages.Message{},
		// 			},
		// 			Operation:           endpoint.OpSendUnicast,
		// 			DestinationEndpoint: "<dest>",
		// 		}

		// 		err = repo.SaveOutbox(
		// 			ctx,
		// 			tx,
		// 			causationID,
		// 			[]endpoint.OutboundEnvelope{env},
		// 		)
		// 		m.Expect(err).Should(m.HaveOccurred())

		// 		err = com.Commit()
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())

		// 		envs, ok, err := repo.LoadOutbox(ctx, store, causationID)
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())
		// 		m.Expect(envs).To(m.BeEmpty())
		// 		m.Expect(ok).To(m.BeTrue())
		// 	})
		// })

		// g.Describe("MarkAsSent", func() {
		// 	g.It("does not return an error if the message has already been marked as sent", func() {
		// 		tx, com, err := store.BeginTx(ctx)
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())
		// 		defer com.Rollback()

		// 		env := endpoint.OutboundEnvelope{
		// 			Envelope: ax.Envelope{
		// 				MessageID:     ax.GenerateMessageID(),
		// 				CausationID:   causationID,
		// 				CorrelationID: correlationID,
		// 				CreatedAt:     time.Now(),
		// 				SendAt:        time.Now(),
		// 				Message:       &testmessages.Message{},
		// 			},
		// 			Operation:           endpoint.OpSendUnicast,
		// 			DestinationEndpoint: "<dest>",
		// 		}

		// 		err = repo.SaveOutbox(
		// 			ctx,
		// 			tx,
		// 			causationID,
		// 			[]endpoint.OutboundEnvelope{env},
		// 		)
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())

		// 		err = com.Commit()
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())

		// 		tx, com, err = store.BeginTx(ctx)
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())
		// 		defer com.Rollback()

		// 		err = repo.MarkAsSent(ctx, tx, env)
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())

		// 		err = repo.MarkAsSent(ctx, tx, env)
		// 		m.Expect(err).ShouldNot(m.HaveOccurred())
		// 	})
		// })
	}
}
