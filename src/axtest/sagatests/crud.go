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
					m.Expect(i.Revision).Should(m.Equal(expected.Revision))
					m.Expect(proto.Equal(i.Data, expected.Data)).Should(m.BeTrue())
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

				g.It("returns an error", func() {
					_, _, err := repo.LoadSagaInstance(
						ctx,
						tx,
						"<unknown>",
						expected.InstanceID,
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

				g.It("returns nil", func() {
					i := saga.Instance{
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
						i,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
				})
			})

			g.Context("when the instance exists (update)", func() {

			})

			g.Context("when the revision is not current for the existing saga instance", func() {

			})

			g.Context("when the instance exists, but belongs to a different saga", func() {

			})
		})
	}
}
