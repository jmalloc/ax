package saga

import (
	"context"
	"database/sql"
	"fmt"

	mysqlpersistence "github.com/jmalloc/ax/axmysql/persistence"
	"github.com/jmalloc/ax/persistence"
	"github.com/jmalloc/ax/saga"
)

// SnapshotRepository is a MySQL-backed implementation of Ax's
// eventsourcing.SnapshotRepository interface.
type SnapshotRepository struct{}

// LoadSagaSnapshot loads the latest available snapshot from the store.
//
// It returns an error if a snapshot of this instance is found, but belongs to
// a different saga, as identified by pk, the saga's persistence key.
func (SnapshotRepository) LoadSagaSnapshot(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	id saga.InstanceID,
) (saga.Instance, bool, error) {
	tx := mysqlpersistence.ExtractTx(ptx)

	var (
		cpk         string
		i           saga.Instance
		contentType string
		data        []byte
	)

	err := tx.QueryRowContext(
		ctx,
		`SELECT
			instance_id,
			revision,
			persistence_key,
			content_type,
			data
		FROM ax_saga_snapshot
		WHERE instance_id = ?
		ORDER BY revision DESC
		LIMIT 1`,
		id,
	).Scan(
		&i.InstanceID,
		&i.Revision,
		&cpk,
		&contentType,
		&data,
	)

	if err == sql.ErrNoRows {
		return saga.Instance{}, false, nil
	} else if err != nil {
		return saga.Instance{}, false, nil
	}

	if cpk != pk {
		return i, false, fmt.Errorf(
			"can not load saga snapshot of %s at revision %d for saga %s, it belongs to %s",
			i.InstanceID,
			i.Revision,
			pk,
			cpk,
		)
	}

	i.Data, err = saga.UnmarshalData(contentType, data)

	return i, true, err
}

// SaveSagaSnapshot saves a snapshot to the store.
//
// This implementation does not verify the saga's persistence key against
// existing snapshots of the same instance.
func (SnapshotRepository) SaveSagaSnapshot(
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

	descr := i.Data.InstanceDescription()
	// Truncate the message to 255 characters to fit within the column, if
	// required.
	if len(descr) > 255 {
		descr = descr[:255]
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO ax_saga_snapshot SET
			instance_id = ?,
			revision = ?,
			persistence_key = ?,
			description = ?,
			content_type = ?,
			data = ?`,
		i.InstanceID,
		i.Revision,
		pk,
		descr,
		contentType,
		data,
	)

	return err
}

// DeleteSagaSnapshots deletes any snapshots associated with a saga instance.
//
// This implementation does not verify the saga's persistence key. It simply
// ignores any snapshots that match the instance ID, but not the persistence key.
func (SnapshotRepository) DeleteSagaSnapshots(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	id saga.InstanceID,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	_, err := tx.ExecContext(
		ctx,
		`DELETE FROM ax_saga_snapshot
		WHERE persistence_key = ?
		AND instance_id = ?`,
		pk,
		id,
	)

	return err
}
