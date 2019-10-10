package axmysql

import (
	mysqlprojection "github.com/jmalloc/ax/axmysql/projection"
	"github.com/jmalloc/ax/projection"
)

// ProjectionOffsetStore is an offset store backed by a MySQL database.
var ProjectionOffsetStore projection.OffsetStore = mysqlprojection.OffsetStore{}

// NewReadModelProjector returns a new projector that builds a MySQL based
// read-model from a stream of events.
func NewReadModelProjector(rm mysqlprojection.ReadModel) projection.Projector {
	return mysqlprojection.NewReadModelProjector(rm)
}
