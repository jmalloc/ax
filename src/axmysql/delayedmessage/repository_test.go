package delayedmessage_test

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmalloc/ax/src/ax/delayedmessage"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axmysql"
	. "github.com/jmalloc/ax/src/axmysql/delayedmessage"
	"github.com/jmalloc/ax/src/axmysql/internal/schema"
	"github.com/jmalloc/ax/src/axtest/delayedmessagetests"
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
		delayedmessagetests.RepositorySuite(
			func() persistence.DataStore {
				return axmysql.NewDataStore(db)
			},
			func() delayedmessage.Repository {
				return Repository{}
			},
		),
	)
})
