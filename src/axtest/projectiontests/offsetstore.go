package projectiontests

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/projection"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// insertOffset inserts an offset into a projection offset store. The value of c
// will be increased by the internal call of repo.IncrementOffset
func insertOffset(
	ctx context.Context,
	store persistence.DataStore,
	repo projection.OffsetStore,
	pk string,
	c uint64,
) uint64 {
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

	if repo.IncrementOffset(ctx, tx, pk, c); err != nil {
		panic(err)
	}

	if err = com.Commit(); err != nil {
		panic(err)
	}
	c++
	return c
}

// OffsetStoreSuite returns a test suite for implementations of
// projection.OffsetStore.
func OffsetStoreSuite(
	getStore func() persistence.DataStore,
	getRepo func() projection.OffsetStore,
) func() {
	return func() {
		var (
			store  persistence.DataStore
			repo   projection.OffsetStore
			ctx    context.Context
			cancel func()
			pk     string
		)

		g.BeforeEach(func() {
			store = getStore()
			repo = getRepo()

			var fn func()
			ctx, fn = context.WithTimeout(context.Background(), 15*time.Second)
			cancel = fn // defeat go vet warning about unused cancel func
			pk = "<test>"
		})

		g.AfterEach(func() {
			cancel()
		})

		g.Describe("LoadOffset", func() {
			g.Context("when previous offset exists", func() {
				var prevOffset uint64
				g.BeforeEach(func() {
					prevOffset = insertOffset(ctx, store, repo, pk, 0)
				})
				g.It("returns no errors", func() {
					_, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					_, err = repo.LoadOffset(ctx, store, pk)
					m.Expect(err).ShouldNot(m.HaveOccurred())
				})
				g.It("returns incremented offset", func() {
					_, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					offset, _ := repo.LoadOffset(ctx, store, pk)
					m.Expect(offset).Should(m.BeNumerically("==", prevOffset))
				})
			})
			g.Context("when previous offset does not exist", func() {
				g.It("returns no errors", func() {
					_, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					_, err = repo.LoadOffset(ctx, store, pk)
					m.Expect(err).ShouldNot(m.HaveOccurred())
				})
				g.It("returns zero offset", func() {
					_, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					offset, _ := repo.LoadOffset(ctx, store, pk)
					m.Expect(offset).Should(m.BeNumerically("==", 0))
				})
			})
		})

		g.Describe("IncrementOffset", func() {
			g.Context("when previous offset exists", func() {
				var prevOffset uint64
				g.BeforeEach(func() {
					prevOffset = insertOffset(ctx, store, repo, pk, 0)
				})
				g.It("returns no errors", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					err = repo.IncrementOffset(ctx, tx, pk, prevOffset)
					m.Expect(err).ShouldNot(m.HaveOccurred())
				})
				g.Context("when the value passed does not equal to offset stored", func() {
					g.It("returns an error", func() {
						tx, com, err := store.BeginTx(ctx)
						m.Expect(err).ShouldNot(m.HaveOccurred())
						defer com.Rollback()

						err = repo.IncrementOffset(ctx, tx, pk, uint64(3))
						m.Expect(err).Should(m.HaveOccurred())
					})
				})
			})
			g.Context("when previous offset does not exist", func() {
				g.It("returns no errors", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					err = repo.IncrementOffset(ctx, tx, pk, 0)
					m.Expect(err).ShouldNot(m.HaveOccurred())
				})
			})
		})
	}
}
