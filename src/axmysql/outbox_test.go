package axmysql_test

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmalloc/ax/src/ax/outbox"
	"github.com/jmalloc/ax/src/ax/persistence"
	. "github.com/jmalloc/ax/src/axmysql"
	"github.com/jmalloc/ax/src/internal/outboxtest"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("OutboxRepository", func() {
	dsn := os.Getenv("AX_MYSQL_DSN")
	var db *sql.DB

	BeforeEach(func() {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}

		if err := createSchema(db, "outbox.sql"); err != nil {
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
		"OutboxRepository",
		outboxtest.RepositorySuite(
			func() persistence.DataStore {
				return NewDataStore(db)
			},
			func() outbox.Repository {
				return &OutboxRepository{}
			},
		),
	)
})
