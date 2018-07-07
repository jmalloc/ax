package delayedmessage

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axmysql/internal/envelopestore"
	"github.com/jmalloc/ax/src/axmysql/internal/sqlutil"
	mysqlpersistence "github.com/jmalloc/ax/src/axmysql/persistence"
)

// Repository is a MySQL-backed implementation of Ax's delayedmessage.Repository
// interface.
type Repository struct{}

// messageTable is the name of the SQL table that stores delayed messages.
const messageTable = "ax_delayed_message"

// LoadNextMessage loads the next that is scheduled to be sent.
func (Repository) LoadNextMessage(
	ctx context.Context,
	ds persistence.DataStore,
) (endpoint.OutboundEnvelope, bool, error) {
	db := mysqlpersistence.ExtractDB(ds)

	row := db.QueryRowContext(
		ctx,
		`SELECT `+envelopestore.Columns+`
		FROM `+messageTable+`
		ORDER BY send_at
		LIMIT 1`,
	)

	env, err := envelopestore.Scan(row)
	if err == sql.ErrNoRows {
		return endpoint.OutboundEnvelope{}, false, nil
	} else if err != nil {
		return endpoint.OutboundEnvelope{}, false, err
	}

	return env, true, nil
}

// SaveMessage saves a message to be sent at a later time.
// If does NOT return an error if the message already exists in the repository.
func (Repository) SaveMessage(
	ctx context.Context,
	ptx persistence.Tx,
	env endpoint.OutboundEnvelope,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	err := envelopestore.Insert(ctx, tx, messageTable, env)

	if sqlutil.IsDuplicateEntry(err) {
		return nil
	}

	return err
}

// MarkAsSent marks a message as sent, removing it from the repository.
func (Repository) MarkAsSent(
	ctx context.Context,
	ptx persistence.Tx,
	env endpoint.OutboundEnvelope,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	return envelopestore.Delete(ctx, tx, messageTable, env)
}
