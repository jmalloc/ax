package aggregate

import "github.com/jmalloc/ax/src/ax"

// Recorder is a function that records the events produced by an aggregate.
type Recorder func(ax.Event)
