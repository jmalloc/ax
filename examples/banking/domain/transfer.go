package domain

import (
	"fmt"

	"github.com/jmalloc/ax/examples/banking/format"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/ident"
)

// DoTransfer begins a new funds transfer between two accounts.
func (t *Transfer) DoTransfer(m *messages.StartTransfer, mctx ax.MessageContext, rec ax.EventRecorder) {
	if t.TransferId != "" {
		mctx.Log("transfer has already been started")
		return
	}

	rec(&messages.TransferStarted{
		TransferId:    m.TransferId,
		FromAccountId: m.FromAccountId,
		ToAccountId:   m.ToAccountId,
		AmountInCents: m.AmountInCents,
	})
}

// DoMarkComplete marks the transfer as completed.
func (t *Transfer) DoMarkComplete(m *messages.MarkTransferComplete, mctx ax.MessageContext, rec ax.EventRecorder) {
	if t.IsComplete {
		mctx.Log("transfer has already been completed")
		return
	}

	rec(&messages.TransferCompleted{
		TransferId: t.TransferId,
	})
}

// WhenStarted updates the transfer to reflect the occurance of m.
func (t *Transfer) WhenStarted(m *messages.TransferStarted) {
	t.TransferId = m.TransferId
	t.FromAccountId = m.FromAccountId
	t.ToAccountId = m.ToAccountId
	t.AmountInCents = m.AmountInCents
}

// WhenCompleted updates the transfer to reflect the occurance of m.
func (t *Transfer) WhenCompleted(m *messages.TransferCompleted) {
	t.IsComplete = true
}

// IsInstanceComplete returns true if the transfer has completed processing.
func (t *Transfer) IsInstanceComplete() bool {
	return t.IsComplete
}

// InstanceDescription returns a human-readable description of the aggregate
// instance.
func (t *Transfer) InstanceDescription() string {
	return fmt.Sprintf(
		"transfer %s of %s from %s to %s",
		ident.Format(t.TransferId),
		format.Amount(t.AmountInCents),
		ident.Format(t.FromAccountId),
		ident.Format(t.ToAccountId),
	)
}
