package axmysql

import (
	"github.com/jmalloc/ax/src/ax/projection"
	mysqlprojection "github.com/jmalloc/ax/src/axmysql/projection"
)

// ProjectionOffsetStore is an offset store backed by an MySQL database.
var ProjectionOffsetStore projection.OffsetStore = mysqlprojection.OffsetStore{}

// NewReadModelProjector returns a new projector that builds a MySQL based
// read-model from a stream of events.
func NewReadModelProjector(rm mysqlprojection.ReadModel) projection.Projector {
	return mysqlprojection.NewReadModelProjector(rm)
}
