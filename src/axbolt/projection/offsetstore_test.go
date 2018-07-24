package projection_test

import (
	"os"

	"github.com/jmalloc/ax/src/ax/projection"

	axboltprojection "github.com/jmalloc/ax/src/axbolt/projection"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axbolt"
	"github.com/jmalloc/ax/src/axtest/projectiontests"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Projection Offset Store", func() {
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
		"Projection Offset Store",
		projectiontests.OffsetStoreSuite(
			func() persistence.DataStore {
				return axbolt.NewDataStore(db)
			},
			func() projection.OffsetStore {
				return axboltprojection.OffsetStore{}
			},
		),
	)
})
