package projectiontests

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/projection"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// InsertOffset inserts an offset into a projection offset store.
// The value of c will be increased by the internal call of repo.IncrementOffset
func InsertOffset(
	ctx context.Context,
	store persistence.DataStore,
	repo projection.OffsetStore,
	pk string,
	c uint64,
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

	repo.IncrementOffset(ctx, tx, pk, c)
	if err != nil {
		panic(err)
	}
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
			g.Context("when offset exists", func() {
				var initialOffset uint64
				g.BeforeEach(func() {
					InsertOffset(ctx, store, repo, pk, initialOffset)
				})
				g.It("returns no errors", func() {
					_, err := repo.LoadOffset(ctx, store, pk)
					m.Expect(err).ShouldNot(m.HaveOccurred())
				})
				// g.It("returns the offset", func() {
				// 	offset, _ := repo.LoadOffset(ctx, store, pk)
				// 	m.Expect(offset).Should(m.BeNumerically("==", initialOffset+1))
				// })
			})
		})

		g.Describe("IncrementOffset", func() {

		})
	}
}
