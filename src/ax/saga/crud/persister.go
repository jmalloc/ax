package crud

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// Persister is an implementation of saga.Persister that persists saga instances
// using "CRUD" semantics.
type Persister struct {
	Repository Repository
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
	_ saga.Saga,
	tx persistence.Tx,
	s ax.Sender,
	id saga.InstanceID,
) (saga.UnitOfWork, error) {
	i, err := p.Repository.LoadSagaInstance(ctx, tx, id)
	if err != nil {
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
		p.Repository,
		tx,
		s,
		proto.Clone(i.Data).(saga.Data),
		i,
	}
}

// unitOfWork is an implementation of saga.UnitOfWork that saves saga instances
// using "CRUD" semantics.
type unitOfWork struct {
	repository Repository

	tx       persistence.Tx
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

	return true, w.repository.SaveSagaInstance(ctx, w.tx, w.instance)
}
