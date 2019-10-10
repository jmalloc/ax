package axmysql

import (
	"github.com/jmalloc/ax/axmysql/saga"
	"github.com/jmalloc/ax/saga/mapping/keyset"
	"github.com/jmalloc/ax/saga/persistence/crud"
	"github.com/jmalloc/ax/saga/persistence/eventsourcing"
)

// SagaKeySetRepository is a key-set repository backed by a MySQL database.
var SagaKeySetRepository keyset.Repository = saga.KeySetRepository{}

// SagaCRUDRepository is a CRUD saga repository backed by a MySQL database.
var SagaCRUDRepository crud.Repository = saga.CRUDRepository{}

// SagaSnapshotRepository is a saga snapshot repository backed by a MySQL database.
var SagaSnapshotRepository eventsourcing.SnapshotRepository = saga.SnapshotRepository{}
