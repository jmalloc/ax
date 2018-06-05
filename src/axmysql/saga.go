package axmysql

import (
	"github.com/jmalloc/ax/src/ax/saga/mapping/keyset"
	"github.com/jmalloc/ax/src/ax/saga/persistence/eventsourcing"
	"github.com/jmalloc/ax/src/axmysql/saga"
)

// SnapshotRepository is a saga snapshot repository backed by an MySQL database.
var SnapshotRepository eventsourcing.SnapshotRepository = saga.SnapshotRepository{}

// KeySetRepository is a key-set repository backed by an MySQL database.
var KeySetRepository keyset.Repository = saga.KeySetRepository{}
