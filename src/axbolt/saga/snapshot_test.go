package saga_test

import (
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga/persistence/eventsourcing"
	"github.com/jmalloc/ax/src/axbolt"
	"github.com/jmalloc/ax/src/axbolt/saga"
	"github.com/jmalloc/ax/src/axtest/sagatests"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Snapshot Repository", func() {
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
		"Snapshost Repository",
		sagatests.SnapshotRepositorySuite(
			func() persistence.DataStore {
				return axbolt.NewDataStore(db)
			},
			func() eventsourcing.SnapshotRepository {
				return saga.SnapshotRepository{}
			},
		),
	)
})
