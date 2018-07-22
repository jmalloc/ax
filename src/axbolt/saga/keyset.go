package saga

import (
	"context"

	bolt "github.com/coreos/bbolt"
	"github.com/golang/protobuf/proto"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// KeySetRepository is a Bolt-backed implementation of Ax's keyset.Repository
// interface.
type KeySetRepository struct{}

// KeySetBktName is name of the Bolt root bucket where all saga keyset data is
// stored.
var KeySetBktName = []byte("ax_saga_keyset")

// FindByKey returns the ID of a saga instance that has a specific key in
// its key set.
//
// pk is the saga's persistence key, mk is the mapping key.
// ok is false if no saga instance has a key set containing mk.
func (KeySetRepository) FindByKey(
	ctx context.Context,
	ptx persistence.Tx,
	pk, mk string,
) (id saga.InstanceID, ok bool, err error) {
	tx := boltpersistence.ExtractTx(ptx)
	bkt := tx.Bucket(KeySetBktName)
	if bkt == nil {
		return
	}

	if bkt = bkt.Bucket([]byte(pk)); bkt == nil {
		return
	}

	if s := bkt.Get([]byte(mk)); s != nil {
		if err = id.Parse(string(s)); err != nil {
			return
		}
		ok = true
	}
	return
}

// SaveKeys associates a set of mapping keys with a saga instance.
//
// Key sets must be disjoint. That is, no two instances of the same saga
// may share any keys.
//
// pk is the saga's persistence key. ks is the set of mapping keys.
//
// SaveKeys() may panic if ks contains duplicate keys.
func (KeySetRepository) SaveKeys(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	ks []string,
	id saga.InstanceID,
) error {
	var (
		bkt *bolt.Bucket
		pb  []byte
		err error
	)
	tx := boltpersistence.ExtractTx(ptx)
	if bkt, err = tx.CreateBucketIfNotExists(KeySetBktName); err != nil {
		return err
	}

	if pb, err = proto.Marshal(&SagaKeySet{
		PersistenceKey: pk,
		MappingKeys:    ks,
	}); err != nil {
		return err
	}

	if err = bkt.Put([]byte(id.Get()), pb); err != nil {
		return err
	}

	bkt, err = bkt.CreateBucketIfNotExists([]byte(pk))
	if err != nil {
		return err
	}

	for _, k := range ks {
		if err = bkt.Put([]byte(k), []byte(id.Get())); err != nil {
			return err
		}
	}

	return nil
}

// DeleteKeys removes any mapping keys associated with a saga instance.
//
// pk is the saga's persistence key.
func (r KeySetRepository) DeleteKeys(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	id saga.InstanceID,
) error {
	var (
		bkt *bolt.Bucket
		err error
		pb  []byte
		ks  SagaKeySet
	)
	tx := boltpersistence.ExtractTx(ptx)
	if bkt = tx.Bucket(KeySetBktName); bkt == nil {
		return nil
	}

	if pb = bkt.Get([]byte(id.Get())); pb == nil {
		return nil
	}

	if err = proto.Unmarshal(pb, &ks); err != nil {
		return err
	}

	if err = bkt.Delete([]byte(id.Get())); err != nil {
		return err
	}

	if bkt = bkt.Bucket([]byte(ks.GetPersistenceKey())); bkt == nil {
		return nil
	}

	for _, k := range ks.GetMappingKeys() {
		if err = bkt.Delete([]byte(k)); err != nil {
			return err
		}
	}

	return nil
}
