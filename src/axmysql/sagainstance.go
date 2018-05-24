package axmysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// SagaInstanceRepository is an implementation of saga.InstanceRepository that
// uses SQL persistence.
type SagaInstanceRepository struct{}

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
func (r SagaInstanceRepository) LoadSagaInstance(
	ctx context.Context,
	ptx persistence.Tx,
	id saga.InstanceID,
) (saga.Instance, error) {
	var (
		i           saga.Instance
		contentType string
		data        []byte
	)

	if err := sqlTx(ptx).QueryRowContext(
		ctx,
		`SELECT
			instance_id,
			revision,
			content_type,
			data
		FROM saga_instance
		WHERE instance_id = ?`,
		id,
	).Scan(
		&i.InstanceID,
		&i.Revision,
		&contentType,
		&data,
	); err != nil {
		return saga.Instance{}, err
	}

	var err error
	i.Data, err = saga.UnmarshalData(contentType, data)

	return i, err
}

// SaveSagaInstance persists a saga instance and its associated mapping
// table to the store as part of tx.
//
// It returns an error if the saga instance has been modified since it was
// loaded, or if there is a problem communicating with the store itself.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r SagaInstanceRepository) SaveSagaInstance(
	ctx context.Context,
	ptx persistence.Tx,
	i saga.Instance,
) error {
	tx := sqlTx(ptx)

	contentType, data, err := saga.MarshalData(i.Data)
	if err != nil {
		return err
	}

	if i.Revision == 0 {
		return r.insertInstance(ctx, tx, i, contentType, data)
	}

	return r.updateInstance(ctx, tx, i, contentType, data)
}

// insertInstance inserts a new saga instance.
func (SagaInstanceRepository) insertInstance(
	ctx context.Context,
	tx *sql.Tx,
	i saga.Instance,
	contentType string,
	data []byte,
) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO saga_instance SET
			instance_id = ?,
			revision = 1,
			description = ?,
			content_type = ?,
			data = ?`,
		i.InstanceID,
		i.Data.SagaDescription(),
		contentType,
		data,
	)

	// TODO: return a more meaningful error if we get a duplicate key error

	return err
}

// updateInstance updates an existing saga instance.
// It returns an error if i.Revision is not the current revision.
func (SagaInstanceRepository) updateInstance(
	ctx context.Context,
	tx *sql.Tx,
	i saga.Instance,
	contentType string,
	data []byte,
) error {
	res, err := tx.ExecContext(
		ctx,
		`UPDATE saga_instance SET
			revision = revision + 1,
			description = ?,
			content_type = ?,
			data = ?
		WHERE instance_id = ?
		AND revision = ?`,
		i.Data.SagaDescription(),
		contentType,
		data,
		i.InstanceID,
		i.Revision,
	)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf(
			"can not update saga instance %s, revision %d is not the current revision",
			i.InstanceID,
			i.Revision,
		)
	}

	return nil
}
