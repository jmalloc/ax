package saga

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"

	bolt "github.com/coreos/bbolt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// SnapshotRepository is a Bolt-backed implementation of Ax's
// eventsourcing.SnapshotRepository interface.
type SnapshotRepository struct{}

// SnapshotBktName is name of the Bolt root bucket where all saga snapshot data
// is stored.
var SnapshotBktName = []byte("ax_saga_snapshot")

// LoadSagaSnapshot loads the latest available snapshot from the store.
//
// It returns an error if a snapshot of this instance is found, but belongs to
// a different saga, as identified by pk, the saga's persistence key.
func (SnapshotRepository) LoadSagaSnapshot(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	id saga.InstanceID,
) (saga.Instance, bool, error) {
	tx := boltpersistence.ExtractTx(ptx)
	bkt := tx.Bucket([]byte("ax_saga_snapshot"))
	if bkt == nil {
		return saga.Instance{}, false, nil
	}

	bkt = bkt.Bucket([]byte(id.Get()))
	if bkt == nil {
		return saga.Instance{}, false, nil
	}

	k, pb := bkt.Cursor().Last()
	if k != nil && pb == nil {
		return saga.Instance{}, false, nil
	}

	var sn SagaSnapshot
	if err := proto.Unmarshal(pb, &sn); err != nil {
		return saga.Instance{}, false, err
	}

	i := saga.Instance{
		Revision: saga.Revision(sn.GetRevision()),
	}

	if err := i.InstanceID.Parse(sn.GetInstanceId()); err != nil {
		return saga.Instance{}, false, err
	}

	var x ptypes.DynamicAny
	if err := ptypes.UnmarshalAny(sn.Data, &x); err != nil {
		return saga.Instance{}, false, err
	}
	i.Data, _ = x.Message.(saga.Data)

	if sn.GetPersistenceKey() != pk {
		return i, false, fmt.Errorf(
			"can not load saga snapshot of %s at revision %d for saga %s, it belongs to %s",
			i.InstanceID,
			i.Revision,
			pk,
			sn.GetPersistenceKey(),
		)
	}
	return i, true, nil
}

// SaveSagaSnapshot saves a snapshot to the store.
//
// This implementation does not verify the saga's persistence key against
// existing snapshots of the same instance.
func (SnapshotRepository) SaveSagaSnapshot(
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
	sn := &SagaSnapshot{
		InstanceId:     i.InstanceID.Get(),
		Revision:       int64(i.Revision),
		PersistenceKey: pk,
		InsertTime:     time.Now().Format(time.RFC3339Nano),
	}

	sn.Data, err = ptypes.MarshalAny(i.Data)
	if err != nil {
		return err
	}

	bkt, err = tx.CreateBucketIfNotExists([]byte("ax_saga_snapshot"))
	if err != nil {
		return err
	}

	bkt, err = bkt.CreateBucketIfNotExists([]byte(i.InstanceID.Get()))
	if err != nil {
		return err
	}

	k := make([]byte, 8)
	binary.PutVarint(k, int64(sn.GetRevision()))

	return boltutil.MarshalProto(bkt, k, sn)
}

// DeleteSagaSnapshots deletes any snapshots associated with a saga instance.
//
// This implementation does not verify the saga's persistence key. It locates a
// child bucket indexed with id as a key and deletes it.
func (SnapshotRepository) DeleteSagaSnapshots(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	id saga.InstanceID,
) error {
	tx := boltpersistence.ExtractTx(ptx)

	bkt := tx.Bucket([]byte("ax_saga_snapshot"))
	if bkt == nil {
		return nil
	}

	return bkt.DeleteBucket([]byte(id.Get()))
}
