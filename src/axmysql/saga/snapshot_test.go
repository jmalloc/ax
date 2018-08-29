package saga_test

import (
	"database/sql"
	"os"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga/persistence/eventsourcing"
	"github.com/jmalloc/ax/src/axmysql"
	"github.com/jmalloc/ax/src/axmysql/internal/schema"
	"github.com/jmalloc/ax/src/axmysql/saga"
	"github.com/jmalloc/ax/src/axtest/sagatests"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("SnapshotRepository", func() {
	dsn := os.Getenv("AX_MYSQL_DSN")
	var db *sql.DB

	BeforeEach(func() {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}

		if err := schema.Create(db, "snapshot.sql"); err != nil {
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
		"SnapshotRepository",
		sagatests.SnapshotRepositorySuite(
			func() persistence.DataStore {
				return axmysql.NewDataStore(db)
			},
			func() eventsourcing.SnapshotRepository {
				return saga.SnapshotRepository{}
			},
		),
	)
})
