package saga

import (
	"context"

	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

const (
	// keySetBktName is name of the Bolt root bucket where all saga keyset data is
	// stored.
	keySetBktName = "ax_saga_keyset"
)

// KeySetRepository is a Bolt-backed implementation of Ax's keyset.Repository
// interface.
type KeySetRepository struct{}

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
	b := boltutil.Get(tx, mk, keySetBktName, pk)
	if b == nil {
		return
	}
	if err = id.Parse(string(b)); err != nil {
		return
	}

	ok = true
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
	tx := boltpersistence.ExtractTx(ptx)
	err := boltutil.PutProto(
		tx,
		id.Get(),
		&SagaKeySet{
			PersistenceKey: pk,
			MappingKeys:    ks,
		},
		keySetBktName,
	)
	if err != nil {
		return err
	}

	for _, k := range ks {
		if err = boltutil.Put(
			tx,
			k,
			[]byte(id.Get()),
			keySetBktName,
			pk,
		); err != nil {
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
		err error
		ok  bool
		ks  SagaKeySet
	)
	tx := boltpersistence.ExtractTx(ptx)
	if ok, err = boltutil.GetProto(
		tx,
		id.Get(),
		&ks,
		keySetBktName,
	); err != nil {
		return err
	}

	if ok {
		for _, k := range ks.GetMappingKeys() {
			if err = boltutil.Delete(
				tx,
				k,
				keySetBktName,
				ks.GetPersistenceKey(),
			); err != nil {
				return err
			}
		}
	}
	return nil
}
