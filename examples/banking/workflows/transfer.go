package workflows

import (
	"fmt"

	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/ident"
)

// BeginWhenTransferStarted responds to m.
func (w *Transfer) BeginWhenTransferStarted(m *messages.TransferStarted, exec ax.CommandExecutor) {
	w.TransferId = m.TransferId
	w.FromAccountId = m.FromAccountId
	w.ToAccountId = m.ToAccountId
	w.AmountInCents = m.AmountInCents

	exec(&messages.DebitAccount{
		AccountId:     w.FromAccountId,
		AmountInCents: w.AmountInCents,
		TransferId:    w.TransferId,
	})
}

// WhenAccountDebited responds to m.
func (w *Transfer) WhenAccountDebited(m *messages.AccountDebited, exec ax.CommandExecutor) {
	exec(&messages.CreditAccount{
		AccountId:     w.ToAccountId,
		AmountInCents: w.AmountInCents,
		TransferId:    w.TransferId,
	})
}

// WhenAccountCredited responds to m.
func (w *Transfer) WhenAccountCredited(m *messages.AccountCredited, mctx ax.MessageContext, exec ax.CommandExecutor) {
	w.IsComplete = true

	exec(&messages.MarkTransferComplete{
		TransferId: w.TransferId,
	})

	mctx.Log("credit and debit have both completed successfully")
}

// IsInstanceComplete returns true if the transfer has completed processing.
func (w *Transfer) IsInstanceComplete() bool {
	return w.IsComplete
}

// InstanceDescription returns a human-readable description of the aggregate
// instance.
func (w *Transfer) InstanceDescription() string {
	return fmt.Sprintf(
		"transfer workflow for %s",
		ident.Format(w.TransferId),
	)
}
