package sagatests

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/ax/saga/persistence/crud"
	"github.com/jmalloc/ax/src/axtest/testmessages"

	"github.com/jmalloc/ax/src/ax/persistence"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// InsertRev0Saga is a helper function that inserts a test SagaInstance
// with revision 0 into the underlying database of the store.
func InsertRev0Saga(
	ctx context.Context,
	store persistence.DataStore,
	repo crud.Repository,
	pk string,
) saga.Instance {
	var (
		err error
		tx  persistence.Tx
		com persistence.Committer
	)

	i := saga.Instance{
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
	defer com.Rollback()

	if err = repo.SaveSagaInstance(
		ctx,
		tx,
		pk,
		i,
	); err != nil {
		panic(err)
	}

	if err = com.Commit(); err != nil {
		panic(err)
	}
	return i
}

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
					r0  saga.Instance
					pk  string
					tx  persistence.Tx
					com persistence.Committer
				)
				g.BeforeEach(func() {
					var err error
					pk = "<test>"
					r0 = InsertRev0Saga(ctx, store, repo, pk)

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
						r0.InstanceID,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).To(m.BeTrue())
				})

				g.It("returns the saga instance from the repository", func() {
					i, _, err := repo.LoadSagaInstance(
						ctx,
						tx,
						pk,
						r0.InstanceID,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(i.InstanceID).Should(m.Equal(r0.InstanceID))
					m.Expect(i.Revision).Should(m.Equal(r0.Revision))
					m.Expect(proto.Equal(i.Data, r0.Data)).Should(m.BeTrue())
				})
			})

			g.Context("when the instance does not exists", func() {
				var (
					pk  string
					tx  persistence.Tx
					com persistence.Committer
				)
				g.BeforeEach(func() {
					var err error
					pk = "<test>"
					tx, com, err = store.BeginTx(ctx)
					if err != nil {
						panic(err)
					}
				})
				g.AfterEach(func() {
					com.Rollback()
				})

				g.It("returns false", func() {
					id := saga.GenerateInstanceID()
					_, ok, err := repo.LoadSagaInstance(
						ctx,
						tx,
						pk,
						id,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).To(m.BeFalse())
				})
			})

			g.Context("instance is found, but belongs to a different saga", func() {
				var (
					r0  saga.Instance
					pk  string
					tx  persistence.Tx
					com persistence.Committer
				)
				g.BeforeEach(func() {
					var err error
					pk = "<test>"
					r0 = InsertRev0Saga(ctx, store, repo, pk)

					tx, com, err = store.BeginTx(ctx)
					if err != nil {
						panic(err)
					}
				})
				g.AfterEach(func() {
					com.Rollback()
				})

				g.It("returns an error", func() {
					_, _, err := repo.LoadSagaInstance(
						ctx,
						tx,
						"<unknown>",
						r0.InstanceID,
					)
					m.Expect(err).Should(m.HaveOccurred())
				})
			})
		})

		g.Describe("SaveSagaInstance", func() {
			g.Context("when the instance doesn't exist (insert)", func() {
				var (
					pk  string
					tx  persistence.Tx
					com persistence.Committer
				)
				g.BeforeEach(func() {
					var err error
					pk = "<test>"
					tx, com, err = store.BeginTx(ctx)
					if err != nil {
						panic(err)
					}
				})

				g.AfterEach(func() {
					com.Rollback()
				})

				g.It("returns no errors", func() {
					r0 := saga.Instance{
						InstanceID: saga.GenerateInstanceID(),
						Revision:   saga.Revision(0),
						Data: &testmessages.Data{
							Value: "<foo>",
						},
					}
					err := repo.SaveSagaInstance(
						ctx,
						tx,
						pk,
						r0,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					i2, ok, err := repo.LoadSagaInstance(
						ctx,
						tx,
						pk,
						r0.InstanceID,
					)

					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).To(m.BeTrue())
					m.Expect(r0.InstanceID).Should(m.Equal(i2.InstanceID))
					m.Expect(r0.Revision).Should(m.Equal(i2.Revision))
					m.Expect(proto.Equal(r0.Data, i2.Data)).Should(m.BeTrue())
				})
			})

			g.Context("when the instance exists (update)", func() {
				var (
					r0  saga.Instance
					pk  string
					tx  persistence.Tx
					com persistence.Committer
				)
				g.BeforeEach(func() {
					var err error
					pk = "<test>"
					r0 = InsertRev0Saga(ctx, store, repo, pk)

					tx, com, err = store.BeginTx(ctx)
					if err != nil {
						panic(err)
					}
				})

				g.AfterEach(func() {
					com.Rollback()
				})

				g.It("returns no errors", func() {
					err := repo.SaveSagaInstance(
						ctx,
						tx,
						pk,
						r0,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					r1, ok, err := repo.LoadSagaInstance(
						ctx,
						tx,
						pk,
						r0.InstanceID,
					)

					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).To(m.BeTrue())
					m.Expect(r1.InstanceID).Should(m.Equal(r0.InstanceID))
					m.Expect(r1.Revision).Should(m.BeNumerically("==", r0.Revision+1))
					m.Expect(proto.Equal(r1.Data, r0.Data)).Should(m.BeTrue())
				})

				g.Context("when the revision is not current for the existing saga instance", func() {
					g.It("returns an error", func() {
						r1 := saga.Instance{
							InstanceID: r0.InstanceID,
							Revision:   saga.Revision(1),
							Data: &testmessages.Data{
								Value: "<foo>",
							},
						}
						err := repo.SaveSagaInstance(
							ctx,
							tx,
							pk,
							r1,
						)
						m.Expect(err).Should(m.HaveOccurred())
					})
				})
				g.Context("when the instance exists, but belongs to a different saga", func() {
					g.It("returns an error", func() {
						err := repo.SaveSagaInstance(
							ctx,
							tx,
							"<unknown>",
							r0,
						)
						m.Expect(err).Should(m.HaveOccurred())
					})
				})
			})
		})

		g.Describe("DeleteSagaInstance", func() {
			var (
				tx  persistence.Tx
				com persistence.Committer
				r0  saga.Instance
				pk  string
			)
			g.BeforeEach(func() {
				var err error
				pk = "<test>"
				r0 = InsertRev0Saga(ctx, store, repo, pk)

				tx, com, err = store.BeginTx(ctx)
				if err != nil {
					panic(err)
				}
			})
			g.AfterEach(func() {
				com.Rollback()
			})

			g.It("returns no errors", func() {
				var (
					err error
					ok  bool
				)
				err = repo.DeleteSagaInstance(
					ctx,
					tx,
					pk,
					r0,
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				_, ok, err = repo.LoadSagaInstance(
					ctx,
					tx,
					pk,
					r0.InstanceID,
				)
				m.Expect(ok).Should(m.BeFalse())
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})

			g.Context("when the revision is not current for the existing saga instance", func() {

			})

			g.Context("when the instance exists, but belongs to a different saga", func() {

			})
		})
	}
}
