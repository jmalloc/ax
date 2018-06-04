package projections

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/axmysql"
)

// AccountProjector is a message handler that builds the "account" read-model.
var AccountProjector accountProjector

type accountProjector struct{}

func (accountProjector) MessageTypes() ax.MessageTypeSet {
	return ax.TypesOf(
		&messages.AccountOpened{},
		&messages.AccountDebited{},
		&messages.AccountCredited{},
	)
}

func (accountProjector) HandleMessage(ctx context.Context, _ ax.Sender, env ax.Envelope) error {
	tx := axmysql.GetTx(ctx)

	switch m := env.Message.(type) {
	case *messages.AccountOpened:
		return insertAccount(ctx, tx, m.AccountId, m.Name)
	case *messages.AccountDebited:
		return updateBalance(ctx, tx, m.AccountId, -m.AmountInCents)
	case *messages.AccountCredited:
		return updateBalance(ctx, tx, m.AccountId, +m.AmountInCents)
	}

	return nil
}

func insertAccount(
	ctx context.Context,
	tx *sql.Tx,
	id string,
	name string,
) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO account SET
			id = ?,
			name = ?
		ON DUPLICATE KEY UPDATE
			name = VALUE(name)`,
		id,
		name,
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
			balance = ?
		ON DUPLICATE KEY UPDATE
			balance = balance + VALUE(balance)`,
		id,
		delta,
	)

	return err
}
