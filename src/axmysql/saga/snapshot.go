package saga

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
	mysqlpersistence "github.com/jmalloc/ax/src/axmysql/persistence"
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
	var (
		ipk         string
		i           saga.Instance
		contentType string
		data        []byte
	)

	err := mysqlpersistence.ExtractTx(ptx).QueryRowContext(
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
		&ipk,
		&contentType,
		&data,
	)

	if err == sql.ErrNoRows {
		return saga.Instance{}, false, nil
	}

	if err != nil {
		return saga.Instance{}, false, nil
	}

	if ipk != pk {
		return i, false, fmt.Errorf(
			"can not load saga snapshot of %s at revision %d for saga %s, it belongs to %s",
			i.InstanceID,
			i.Revision,
			pk,
			ipk,
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
	tx persistence.Tx,
	pk string,
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
			persistence_key = ?,
			description = ?,
			content_type = ?,
			data = ?`,
		i.InstanceID,
		i.Revision,
		pk,
		i.Data.InstanceDescription(),
		contentType,
		data,
	)

	return err
}
