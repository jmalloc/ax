package saga

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
	mysqlpersistence "github.com/jmalloc/ax/src/axmysql/persistence"
)

// CRUDRepository is a MySQL-backed implementation of Ax's crud.Repository
// interface.
type CRUDRepository struct{}

// LoadSagaInstance fetches a saga instance by its ID.
//
// It returns an false if the instance does not exist. It returns an error
// if a problem occurs with the underlying data store.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r CRUDRepository) LoadSagaInstance(
	ctx context.Context,
	ptx persistence.Tx,
	id saga.InstanceID,
) (saga.Instance, bool, error) {
	var (
		i           saga.Instance
		contentType string
		data        []byte
	)

	if err := mysqlpersistence.ExtractTx(ptx).QueryRowContext(
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
		if err == sql.ErrNoRows {
			err = nil
		}

		return saga.Instance{}, false, err
	}

	var err error
	i.Data, err = saga.UnmarshalData(contentType, data)

	return i, true, err
}

// SaveSagaInstance persists a saga instance.
//
// It returns an error if i.Revision is not the current revision of the
// instance as it exists within the store, or a problem occurs with the
// underlying data store.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r CRUDRepository) SaveSagaInstance(
	ctx context.Context,
	ptx persistence.Tx,
	i saga.Instance,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

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
func (CRUDRepository) insertInstance(
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
			time = ?,
			data = ?`,
		i.InstanceID,
		i.Data.InstanceDescription(),
		contentType,
		time.Now().Format(time.RFC3339Nano),
		data,
	)

	// TODO: return a more meaningful error if we get a duplicate key error

	return err
}

// updateInstance updates an existing saga instance.
// It returns an error if i.Revision is not the current revision.
func (CRUDRepository) updateInstance(
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
			time = ?,
			data = ?
		WHERE instance_id = ?
		AND revision = ?`,
		i.Data.InstanceDescription(),
		contentType,
		time.Now().Format(time.RFC3339Nano),
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
