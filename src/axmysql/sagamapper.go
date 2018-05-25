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
//
// ok is false if no saga instance is found.
func (SagaMapper) FindByKey(
	ctx context.Context,
	ptx persistence.Tx,
	sn, k string,
) (id saga.InstanceID, ok bool, err error) {
	err = sqlTx(ptx).QueryRowContext(
		ctx,
		`SELECT
			instance_id
		FROM saga_map
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

// SaveKeys persists the changes to a saga instance's mapping key set.
//
// sn is the name of the saga.
func (SagaMapper) SaveKeys(
	ctx context.Context,
	ptx persistence.Tx,
	sn string,
	id saga.InstanceID,
	ks saga.KeySet,
) error {
	tx := sqlTx(ptx)

	if _, err := tx.ExecContext(
		ctx,
		`DELETE FROM saga_map
		WHERE instance_id = ?`,
		id,
	); err != nil {
		return err
	}

	for k := range ks {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO saga_map SET
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
