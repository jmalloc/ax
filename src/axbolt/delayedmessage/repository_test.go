package delayedmessage_test

import (
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/ax/delayedmessage"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axbolt"
	boltdelayedmessage "github.com/jmalloc/ax/src/axbolt/delayedmessage"
	"github.com/jmalloc/ax/src/axtest/delayedmessagetests"
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
		delayedmessagetests.RepositorySuite(
			func() persistence.DataStore {
				return axbolt.NewDataStore(db)
			},
			func() delayedmessage.Repository {
				return boltdelayedmessage.Repository{}
			},
		),
	)
})
