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
	if dsn == "" {
		return
	}

	var db *sql.DB

	BeforeEach(func() {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}

		_, err = db.Exec("DROP TABLE IF EXISTS outbox")
		if err != nil {
			panic(err)
		}

		_, err = db.Exec("DROP TABLE IF EXISTS outbox_message")
		if err != nil {
			panic(err)
		}

		for _, q := range OutboxSchema {
			_, err = db.Exec(q)
			if err != nil {
				panic(err)
			}
		}
	})

	AfterEach(func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	})

	outboxtest.DescribeRepository(
		func() persistence.DataStore {
			return &DataStore{DB: db}
		},
		func() outbox.Repository {
			return &OutboxRepository{}
		},
	)
})
