package sagatests

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax/saga/mapping/keyset"

	"github.com/jmalloc/ax/src/ax/persistence"
	g "github.com/onsi/ginkgo"
)

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
			_ = ctx
			_ = store
			_ = repo
		})

		g.Describe("SaveKeys", func() {
		})

		g.Describe("DeleteKeys", func() {
		})
	}
}
