package outbox_test

import (
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/ax/outbox"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axbolt"
	. "github.com/jmalloc/ax/src/axbolt/outbox"
	"github.com/jmalloc/ax/src/axtest/outboxtests"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Repository", func() {
	fname := os.Getenv("AX_BOLT_DB")
	var db *bolt.DB

	BeforeEach(func() {
		var err error
		db, err = bolt.Open(fname, 0600, nil)
		if err != nil {
			panic(err)
		}
	})

	AfterEach(func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
		_ = os.Remove(fname)
	})

	fn := Describe
	if fname == "" {
		fn = XDescribe
	}

	fn(
		"Repository",
		outboxtests.RepositorySuite(
			func() persistence.DataStore {
				return axbolt.NewDataStore(db)
			},
			func() outbox.Repository {
				return Repository{}
			},
		),
	)
})
