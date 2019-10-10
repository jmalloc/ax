package projection_test

import "github.com/jmalloc/ax/projection"
import . "github.com/jmalloc/ax/axmysql/projection"

var _ projection.Projector = (*ReadModelProjector)(nil) // ensure ReadModelProjector implements Projector
