package axmysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

// SagaRepository is an implementation of saga.Repository that uses SQL
// persistence.
//
// It requires that Saga instances be implemented as protocol buffers messages.
type SagaRepository struct{}

// LoadSagaInstance fetches a saga instance that has a specific mapping key
// in its key set.
//
// sn is the saga name. k is the message mapping key.
//
// If a saga instance is found; ok is true, otherwise it is false. A
// non-nil error indicates a problem with the store itself.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (*SagaRepository) LoadSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	id saga.InstanceID,
) (i saga.Instance, err error) {
	stx := tx.(*Tx).sqlTx

	row := stx.QueryRowContext(
		ctx,
		`SELECT
			i.instance_id,
			i.current_revision,
			i.snapshot_revision,
			i.snapshot_content_type,
			i.snapshot_data
		FROM saga_instance AS i
		WHERE i.instance_id = ?`,
		id,
	)

	var (
		snapRev  saga.Revision
		snapType string
		snapData []byte
	)

	err = row.Scan(
		&i.InstanceID,
		&i.Revision,
		&snapRev,
		&snapType,
		&snapData,
	)

	if err != nil {
		return
	}

	if snapRev != i.Revision {
		err = errors.New("saga snapshot revision is not the current revision")
		return
	}

	i.Data, err = saga.UnmarshalData(snapType, snapData)
	return
}

// SaveSagaInstance persists a saga instance and its associated mapping
// table to the store as part of tx.
//
// It returns an error if the saga instance has been modified since it was
// loaded, or if there is a problem communicating with the store itself.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (*SagaRepository) SaveSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	i saga.Instance,
) error {
	stx := tx.(*Tx).sqlTx

	snapType, snapData, err := saga.MarshalData(i.Data)
	if err != nil {
		return err
	}

	if i.Revision == 0 {
		if _, err := stx.ExecContext(
			ctx,
			`INSERT INTO saga_instance SET
				instance_id = ?,
				current_revision = 1,
				snapshot_revision = 1,
				snapshot_description = ?,
				snapshot_content_type = ?,
				snapshot_data = ?`,
			i.InstanceID,
			i.Data.SagaDescription(),
			snapType,
			snapData,
		); err != nil {
			return err
		}
	} else {
		row := stx.QueryRowContext(
			ctx,
			`SELECT
				current_revision
			FROM saga_instance
			WHERE instance_id = ?
			FOR UPDATE`,
			i.InstanceID,
		)

		var rev saga.Revision
		if err := row.Scan(&rev); err != nil {
			return err
		}

		if i.Revision != rev {
			return fmt.Errorf(
				"can not update saga instance %s, revision is out of date",
				i.InstanceID,
			)
		}

		if _, err := stx.ExecContext(
			ctx,
			`UPDATE saga_instance SET
				current_revision = ?,
				snapshot_revision = ?,
				snapshot_description = ?,
				snapshot_content_type = ?,
				snapshot_data = ?
			WHERE instance_id = ?`,
			rev+1,
			rev+1,
			i.Data.SagaDescription(),
			snapType,
			snapData,
			i.InstanceID,
		); err != nil {
			return err
		}
	}

	return nil
}

// SagaSchema is a collection of DDL queries that create the schema
// used by OutboxRepository.
var SagaSchema = []string{
	`CREATE TABLE IF NOT EXISTS saga_instance (
		instance_id           VARBINARY(255) NOT NULL PRIMARY KEY,
		current_revision      BIGINT UNSIGNED NOT NULL,

		snapshot_revision     BIGINT UNSIGNED,
		snapshot_description  VARBINARY(255),
		snapshot_content_type VARBINARY(255),
		snapshot_data         BLOB,

		INDEX (saga)
	)`,
	`CREATE TABLE IF NOT EXISTS saga_keyset (
		saga          VARBINARY(255) NOT NULL,
		mapping_key   VARBINARY(255) NOT NULL,
		instance_id   VARBINARY(255) NOT NULL,

		PRIMARY KEY (saga, mapping_key),
		INDEX (instance_id)
	)`,
}
