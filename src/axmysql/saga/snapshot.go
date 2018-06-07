package saga

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
	mysqlpersistence "github.com/jmalloc/ax/src/axmysql/persistence"
)

// SnapshotRepository is a MySQL-backed implementation of Ax's
// eventsourcing.SnapshotRepository interface.
type SnapshotRepository struct{}

// LoadSagaSnapshot loads the latest available snapshot from the store.
func (SnapshotRepository) LoadSagaSnapshot(
	ctx context.Context,
	ptx persistence.Tx,
	id saga.InstanceID,
) (saga.Instance, bool, error) {
	var (
		i           saga.Instance
		contentType string
		data        []byte
	)

	err := mysqlpersistence.ExtractTx(ptx).QueryRowContext(
		ctx,
		`SELECT
			instance_id,
			revision,
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
		&contentType,
		&data,
	)

	if err == sql.ErrNoRows {
		return saga.Instance{}, false, nil
	}

	if err != nil {
		return saga.Instance{}, false, nil
	}

	i.Data, err = saga.UnmarshalData(contentType, data)

	return i, true, err
}

// SaveSagaSnapshot saves a snapshot to the store.
func (SnapshotRepository) SaveSagaSnapshot(
	ctx context.Context,
	tx persistence.Tx,
	i saga.Instance,
) error {
	contentType, data, err := saga.MarshalData(i.Data)
	if err != nil {
		return err
	}

	_, err = mysqlpersistence.ExtractTx(tx).ExecContext(
		ctx,
		`INSERT INTO ax_saga_snapshot SET
			instance_id = ?,
			revision = ?,
			description = ?,
			content_type = ?,
			data = ?`,
		i.InstanceID,
		i.Revision,
		i.Data.InstanceDescription(),
		contentType,
		data,
	)

	return err
}
