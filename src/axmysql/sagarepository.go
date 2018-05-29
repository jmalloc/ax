package axmysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// SagaRepository is an implementation of crud.Repository that stores saga
// instances using SQL persistence.
type SagaRepository struct{}

// LoadSagaInstance fetches a saga instance by its ID.
//
// It returns an error if the instance does not exist, or a problem occurs
// with the underlying data store.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r SagaRepository) LoadSagaInstance(
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

// SaveSagaInstance persists a saga instance.
//
// It returns an error if i.Revision is not the current revision of the
// instance as it exists within the store, or a problem occurs with the
// underlying data store.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r SagaRepository) SaveSagaInstance(
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
func (SagaRepository) insertInstance(
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
		i.Data.InstanceDescription(),
		contentType,
		data,
	)

	// TODO: return a more meaningful error if we get a duplicate key error

	return err
}

// updateInstance updates an existing saga instance.
// It returns an error if i.Revision is not the current revision.
func (SagaRepository) updateInstance(
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
		i.Data.InstanceDescription(),
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
