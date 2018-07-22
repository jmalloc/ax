package sagatests

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/ax/saga/mapping/keyset"

	"github.com/jmalloc/ax/src/ax/persistence"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// InsertMappingKeys inserts saga mapping keys into the repository
func InsertMappingKeys(
	ctx context.Context,
	store persistence.DataStore,
	repo keyset.Repository,
	pk string,
	id saga.InstanceID,
	ks []string,
) {
	var (
		err error
		tx  persistence.Tx
		com persistence.Committer
	)
	tx, com, err = store.BeginTx(ctx)
	if err != nil {
		panic(err)
	}
	defer com.Rollback()

	if err = repo.SaveKeys(
		ctx,
		tx,
		pk,
		ks,
		id,
	); err != nil {
		panic(err)
	}

	if err = com.Commit(); err != nil {
		panic(err)
	}
}

// KeySetRepositorySuite returns a test suite for implementations of keyset.Repository.
func KeySetRepositorySuite(
	getStore func() persistence.DataStore,
	getRepo func() keyset.Repository,
) func() {
	return func() {
		const (
			pk = "<test>"
		)
		var (
			store  persistence.DataStore
			repo   keyset.Repository
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

		g.Describe("FindByKey", func() {
			g.Context("saga instance's id with mapping key exists", func() {
				var (
					mk []string
					id saga.InstanceID
				)
				g.BeforeEach(func() {
					mk = []string{"<test1>", "<test2>", "<test3>"}
					id = saga.GenerateInstanceID()
					InsertMappingKeys(
						ctx,
						store,
						repo,
						pk,
						id,
						mk,
					)
				})
				g.It("returns true", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					for _, k := range mk {
						_, ok, err := repo.FindByKey(
							ctx,
							tx,
							pk,
							k,
						)
						m.Expect(err).ShouldNot(m.HaveOccurred())
						m.Expect(ok).Should(m.BeTrue())
					}
				})

				g.It("returns saga instance id", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					for _, k := range mk {
						iid, _, err := repo.FindByKey(
							ctx,
							tx,
							pk,
							k,
						)
						m.Expect(err).ShouldNot(m.HaveOccurred())
						m.Expect(iid.Get()).Should(m.Equal(id.Get()))
					}
				})
			})
			g.Context("saga instance's id with mapping key does not exist", func() {
				g.It("returns false", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					_, ok, err := repo.FindByKey(
						ctx,
						tx,
						pk,
						"<unknown>",
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).Should(m.BeFalse())
				})
			})
		})

		g.Describe("SaveKeys", func() {
			g.It("returns nil", func() {
				tx, com, err := store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				defer com.Rollback()

				err = repo.SaveKeys(
					ctx,
					tx,
					pk,
					[]string{"<test1>", "<test2>", "<test3>"},
					saga.GenerateInstanceID(),
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				err = com.Commit()
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})
		})

		g.Describe("DeleteKeys", func() {
			var (
				mk []string
				id saga.InstanceID
			)
			g.BeforeEach(func() {
				mk = []string{"<test1>", "<test2>", "<test3>"}
				id = saga.GenerateInstanceID()
				InsertMappingKeys(
					ctx,
					store,
					repo,
					pk,
					id,
					mk,
				)
			})
			g.It("returns nil", func() {
				tx, com, err := store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				defer com.Rollback()

				err = repo.DeleteKeys(
					ctx,
					tx,
					pk,
					id,
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				for _, k := range mk {
					_, ok, err := repo.FindByKey(
						ctx,
						tx,
						pk,
						k,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).Should(m.BeFalse())
				}

				err = com.Commit()
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})
		})
	}
}
