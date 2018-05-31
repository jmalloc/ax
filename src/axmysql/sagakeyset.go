package axmysql

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// SagaKeySetRepository is an implementation of saga.KeySetRepository that uses
// SQL persistence.
type SagaKeySetRepository struct{}

// FindByKey returns the ID of the saga instance that contains k in its
// key set for the saga named sn.
//
// ok is false if no saga instance has a key set containing k.
func (SagaKeySetRepository) FindByKey(
	ctx context.Context,
	ptx persistence.Tx,
	sn, k string,
) (id saga.InstanceID, ok bool, err error) {
	err = sqlTx(ptx).QueryRowContext(
		ctx,
		`SELECT
			instance_id
		FROM saga_keyset
		WHERE saga = ?
		AND mapping_key = ?`,
		sn,
		k,
	).Scan(
		&id,
	)

	if err == nil {
		ok = true
	} else if err == sql.ErrNoRows {
		err = nil
	}

	return
}

// SaveKeys associates a key set with the saga instance identified by id
// for the saga named sn.
//
// Key sets must be disjoint. That is, no two instances of the same saga
// may share any keys.
func (SagaKeySetRepository) SaveKeys(
	ctx context.Context,
	ptx persistence.Tx,
	sn string,
	id saga.InstanceID,
	ks []string,
) error {
	tx := sqlTx(ptx)

	if _, err := tx.ExecContext(
		ctx,
		`DELETE FROM saga_keyset
		WHERE instance_id = ?`,
		id,
	); err != nil {
		return err
	}

	for _, k := range ks {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO saga_keyset SET
				saga = ?,
				mapping_key = ?,
				instance_id = ?`,
			sn,
			k,
			id,
		); err != nil {
			// TODO: return a more meaningful error if we get a duplicate key error
			return err
		}
	}

	return nil
}
