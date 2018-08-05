package axprotavo

import (
	"github.com/jmalloc/ax/src/ax/projection"
	protavoprojection "github.com/jmalloc/ax/src/axprotavo/projection"
)

// ProjectionOffsetStore is an offset store backed by a Protavo database.
var ProjectionOffsetStore projection.OffsetStore = protavoprojection.OffsetStore{}

// NewReadModelProjector returns a new projector that builds a Protavo based
// read-model from a stream of events.
func NewReadModelProjector(rm protavoprojection.ReadModel) projection.Projector {
	return protavoprojection.NewReadModelProjector(rm)
}
