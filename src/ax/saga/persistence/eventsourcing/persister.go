package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/messagestore"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// DefaultSnapshotFrequency is the default number of revisions to allow between
// storing snapshots.
const DefaultSnapshotFrequency saga.Revision = 1000

// Persister is an implementation of saga.Persister that stores saga instances
// using event-sourcing semantics.
//
// The saga data MUST implement saga.EventedData.
type Persister struct {
	MessageStore      messagestore.Store
	Snapshots         SnapshotRepository
	SnapshotFrequency saga.Revision
}

// BeginUnitOfWork starts a new unit-of-work that modifies a saga instance.
//
// If the saga instance does not exist, it returns a UnitOfWork with an
// instance at revision zero.
func (p *Persister) BeginUnitOfWork(
	ctx context.Context,
	sg saga.Saga,
	tx persistence.Tx,
	s ax.Sender,
	id saga.InstanceID,
) (saga.UnitOfWork, error) {

	var (
		ok  bool
		err error
	)

	uow := &unitOfWork{
		messageStore: p.MessageStore,
		snapshots:    p.Snapshots,
		frequency:    p.SnapshotFrequency,
		tx:           tx,
		recorder:     &Recorder{Next: s},
	}

	if p.Snapshots != nil {
		uow.key = sg.PersistenceKey()
		uow.instance, ok, err = p.Snapshots.LoadSagaSnapshot(ctx, uow.tx, uow.key, id)
		if err != nil {
			return nil, err
		}
	}

	if !ok {
		uow.instance = saga.Instance{
			InstanceID: id,
			Data:       sg.NewData(),
		}
	}

	uow.lastKnownSnapshot = uow.instance.Revision

	if err = applyEvents(
		ctx,
		tx,
		p.MessageStore,
		sg.(saga.EventedSaga),
		&uow.instance,
	); err != nil {
		return nil, err
	}

	return uow, nil
}

// unitOfWork is an implementation of saga.UnitOfWork that perists saga
// instances as a stream of events, with optional snapshots.
type unitOfWork struct {
	messageStore messagestore.Store
	snapshots    SnapshotRepository
	frequency    saga.Revision

	lastKnownSnapshot saga.Revision
	tx                persistence.Tx
	key               string
	recorder          *Recorder
	instance          saga.Instance
}

// Sender returns the ax.Sender that the saga must use to send messages.
func (w *unitOfWork) Sender() ax.Sender {
	return w.recorder
}

// Instance returns the saga instance that the unit-of-work applies to.
func (w *unitOfWork) Instance() saga.Instance {
	return w.instance
}

// Save persists changes to the instance.
// It returns true if any changes have occurred.
func (w *unitOfWork) Save(ctx context.Context) (bool, error) {
	n := len(w.recorder.Events)

	if n == 0 {
		return false, nil
	}

	if err := appendEvents(
		ctx,
		w.tx,
		w.messageStore,
		w.instance,
		w.recorder.Events,
	); err != nil {
		return false, err
	}

	w.instance.Revision += saga.Revision(n)

	if w.shouldSnapshot() {
		if err := w.snapshots.SaveSagaSnapshot(ctx, w.tx, w.key, w.instance); err != nil {
			return false, err
		}
	}

	return true, nil
}

// Close is called when the unit-of-work has ended, regardless of whether
// Save() has been called.
func (w *unitOfWork) Close() {
}

// shouldSnapshot returns true if a new snapshot should be stored.
func (w *unitOfWork) shouldSnapshot() bool {
	if w.snapshots == nil {
		return false
	}

	freq := w.frequency
	if freq == 0 {
		freq = DefaultSnapshotFrequency
	}

	return (w.instance.Revision - w.lastKnownSnapshot) >= freq
}
