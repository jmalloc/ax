package outbox_test

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmalloc/ax/axmysql"
	"github.com/jmalloc/ax/axmysql/internal/schema"
	. "github.com/jmalloc/ax/axmysql/outbox"
	"github.com/jmalloc/ax/axtest/outboxtests"
	"github.com/jmalloc/ax/outbox"
	"github.com/jmalloc/ax/persistence"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Repository", func() {
	dsn := os.Getenv("AX_MYSQL_DSN")
	var db *sql.DB

	BeforeEach(func() {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}

		if err := schema.Create(db, "schema.sql"); err != nil {
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
		"Repository",
		outboxtests.RepositorySuite(
			func() persistence.DataStore {
				return axmysql.NewDataStore(db)
			},
			func() outbox.Repository {
				return Repository{}
			},
		),
	)
})
