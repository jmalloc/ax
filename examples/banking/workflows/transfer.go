package workflows

import (
	"fmt"

	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/ident"
)

// StartWhenTransferStarted responds to m.
func (w *Transfer) StartWhenTransferStarted(m *messages.TransferStarted, x ax.CommandExecutor) {
	w.TransferId = m.TransferId
	w.FromAccountId = m.FromAccountId
	w.ToAccountId = m.ToAccountId
	w.AmountInCents = m.AmountInCents

	x(&messages.DebitAccount{
		AccountId:     w.FromAccountId,
		AmountInCents: w.AmountInCents,
		TransferId:    w.TransferId,
	})
}

// WhenAccountDebited responds to m.
func (w *Transfer) WhenAccountDebited(m *messages.AccountDebited, x ax.CommandExecutor) {
	x(&messages.CreditAccount{
		AccountId:     w.ToAccountId,
		AmountInCents: w.AmountInCents,
		TransferId:    w.TransferId,
	})
}

// WhenAccountCredited responds to m.
func (w *Transfer) WhenAccountCredited(m *messages.AccountCredited, x ax.CommandExecutor) {
	w.IsApproved = true

	x(&messages.MarkTransferApproved{
		TransferId: w.TransferId,
	})
}

// IsInstanceComplete returns true if the transfer has completed processing.
func (w *Transfer) IsInstanceComplete() bool {
	return w.IsApproved
}

// InstanceDescription returns a human-readable description of the aggregate
// instance.
func (w *Transfer) InstanceDescription() string {
	return fmt.Sprintf(
		"transfer workflow for %s",
		ident.Format(w.TransferId),
	)
}
