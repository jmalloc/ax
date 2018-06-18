package domain

import (
	"fmt"

	"github.com/jmalloc/ax/examples/banking/format"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax/aggregate"
	"github.com/jmalloc/ax/src/ax/ident"
)

// Start begins a new funds transfer between two accounts.
func (t *Transfer) Start(m *messages.StartTransfer, rec aggregate.Recorder) {
	if t.TransferId != "" {
		return
	}

	rec(&messages.TransferStarted{
		TransferId:    m.TransferId,
		FromAccountId: m.FromAccountId,
		ToAccountId:   m.ToAccountId,
		AmountInCents: m.AmountInCents,
	})
}

// MarkApproved marks the transfer as approved.
func (t *Transfer) MarkApproved(m *messages.MarkTransferApproved, rec aggregate.Recorder) {
	if t.IsApproved {
		return
	}

	rec(&messages.TransferApproved{
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

// WhenApproved updates the transfer to reflect the occurance of m.
func (t *Transfer) WhenApproved(m *messages.TransferApproved) {
	t.IsApproved = true
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
