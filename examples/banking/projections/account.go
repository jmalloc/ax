package projections

import (
	"context"
	"database/sql"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmalloc/ax/src/ax"

	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax/projection"
	"github.com/jmalloc/ax/src/axmysql"
)

type account struct{}

func (account) PersistenceKey() string {
	return "Account"
}

func (account) WhenAccountOpened(ctx context.Context, tx *sql.Tx, ev *messages.AccountOpened, env ax.Envelope) error {
	return insertAccount(ctx, tx, ev.AccountId, ev.Name, env.Time)
}

func (account) WhenAccountDebited(ctx context.Context, tx *sql.Tx, ev *messages.AccountDebited) error {
	return updateBalance(ctx, tx, ev.AccountId, -ev.AmountInCents)
}

func (account) WhenAccountCredited(ctx context.Context, tx *sql.Tx, ev *messages.AccountCredited) error {
	return updateBalance(ctx, tx, ev.AccountId, +ev.AmountInCents)
}

func insertAccount(
	ctx context.Context,
	tx *sql.Tx,
	id string,
	name string,
	at time.Time,
) error {
	spew.Dump(at)

	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO account SET
			id = ?,
			name = ?,
			opened_at = ?
		ON DUPLICATE KEY UPDATE
			name = VALUE(name),
			opened_at = VALUE(opened_at)`,
		id,
		name,
		at,
	)

	return err
}

func updateBalance(
	ctx context.Context,
	tx *sql.Tx,
	id string,
	delta int32,
) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO account SET
			id = ?,
			balance = ? / 100
		ON DUPLICATE KEY UPDATE
			balance = balance + VALUE(balance)`,
		id,
		delta,
	)

	return err
}

// AccountProjector is a message handler that builds the "account" read-model.
var AccountProjector projection.Projector = axmysql.NewReadModelProjector(
	account{},
)
