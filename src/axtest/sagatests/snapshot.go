package sagatests

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/ax/saga/persistence/eventsourcing"
	"github.com/jmalloc/ax/src/axtest/testmessages"

	"github.com/jmalloc/ax/src/ax/persistence"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// insertSagaSnapshot insert a saga snapshot into the snapshot repository It
// panics if any error might result in the process of the snapshot insertion.
func insertSagaSnapshot(
	ctx context.Context,
	store persistence.DataStore,
	i saga.Instance,
	repo eventsourcing.SnapshotRepository,
	pk string,
) saga.Instance {
	tx, com, err := store.BeginTx(ctx)
	if err != nil {
		panic(err)
	}

	if err = repo.SaveSagaSnapshot(
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

// SnapshotRepositorySuite returns a test suite for implementations of
// snapshot.Repository.
func SnapshotRepositorySuite(
	getStore func() persistence.DataStore,
	getRepo func() eventsourcing.SnapshotRepository,
) func() {
	return func() {
		const (
			pk = "<test>"
		)
		var (
			store  persistence.DataStore
			repo   eventsourcing.SnapshotRepository
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

		g.Describe("LoadSagaSnapshot", func() {
			g.Context("when the latest snapshot exists", func() {
				var (
					i      saga.Instance
					latest saga.Revision
				)
				g.BeforeEach(func() {
					i = saga.Instance{
						InstanceID: saga.GenerateInstanceID(),
						Data: &testmessages.Data{
							Value: "<foo>",
						},
					}
					latest = saga.Revision(3)
					for r := latest; r > saga.Revision(0); r -= saga.Revision(1) {
						i.Revision = r
						i = insertSagaSnapshot(
							ctx,
							store,
							i,
							repo,
							pk,
						)
					}
				})

				g.It("returns true", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					_, ok, err := repo.LoadSagaSnapshot(
						ctx,
						tx,
						pk,
						i.InstanceID,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).Should(m.BeTrue())
				})

				g.It("returns the latest snapshot of the saga instance from repository", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					l, _, err := repo.LoadSagaSnapshot(
						ctx,
						tx,
						pk,
						i.InstanceID,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(l.InstanceID).Should(m.Equal(i.InstanceID))
					m.Expect(l.Revision).Should(m.BeNumerically("==", latest))
					m.Expect(proto.Equal(l.Data, i.Data)).Should(m.BeTrue())
				})

				g.Context("but it belongs to a different saga", func() {
					g.It("returns error", func() {
						tx, com, err := store.BeginTx(ctx)
						m.Expect(err).ShouldNot(m.HaveOccurred())
						defer com.Rollback()

						_, _, err = repo.LoadSagaSnapshot(
							ctx,
							tx,
							"<unknown>",
							i.InstanceID,
						)
						m.Expect(err).Should(m.HaveOccurred())
					})
				})
			})
		})

		g.Describe("SaveSagaSnapshot", func() {
			g.It("returns no errors", func() {
				var (
					tx  persistence.Tx
					com persistence.Committer
					i   saga.Instance
					err error
				)
				tx, com, err = store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				i = saga.Instance{
					InstanceID: saga.GenerateInstanceID(),
					Data: &testmessages.Data{
						Value: "<foo>",
					},
				}
				err = repo.SaveSagaSnapshot(
					ctx,
					tx,
					pk,
					i,
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				err = com.Commit()
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})
		})

		g.Describe("DeleteSagaSnapshots", func() {
			var (
				i saga.Instance
			)
			g.BeforeEach(func() {
				i = saga.Instance{
					InstanceID: saga.GenerateInstanceID(),
					Data: &testmessages.Data{
						Value: "<foo>",
					},
				}
				i = insertSagaSnapshot(
					ctx,
					store,
					i,
					repo,
					pk,
				)
			})

			g.It("returns no errors", func() {
				tx, com, err := store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				err = repo.DeleteSagaSnapshots(
					ctx,
					tx,
					pk,
					i.InstanceID,
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				_, ok, err := repo.LoadSagaSnapshot(
					ctx,
					tx,
					pk,
					i.InstanceID,
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				m.Expect(ok).Should(m.BeFalse())

				err = com.Commit()
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})
		})
	}
}
