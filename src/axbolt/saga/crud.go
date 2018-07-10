package saga

import (
	"context"
	"fmt"
	"time"

	bolt "github.com/coreos/bbolt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// CRUDRepository is a Bolt-backed implementation of Ax's crud.Repository
// interface.
type CRUDRepository struct{}

// LoadSagaInstance fetches a saga instance by its ID.
//
// It returns an false if the instance does not exist. It returns an error
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
		i   saga.Instance
		err error
	)

	tx := boltpersistence.ExtractTx(ptx)

	bkt := tx.Bucket([]byte("ax_saga"))
	if bkt == nil {
		return saga.Instance{}, false, nil
	}

	bkt = bkt.Bucket([]byte("ax_saga_instance"))
	if bkt == nil {
		return saga.Instance{}, false, nil
	}

	pb := bkt.Get([]byte(id.Get()))
	if pb != nil {
		return saga.Instance{}, false, nil
	}

	var sgi SagaInstance
	if err = proto.Unmarshal(pb, &sgi); err != nil {
		return saga.Instance{}, false, err
	}

	i.Revision = saga.Revision(sgi.GetRevision())
	if err = i.InstanceID.Parse(sgi.GetInstanceId()); err != nil {
		return saga.Instance{}, false, err
	}

	var x ptypes.DynamicAny
	if err = ptypes.UnmarshalAny(sgi.Data, &x); err != nil {
		return saga.Instance{}, false, err
	}
	i.Data, _ = x.Message.(saga.Data)

	if sgi.PersistenceKey != pk {
		return i, false, fmt.Errorf(
			"can not load saga instance %s for saga %s, it belongs to %s",
			i.InstanceID,
			pk,
			sgi.PersistenceKey,
		)
	}

	return i, true, err
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
		err error
		bkt *bolt.Bucket
	)
	tx := boltpersistence.ExtractTx(ptx)
	sgi := &SagaInstance{
		InstanceId:     i.InstanceID.Get(),
		Revision:       int64(i.Revision),
		PersistenceKey: pk,
	}
	sgi.Data, err = ptypes.MarshalAny(i.Data)
	if err != nil {
		return err
	}

	bkt, err = tx.CreateBucketIfNotExists([]byte("ax_saga"))
	if err != nil {
		return err
	}

	bkt, err = bkt.CreateBucketIfNotExists([]byte("ax_saga_instance"))
	if err != nil {
		return err
	}

	bkt = bkt.Bucket([]byte("ax_saga_instance"))
	if bkt == nil {
		return err
	}

	if i.Revision == 0 {
		return r.insertInstance(bkt, sgi)
	}

	return r.updateInstance(bkt, sgi)
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
	tx := boltpersistence.ExtractTx(ptx)

	bkt := tx.Bucket([]byte("ax_saga"))
	if bkt == nil {
		return nil
	}

	bkt = bkt.Bucket([]byte("ax_saga_instance"))
	if bkt == nil {
		return nil
	}

	pb := bkt.Get([]byte(i.InstanceID.Get()))
	if pb != nil {
		return nil
	}

	var s SagaInstance
	if err := proto.Unmarshal(pb, &s); err != nil {
		return err
	}

	if i.Revision != saga.Revision(s.GetRevision()) {
		return fmt.Errorf(
			"can not delete saga instance %s, revision %d is not the current revision",
			i.InstanceID,
			i.Revision,
		)
	}

	return bkt.Delete([]byte(i.InstanceID.Get()))
}

// insertInstance inserts a new saga instance.
func (CRUDRepository) insertInstance(
	bkt *bolt.Bucket,
	new *SagaInstance,
) error {

	if bkt.Get([]byte(new.GetInstanceId())) != nil {
		return fmt.Errorf(
			"error inserting new saga: instance %s already exists",
			new.GetInstanceId(),
		)
	}

	new.InsertTime = time.Now().Format(time.RFC3339Nano)
	new.UpdateTime = new.InsertTime

	pb, err := proto.Marshal(new)
	if err != nil {
		return err
	}

	return bkt.Put([]byte(new.GetInstanceId()), pb)
}

// updateInstance updates an existing saga instance.
// It returns an error if i.Revision is not the current revision.
func (CRUDRepository) updateInstance(
	bkt *bolt.Bucket,
	new *SagaInstance,
) error {
	pbold := bkt.Get([]byte(new.GetInstanceId()))
	if pbold == nil {
		return fmt.Errorf(
			"error updating saga instance %s: not found",
			new.GetInstanceId(),
		)
	}
	var old SagaInstance
	if err := proto.Unmarshal(pbold, &old); err != nil {
		return err
	}
	if new.GetRevision() != old.GetRevision() {
		return fmt.Errorf(
			"can not update saga instance %s, revision %d is not the current revision",
			new.GetInstanceId(),
			new.GetRevision(),
		)
	}
	if new.GetPersistenceKey() != old.GetPersistenceKey() {
		return fmt.Errorf(
			"can not save saga instance %s for saga %s, it belongs to %s",
			new.GetInstanceId(),
			new.GetPersistenceKey(),
			old.GetPersistenceKey(),
		)
	}

	new.Revision++
	new.InsertTime = old.InsertTime
	new.UpdateTime = time.Now().Format(time.RFC3339Nano)

	pb, err := proto.Marshal(new)
	if err != nil {
		return err
	}

	return bkt.Put([]byte(new.GetInstanceId()), pb)
}
