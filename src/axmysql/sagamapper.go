package axmysql

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// SagaMapper is an implementation of saga.Mapper that uses SQL persistence.
type SagaMapper struct{}

// FindByKey returns the instance ID of the saga instance that handles
// messages with a specific mapping key.
//
// sn is the name of the saga, and k is the message's mapping key.
func (SagaMapper) FindByKey(
	ctx context.Context,
	tx persistence.Tx,
	sn, k string,
) (saga.InstanceID, bool, error) {
	stx := tx.(*Tx).sqlTx

	row := stx.QueryRowContext(
		ctx,
		`SELECT k.instance_id
		FROM saga_keyset AS k
		WHERE k.saga = ?
		AND k.mapping_key = ?`,
		sn,
		k,
	)

	var id saga.InstanceID

	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return id, false, nil
		}

		return id, false, err
	}

	return id, true, nil
}

// SaveKeys persists the changes to a saga instance's mapping key set.
//
// sn is the name of the saga.
func (SagaMapper) SaveKeys(
	ctx context.Context,
	tx persistence.Tx,
	sn string,
	id saga.InstanceID,
	ks saga.KeySet,
) error {
	stx := tx.(*Tx).sqlTx

	if _, err := stx.ExecContext(
		ctx,
		`DELETE FROM saga_keyset
		WHERE instance_id = ?`,
		id,
	); err != nil {
		return err
	}

	for k := range ks {
		if _, err := stx.ExecContext(
			ctx,
			`INSERT INTO saga_keyset SET
				saga = ?,
				mapping_key = ?,
				instance_id = ?`,
			sn,
			k,
			id,
		); err != nil {
			return err
		}
	}

	return nil
}
