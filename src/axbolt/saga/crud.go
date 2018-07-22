package saga

import (
	"context"
	"fmt"
	"time"

	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"

	bolt "github.com/coreos/bbolt"
	"github.com/golang/protobuf/ptypes"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// CRUDRepository is a Bolt-backed implementation of Ax's crud.Repository
// interface.
type CRUDRepository struct{}

// InstanceBktName is the name of of the Bolt root bucket where all saga instance
// data is stored.
var InstanceBktName = []byte("ax_saga_instance")

// LoadSagaInstance fetches a saga instance by its ID.
//
// It returns false if the instance does not exist. It returns an error
// if a problem occurs with the underlying data store.
//
// It returns an error if the instance is found, but belongs to a different
// saga, as identified by pk, the saga's persistence key.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r CRUDRepository) LoadSagaInstance(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	id saga.InstanceID,
) (saga.Instance, bool, error) {
	var (
		x ptypes.DynamicAny
		s SagaInstance
	)
	tx := boltpersistence.ExtractTx(ptx)
	bkt := tx.Bucket(InstanceBktName)
	if bkt == nil {
		return saga.Instance{}, false, nil
	}

	if ok, err := boltutil.UnmarshalProto(
		bkt,
		[]byte(id.Get()),
		&s,
	); err != nil || !ok {
		return saga.Instance{}, ok, err
	}

	i := saga.Instance{
		InstanceID: id,
		Revision:   saga.Revision(s.GetRevision()),
	}
	if err := ptypes.UnmarshalAny(s.Data, &x); err != nil {
		return saga.Instance{}, false, err
	}
	i.Data, _ = x.Message.(saga.Data)

	if s.PersistenceKey != pk {
		return i, false, fmt.Errorf(
			"can not load saga instance %s for saga %s, it belongs to %s",
			i.InstanceID,
			pk,
			s.GetPersistenceKey(),
		)
	}

	return i, true, nil
}

// SaveSagaInstance persists a saga instance.
//
// It returns an error if i.Revision is not the current revision of the
// instance as it exists within the store, or a problem occurs with the
// underlying data store.
//
// It returns an error if the instance already exists, but belongs to a
// different saga, as identified by pk, the saga's persistence key.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r CRUDRepository) SaveSagaInstance(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	i saga.Instance,
) error {
	var (
		err  error
		ok   bool
		bkt  *bolt.Bucket
		prev SagaInstance
	)
	tx := boltpersistence.ExtractTx(ptx)
	new := &SagaInstance{
		InstanceId:     i.InstanceID.Get(),
		Revision:       int64(i.Revision),
		PersistenceKey: pk,
	}
	new.Data, err = ptypes.MarshalAny(i.Data)
	if err != nil {
		return err
	}

	bkt, err = tx.CreateBucketIfNotExists(InstanceBktName)
	if err != nil {
		return err
	}

	if ok, err = boltutil.UnmarshalProto(
		bkt,
		[]byte(new.GetInstanceId()),
		&prev,
	); err != nil {
		return err
	}

	if ok {
		if new.GetRevision() != prev.GetRevision() {
			return fmt.Errorf(
				"can not update saga instance %s, revision %d is not the current revision",
				i.InstanceID,
				i.Revision,
			)
		}
		if pk != prev.GetPersistenceKey() {
			return fmt.Errorf(
				"can not save saga instance %s for saga %s, it belongs to %s",
				i.InstanceID,
				pk,
				prev.GetPersistenceKey(),
			)
		}

		return r.updateInstance(bkt, new, &prev)
	}

	return r.insertInstance(bkt, new)
}

// DeleteSagaInstance deletes a saga instance.
//
// It returns an error if i.Revision is not the current revision of the
// instance as it exists within the store, or a problem occurs with the
// underlying data store.
//
// It returns an error if the instance belongs to a different saga, as
// identified by pk, the saga's persistence key.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r CRUDRepository) DeleteSagaInstance(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	i saga.Instance,
) error {
	var (
		s   SagaInstance
		err error
		ok  bool
	)
	tx := boltpersistence.ExtractTx(ptx)
	bkt := tx.Bucket(InstanceBktName)
	if bkt == nil {
		return nil
	}

	if ok, err = boltutil.UnmarshalProto(
		bkt,
		[]byte(i.InstanceID.Get()),
		&s,
	); err != nil || !ok {
		return err
	}

	if i.Revision != saga.Revision(s.GetRevision()) {
		return fmt.Errorf(
			"can not delete saga instance %s, revision %d is not the current revision",
			i.InstanceID,
			i.Revision,
		)
	}

	if pk != s.GetPersistenceKey() {
		return fmt.Errorf(
			"can not save saga instance %s for saga %s, it belongs to %s",
			i.InstanceID,
			pk,
			s.GetPersistenceKey(),
		)
	}

	return bkt.Delete([]byte(i.InstanceID.Get()))
}

// insertInstance inserts a new saga instance.
func (CRUDRepository) insertInstance(
	bkt *bolt.Bucket,
	new *SagaInstance,
) error {

	new.Revision = 1
	new.InsertTime = time.Now().Format(time.RFC3339Nano)
	new.UpdateTime = new.InsertTime

	return boltutil.MarshalProto(
		bkt,
		[]byte(new.GetInstanceId()),
		new,
	)
}

// updateInstance updates an existing saga instance.
// It returns an error if i.Revision is not the current revision.
func (CRUDRepository) updateInstance(
	bkt *bolt.Bucket,
	new, prev *SagaInstance,
) error {

	new.Revision = prev.Revision + 1
	new.InsertTime = prev.InsertTime
	new.UpdateTime = time.Now().Format(time.RFC3339Nano)

	return boltutil.MarshalProto(
		bkt,
		[]byte(new.GetInstanceId()),
		new,
	)
}
