package saga

import (
	"context"
	"database/sql"
	"fmt"

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
// It returns an error if the instance is found, but belongs to a different
// saga, as identified by pk, the saga's persistence key.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r CRUDRepository) LoadSagaInstance(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	id saga.InstanceID,
) (saga.Instance, bool, error) {
	var (
		ipk         string
		i           saga.Instance
		contentType string
		data        []byte
	)

	if err := mysqlpersistence.ExtractTx(ptx).QueryRowContext(
		ctx,
		`SELECT
			instance_id,
			persistence_key,
			revision,
			content_type,
			data
		FROM ax_saga_instance
		WHERE instance_id = ?`,
		id,
	).Scan(
		&i.InstanceID,
		&ipk,
		&i.Revision,
		&contentType,
		&data,
	); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return saga.Instance{}, false, err
	}

	if ipk != pk {
		return i, false, fmt.Errorf(
			"can not load saga instance %s for saga %s, it belongs to %s",
			id,
			pk,
			ipk,
		)
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
// It returns an error if the instance already exists, but belongs to a
// different saga, as identified by pk, the saga's persistence key.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r CRUDRepository) SaveSagaInstance(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	i saga.Instance,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	contentType, data, err := saga.MarshalData(i.Data)
	if err != nil {
		return err
	}

	if i.Revision == 0 {
		return r.insertInstance(ctx, tx, pk, i, contentType, data)
	}

	return r.updateInstance(ctx, tx, pk, i, contentType, data)
}

// insertInstance inserts a new saga instance.
func (CRUDRepository) insertInstance(
	ctx context.Context,
	tx *sql.Tx,
	pk string,
	i saga.Instance,
	contentType string,
	data []byte,
) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO ax_saga_instance SET
			instance_id = ?,
			persistence_key = ?,
			revision = 1,
			description = ?,
			content_type = ?,
			data = ?`,
		i.InstanceID,
		pk,
		i.Data.InstanceDescription(),
		contentType,
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
	pk string,
	i saga.Instance,
	contentType string,
	data []byte,
) error {
	var (
		ipk string
		rev saga.Revision
	)
	err := tx.QueryRowContext(
		ctx,
		`SELECT
			persistence_key,
			revision
		FROM ax_saga_instance
		WHERE instance_id = ?
		FOR UPDATE`,
		i.InstanceID,
	).Scan(
		&ipk,
		&rev,
	)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if i.Revision != rev {
		return fmt.Errorf(
			"can not update saga instance %s, revision %d is not the current revision",
			i.InstanceID,
			i.Revision,
		)
	}

	if pk != ipk {
		return fmt.Errorf(
			"can not save saga instance %s for saga %s, it belongs to %s",
			i.InstanceID,
			pk,
			ipk,
		)
	}

	res, err := tx.ExecContext(
		ctx,
		`UPDATE ax_saga_instance SET
			revision = revision + 1,
			description = ?,
			content_type = ?,
			data = ?
		WHERE instance_id = ?`,
		i.Data.InstanceDescription(),
		contentType,
		data,
		i.InstanceID,
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
			"no rows affected when updating saga instance %s at revision %d",
			i.InstanceID,
			i.Revision,
		)
	}

	return nil
}
