package axbolt

import (
	"github.com/jmalloc/ax/src/ax/projection"
	boltprojection "github.com/jmalloc/ax/src/axbolt/projection"
)

// ProjectionOffsetStore is an offset store backed by a Bolt database.
var ProjectionOffsetStore projection.OffsetStore = boltprojection.OffsetStore{}

// NewReadModelProjector returns a new projector that builds a Bolt based
// read-model from a stream of events.
func NewReadModelProjector(rm boltprojection.ReadModel) projection.Projector {
	return boltprojection.NewReadModelProjector(rm)
}
