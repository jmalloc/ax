package axmysql

import (
	"github.com/jmalloc/ax/src/ax/projection"
	mysqlprojection "github.com/jmalloc/ax/src/axmysql/projection"
)

// ProjectionOffsetStore is an offset store backed by an MySQL database.
var ProjectionOffsetStore projection.OffsetStore = mysqlprojection.OffsetStore{}
