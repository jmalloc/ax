package messagestore_test

import (
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/ax/messagestore"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axbolt"
	axboltmessagestore "github.com/jmalloc/ax/src/axbolt/messagestore"
	"github.com/jmalloc/ax/src/axtest/messagestoretests"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("MessageStore", func() {
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
		"MessageStore",
		messagestoretests.MessageStoreSuite(
			func() persistence.DataStore {
				return axbolt.NewDataStore(db)
			},
			func() messagestore.GloballyOrderedStore {
				return axboltmessagestore.Store{}
			},
		),
	)
})
