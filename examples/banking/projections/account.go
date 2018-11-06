package projections

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/projection"
	"github.com/jmalloc/ax/src/axmysql"
)

type account struct{}

func (account) PersistenceKey() string {
	return "Account"
}

func (account) WhenAccountOpened(ctx context.Context, tx *sql.Tx, ev *messages.AccountOpened, mctx ax.MessageContext) error {
	new, err := insertAccount(ctx, tx, ev.AccountId, ev.Name, mctx.Envelope.CreatedAt)
	if err != nil {
		return err
	}

	if new {
		mctx.Log("account now reflected in read model")
	}

	return nil
}

func (account) WhenAccountDebited(ctx context.Context, tx *sql.Tx, ev *messages.AccountDebited, mctx ax.MessageContext) error {
	new, err := updateBalance(ctx, tx, ev.AccountId, -ev.AmountInCents)
	if err != nil {
		return err
	}

	if new {
		mctx.Log("account now reflected in read model")
	}

	return nil
}

func (account) WhenAccountCredited(ctx context.Context, tx *sql.Tx, ev *messages.AccountCredited, mctx ax.MessageContext) error {
	new, err := updateBalance(ctx, tx, ev.AccountId, +ev.AmountInCents)
	if err != nil {
		return err
	}

	if new {
		mctx.Log("account now reflected in read model")
	}

	return nil
}

func insertAccount(
	ctx context.Context,
	tx *sql.Tx,
	id string,
	name string,
	at time.Time,
) (bool, error) {
	r, err := tx.ExecContext(
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
	if err != nil {
		return false, err
	}

	n, err := r.RowsAffected()
	if err != nil {
		return false, err
	}

	return n == 1, nil
}

func updateBalance(
	ctx context.Context,
	tx *sql.Tx,
	id string,
	delta int32,
) (bool, error) {
	r, err := tx.ExecContext(
		ctx,
		`INSERT INTO account SET
			id = ?,
			balance = ? / 100
		ON DUPLICATE KEY UPDATE
			balance = balance + VALUE(balance)`,
		id,
		delta,
	)

	if err != nil {
		return false, err
	}

	n, err := r.RowsAffected()
	if err != nil {
		return false, err
	}

	return n == 1, nil
}

// AccountProjector is a message handler that builds the "account" read-model.
var AccountProjector projection.Projector = axmysql.NewReadModelProjector(
	account{},
)
