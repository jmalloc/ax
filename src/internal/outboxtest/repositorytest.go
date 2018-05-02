package outboxtest

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/outbox"
	"github.com/jmalloc/ax/src/ax/persistence"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// DescribeRepository describes an integration test suite for implementations of
// outbox.Repository.
func DescribeRepository(
	getStore func() persistence.DataStore,
	getRepo func() outbox.Repository,
) {
	var (
		store  persistence.DataStore
		repo   outbox.Repository
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

	g.Describe("LoadOutbox", func() {
		g.It("returns false if no outbox has been saved", func() {
			var id ax.MessageID
			id.GenerateUUID()

			_, ok, err := repo.LoadOutbox(ctx, store, id)
			m.Expect(err).ShouldNot(m.HaveOccurred())
			m.Expect(ok).To(m.BeFalse())
		})
	})
}
