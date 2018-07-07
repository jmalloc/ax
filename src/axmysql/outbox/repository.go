package outbox

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axmysql/internal/envelopestore"
	mysqlpersistence "github.com/jmalloc/ax/src/axmysql/persistence"
)

// Repository is a MySQL-backed implementation of Ax's outbox.Repository
// interface.
type Repository struct{}

// messageTable is the name of the SQL table that stores outbox messages.
const messageTable = "ax_outbox_message"

// LoadOutbox loads the unsent outbound messages that were produced when the
// message identified by id was first delivered.
func (Repository) LoadOutbox(
	ctx context.Context,
	ds persistence.DataStore,
	id ax.MessageID,
) ([]endpoint.OutboundEnvelope, bool, error) {
	db := mysqlpersistence.ExtractDB(ds)

	row := db.QueryRowContext(
		ctx,
		`SELECT EXISTS (
			SELECT * FROM ax_outbox WHERE causation_id = ?
		)`,
		id,
	)

	var ok bool
	if err := row.Scan(&ok); err != nil {
		return nil, false, err
	}

	if !ok {
		return nil, false, nil
	}

	rows, err := db.QueryContext(
		ctx,
		`SELECT `+envelopestore.Columns+`
		FROM `+messageTable+`
		WHERE causation_id = ?`,
		id,
	)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var envelopes []endpoint.OutboundEnvelope

	for rows.Next() {
		env, err := envelopestore.Scan(rows)
		if err != nil {
			return nil, false, err
		}

		envelopes = append(envelopes, env)
	}

	return envelopes, true, nil
}

// SaveOutbox saves a set of unsent outbound messages that were produced
// when the message identified by id was delivered.
func (Repository) SaveOutbox(
	ctx context.Context,
	ptx persistence.Tx,
	id ax.MessageID,
	envs []endpoint.OutboundEnvelope,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO ax_outbox SET causation_id = ?`,
		id,
	); err != nil {
		return err
	}

	for _, env := range envs {
		if err := envelopestore.Insert(ctx, tx, messageTable, env); err != nil {
			return err
		}
	}

	return nil
}

// MarkAsSent marks a message as sent, removing it from the outbox.
func (Repository) MarkAsSent(
	ctx context.Context,
	ptx persistence.Tx,
	env endpoint.OutboundEnvelope,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	return envelopestore.Delete(ctx, tx, messageTable, env)
}
