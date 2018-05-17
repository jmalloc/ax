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

// LoadSagaInstance fetches a saga instance from the store based on a
// mapping key for a particular message type.
//
// sn is the saga name. mt is the message type, and k is the mapping key
// associated with that message type. i must be a non-nil pointer to an
// empty saga instance, which is populated with the loaded data.
//
// ok is true if the instance is found, in which case i is populated with
// data from the store.
//
// err is non-nil if there is a problem communicating with the store itself.
func (*SagaRepository) LoadSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	req saga.LoadRequest,
) (saga.LoadResult, bool, error) {
	stx := tx.(*Tx).tx

	row := stx.QueryRowContext(
		ctx,
		`SELECT
			i.id,
			i.revision,
			i.content_type,
			i.data
		FROM saga_instance AS i
		INNER JOIN saga_map AS m
		ON m.instance_id = i.id
		WHERE m.saga_name = ?
		AND m.message_type = ?
		AND m.mapping_key = ?`,
		req.SagaName,
		req.MessageType.Name,
		req.MappingKey,
	)

	var (
		res  saga.LoadResult
		ct   string
		data []byte
	)

	err := row.Scan(
		&res.InstanceID,
		&res.CurrentRevision,
		&ct,
		&data,
	)

	if err == sql.ErrNoRows {
		return res, false, nil
	}

	if err != nil {
		return res, false, err
	}

	res.Instance, err = saga.UnmarshalInstance(ct, data)
	if err != nil {
		return res, false, err
	}

	return res, true, nil
}

// SaveSagaInstance persists a saga instance and its associated mapping
// table to the store as part of tx.
//
// It returns an error if the saga instance has been modified since it was
// loaded, or if there is a problem communicating with the store itself.
//
// sn is the saga name. i is the saga instance to save, and t is the complete
// mapping table for i.
//
// Save() panics if the repository is not able to enlist in tx because it
// uses a different underlying storage system.
func (*SagaRepository) SaveSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	req saga.SaveRequest,
) error {
	stx := tx.(*Tx).tx

	ct, data, err := saga.MarshalInstance(req.Instance)
	if err != nil {
		return err
	}

	if req.CurrentRevision == 0 {
		if _, err := stx.ExecContext(
			ctx,
			`INSERT INTO saga_instance SET
				id = ?,
				saga_name = ?,
				description = ?,
				content_type = ?,
				data = ?,
				revision = 1`,
			req.InstanceID,
			req.SagaName,
			req.Instance.InstanceDescription(),
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
			req.InstanceID,
		)

		var rev uint64
		if err := row.Scan(&rev); err != nil {
			return err
		}

		if req.CurrentRevision != rev {
			return fmt.Errorf(
				"can not update saga instance %s, revision is out of date",
				req.InstanceID,
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
			req.Instance.InstanceDescription(),
			ct,
			data,
			req.InstanceID,
		); err != nil {
			return err
		}

		if _, err := stx.ExecContext(
			ctx,
			`DELETE FROM saga_map
			WHERE instance_id = ?`,
			req.InstanceID,
		); err != nil {
			return err
		}
	}

	for mt, mk := range req.MappingTable {
		if _, err := stx.ExecContext(
			ctx,
			`INSERT INTO saga_map SET
				saga_name = ?,
				message_type = ?,
				mapping_key = ?,
				instance_id = ?`,
			req.SagaName,
			mt.Name,
			mk,
			req.InstanceID,
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
	`CREATE TABLE IF NOT EXISTS saga_map (
		saga_name     VARBINARY(255) NOT NULL,
		message_type  VARBINARY(255) NOT NULL,
		mapping_key   VARBINARY(255) NOT NULL,
		instance_id   VARBINARY(255) NOT NULL,

		PRIMARY KEY (saga_name, message_type, mapping_key),
		INDEX (instance_id)
	)`,
}
