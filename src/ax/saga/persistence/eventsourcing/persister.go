package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
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
	MessageStore      persistence.MessageStore
	Snapshots         SnapshotRepository
	SnapshotFrequency saga.Revision
}

// BeginCreate starts a new unit-of-work that persists a new saga instance.
func (p *Persister) BeginCreate(
	ctx context.Context,
	_ saga.Saga,
	tx persistence.Tx,
	s ax.Sender,
	i saga.Instance,
) (saga.UnitOfWork, error) {
	return p.newUnitOfWork(tx, s, i), nil
}

// BeginUpdate starts a new unit-of-work that updates an existing saga
// instance.
func (p *Persister) BeginUpdate(
	ctx context.Context,
	sg saga.Saga,
	tx persistence.Tx,
	s ax.Sender,
	id saga.InstanceID,
) (saga.UnitOfWork, error) {
	var (
		i  saga.Instance
		ok bool
	)

	if p.Snapshots != nil {
		var err error
		i, ok, err = p.Snapshots.LoadSagaSnapshot(ctx, tx, id)
		if err != nil {
			return nil, err
		}
	}

	if !ok {
		i = saga.Instance{
			InstanceID: id,
			Data:       sg.NewData(),
		}
	}

	if err := applyEvents(
		ctx,
		tx,
		p.MessageStore,
		sg.(saga.EventedSaga),
		&i,
	); err != nil {
		return nil, err
	}

	return p.newUnitOfWork(tx, s, i), nil
}

// newUnitOfWork returns a new unit-of-work.
func (p *Persister) newUnitOfWork(
	tx persistence.Tx,
	s ax.Sender,
	i saga.Instance,
) *unitOfWork {
	return &unitOfWork{
		p.MessageStore,
		p.Snapshots,
		p.SnapshotFrequency,
		tx,
		&Recorder{Next: s},
		i,
	}
}

// unitOfWork is an implementation of saga.UnitOfWork that perists saga
// instances as a stream of events, with optional snapshots.
type unitOfWork struct {
	messageStore persistence.MessageStore
	snapshots    SnapshotRepository
	frequency    saga.Revision

	tx       persistence.Tx
	recorder *Recorder
	instance saga.Instance
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

	before := w.instance.Revision
	w.instance.Revision += saga.Revision(n)

	if w.shouldSnapshot(before, w.instance.Revision) {
		if err := w.snapshots.SaveSagaSnapshot(ctx, w.tx, w.instance); err != nil {
			return false, err
		}
	}

	return true, nil
}

// shouldSnapshot returns true if a new snapshot should be stored.
func (w *unitOfWork) shouldSnapshot(before, after saga.Revision) bool {
	if w.snapshots == nil {
		return false
	}

	freq := w.frequency
	if freq == 0 {
		freq = DefaultSnapshotFrequency
	}

	return (before / freq) != (after / freq)
}
