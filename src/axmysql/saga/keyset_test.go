package saga_test

import (
	"database/sql"
	"os"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga/mapping/keyset"
	"github.com/jmalloc/ax/src/axmysql"
	"github.com/jmalloc/ax/src/axmysql/internal/schema"
	"github.com/jmalloc/ax/src/axmysql/saga"
	"github.com/jmalloc/ax/src/axtest/sagatests"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Keyset Repository", func() {
	dsn := os.Getenv("AX_MYSQL_DSN")
	var db *sql.DB

	BeforeEach(func() {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}

		if err := schema.Create(db, "keyset.sql"); err != nil {
			panic(err)
		}
	})

	AfterEach(func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	})

	fn := Describe
	if dsn == "" {
		fn = XDescribe
	}

	fn(
		"Keyset Repository",
		sagatests.KeySetRepositorySuite(
			func() persistence.DataStore {
				return axmysql.NewDataStore(db)
			},
			func() keyset.Repository {
				return saga.KeySetRepository{}
			},
		),
	)
})
