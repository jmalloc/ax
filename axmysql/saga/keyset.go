package saga

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmalloc/ax/axmysql/internal/sqlutil"
	mysqlpersistence "github.com/jmalloc/ax/axmysql/persistence"
	"github.com/jmalloc/ax/persistence"
	"github.com/jmalloc/ax/saga"
)

// KeySetRepository is a MySQL-backed implementation of Ax's keyset.Repository
// interface.
type KeySetRepository struct{}

// FindByKey returns the ID of a saga instance that has a specific key in
// its key set.
//
// pk is the saga's persistence key, mk is the mapping key.
// ok is false if no saga instance has a key set containing mk.
func (KeySetRepository) FindByKey(
	ctx context.Context,
	ptx persistence.Tx,
	pk, mk string,
) (id saga.InstanceID, ok bool, err error) {
	tx := mysqlpersistence.ExtractTx(ptx)

	err = tx.QueryRowContext(
		ctx,
		`SELECT
			instance_id
		FROM ax_saga_keyset
		WHERE persistence_key = ?
		AND mapping_key = ?`,
		pk,
		mk,
	).Scan(
		&id,
	)

	if err == nil {
		ok = true
	} else if err == sql.ErrNoRows {
		err = nil
	}

	return
}

// SaveKeys associates a set of mapping keys with a saga instance.
//
// Key sets must be disjoint. That is, no two instances of the same saga
// may share any keys.
//
// pk is the saga's persistence key. ks is the set of mapping keys.
//
// SaveKeys() may panic if ks contains duplicate keys.
func (r KeySetRepository) SaveKeys(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	ks []string,
	id saga.InstanceID,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	if err := r.deleteKeys(ctx, tx, pk, id); err != nil {
		return err
	}

	for _, mk := range ks {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO ax_saga_keyset SET
				persistence_key = ?,
				mapping_key = ?,
				instance_id = ?`,
			pk,
			mk,
			id,
		); err != nil {
			if sqlutil.IsDuplicateEntry(err) {
				return fmt.Errorf(
					"can not save mapping keys for instance %s, the '%s' key is mapped to another instance",
					id,
					mk,
				)
			}

			return err
		}
	}

	return nil
}

// DeleteKeys removes any mapping keys associated with a saga instance.
//
// pk is the saga's persistence key.
func (r KeySetRepository) DeleteKeys(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	id saga.InstanceID,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	return r.deleteKeys(ctx, tx, pk, id)
}

func (KeySetRepository) deleteKeys(
	ctx context.Context,
	tx *sql.Tx,
	pk string,
	id saga.InstanceID,
) error {
	_, err := tx.ExecContext(
		ctx,
		`DELETE FROM ax_saga_keyset
		WHERE persistence_key = ?
		AND instance_id = ?`,
		pk,
		id,
	)

	return err
}
