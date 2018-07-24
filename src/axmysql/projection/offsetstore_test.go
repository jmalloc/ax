package projection_test

import (
	"database/sql"
	"os"

	"github.com/jmalloc/ax/src/ax/projection"
	"github.com/jmalloc/ax/src/axmysql"
	"github.com/jmalloc/ax/src/axmysql/internal/schema"

	axmysqlprojection "github.com/jmalloc/ax/src/axmysql/projection"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axtest/projectiontests"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Projection Offset Store", func() {
	dsn := os.Getenv("AX_MYSQL_DSN")
	var db *sql.DB

	BeforeEach(func() {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}

		if err := schema.Create(db, "offsetstore.sql"); err != nil {
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
		"Projection Offset Store",
		projectiontests.OffsetStoreSuite(
			func() persistence.DataStore {
				return axmysql.NewDataStore(db)
			},
			func() projection.OffsetStore {
				return axmysqlprojection.OffsetStore{}
			},
		),
	)
})
