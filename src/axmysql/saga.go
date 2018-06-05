package axmysql

import (
	"github.com/jmalloc/ax/src/ax/saga/mapping/keyset"
	"github.com/jmalloc/ax/src/ax/saga/persistence/crud"
	"github.com/jmalloc/ax/src/ax/saga/persistence/eventsourcing"
	"github.com/jmalloc/ax/src/axmysql/saga"
)

// SagaKeySetRepository is a key-set repository backed by an MySQL database.
var SagaKeySetRepository keyset.Repository = saga.KeySetRepository{}

// SagaCRUDRepository is a CRUD saga repository backed by an MySQL database.
var SagaCRUDRepository crud.Repository = saga.CRUDRepository{}

// SagaSnapshotRepository is a saga snapshot repository backed by an MySQL database.
var SagaSnapshotRepository eventsourcing.SnapshotRepository = saga.SnapshotRepository{}
