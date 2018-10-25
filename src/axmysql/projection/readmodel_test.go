package projection_test

import "github.com/jmalloc/ax/src/ax/projection"
import . "github.com/jmalloc/ax/src/axmysql/projection"

var _ projection.Projector = (*ReadModelProjector)(nil) // ensure ReadModelProjector implements Projector
