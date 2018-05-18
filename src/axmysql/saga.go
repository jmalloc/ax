package axmysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// SagaRepository is an implementation of saga.Repository that uses SQL
// persistence.
//
// It requires that Saga instances be implemented as protocol buffers messages.
type SagaRepository struct{}

// LoadSagaInstance fetches a saga instance that has a specific key/value
// pair in its mapping table.
//
// sn is the saga name. k and v are the key and value in the mapping table,
// respectively.
//
// If a saga instance is found; ok is true, otherwise it is false. A
// non-nil error indicates a problem with the store itself.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (*SagaRepository) LoadSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	sn, k string,
) (i saga.Instance, ok bool, err error) {
	stx := tx.(*Tx).sqlTx

	row := stx.QueryRowContext(
		ctx,
		`SELECT
			i.id,
			i.revision,
			i.content_type,
			i.data
		FROM saga_instance AS i
		INNER JOIN saga_key AS k
		ON k.instance_id = i.id
		WHERE k.saga_name = ?
		AND k.mapping_key = ?`,
		sn,
		k,
	)

	var (
		ct   string
		data []byte
	)

	err = row.Scan(
		&i.InstanceID,
		&i.Revision,
		&ct,
		&data,
	)

	if err == sql.ErrNoRows {
		return i, false, nil
	}

	if err != nil {
		return i, false, err
	}

	i.Data, err = saga.UnmarshalData(ct, data)
	if err != nil {
		return i, false, err
	}

	return i, true, nil
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
	sn string,
	i saga.Instance,
	ks saga.KeySet,
) error {
	stx := tx.(*Tx).sqlTx

	ct, data, err := saga.MarshalData(i.Data)
	if err != nil {
		return err
	}

	if i.Revision == 0 {
		if _, err := stx.ExecContext(
			ctx,
			`INSERT INTO saga_instance SET
				id = ?,
				saga_name = ?,
				description = ?,
				content_type = ?,
				data = ?,
				revision = 1`,
			i.InstanceID,
			sn,
			i.Data.SagaDescription(),
			ct,
			data,
		); err != nil {
			return err
		}
	} else {
		row := stx.QueryRowContext(
			ctx,
			`SELECT
				revision
			FROM saga_instance
			WHERE id = ?
			FOR UPDATE`,
			i.InstanceID,
		)

		var rev uint64
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
				description = ?,
				content_type = ?,
				data = ?,
				revision = revision + 1
			WHERE id = ?`,
			i.Data.SagaDescription(),
			ct,
			data,
			i.InstanceID,
		); err != nil {
			return err
		}

		if _, err := stx.ExecContext(
			ctx,
			`DELETE FROM saga_key
			WHERE instance_id = ?`,
			i.InstanceID,
		); err != nil {
			return err
		}
	}

	for k := range ks {
		if _, err := stx.ExecContext(
			ctx,
			`INSERT INTO saga_key SET
				saga_name = ?,
				mapping_key = ?,
				instance_id = ?`,
			sn,
			k,
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
		id             VARBINARY(255) NOT NULL PRIMARY KEY,
		saga_name      VARBINARY(255) NOT NULL,
		description    VARBINARY(255) NOT NULL,
		content_type   VARBINARY(255) NOT NULL,
		data           BLOB NOT NULL,
		revision       BIGINT UNSIGNED NOT NULL,

		INDEX (saga_name)
	)`,
	`CREATE TABLE IF NOT EXISTS saga_key (
		saga_name     VARBINARY(255) NOT NULL,
		mapping_key   VARBINARY(255) NOT NULL,
		instance_id   VARBINARY(255) NOT NULL,

		PRIMARY KEY (saga_name, mapping_key),
		INDEX (instance_id)
	)`,
}
