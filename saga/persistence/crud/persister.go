package crud

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/persistence"
	"github.com/jmalloc/ax/saga"
)

// Persister is an implementation of saga.Persister that persists saga instances
// using "CRUD" semantics.
type Persister struct {
	Repository Repository
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
	pk := sg.PersistenceKey()
	i, ok, err := p.Repository.LoadSagaInstance(ctx, tx, pk, id)
	if err != nil {
		return nil, err
	}

	if !ok {
		i = saga.Instance{
			InstanceID: id,
			Data:       sg.NewData(),
		}
	}

	return &unitOfWork{
		p.Repository,
		tx,
		pk,
		s,
		proto.Clone(i.Data).(saga.Data),
		i,
	}, nil
}

// unitOfWork is an implementation of saga.UnitOfWork that saves saga instances
// using "CRUD" semantics.
type unitOfWork struct {
	repository Repository

	tx       persistence.Tx
	key      string
	sender   ax.Sender
	original saga.Data
	instance saga.Instance
}

// Sender returns the ax.Sender that the saga must use to send messages.
func (w *unitOfWork) Sender() ax.Sender {
	return w.sender
}

// Instance returns the saga instance that the unit-of-work applies to.
func (w *unitOfWork) Instance() saga.Instance {
	return w.instance
}

// Save persists changes to the instance.
// It returns true if any changes have occurred.
func (w *unitOfWork) Save(ctx context.Context) (bool, error) {
	if proto.Equal(w.instance.Data, w.original) {
		return false, nil
	}

	if err := w.repository.SaveSagaInstance(ctx, w.tx, w.key, w.instance); err != nil {
		return false, err
	}

	w.instance.Revision++

	return true, nil
}

// SaveAndComplete deletes the saga instance.
func (w *unitOfWork) SaveAndComplete(ctx context.Context) error {
	if err := w.repository.DeleteSagaInstance(ctx, w.tx, w.key, w.instance); err != nil {
		return err
	}

	w.instance.Revision++

	return nil
}

// Close is called when the unit-of-work has ended, regardless of whether
// Save() has been called.
func (w *unitOfWork) Close() {
}
