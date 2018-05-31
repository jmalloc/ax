package messages

import (
	"fmt"

	"github.com/jmalloc/ax/examples/banking/format"
	"github.com/jmalloc/ax/src/ax/ident"
)

// IsCommand marks the message as a command.
func (*StartTransfer) IsCommand() {}

// MessageDescription returns a human-readable description of the message.
func (m *StartTransfer) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s from %s to %s identified by %s",
		format.Amount(m.AmountInCents),
		ident.Format(m.FromAccountId),
		ident.Format(m.ToAccountId),
		ident.Format(m.TransferId),
	)
}

// IsEvent marks the message as an event.
func (*TransferStarted) IsEvent() {}

// MessageDescription returns a human-readable description of the message.
func (m *TransferStarted) MessageDescription() string {
	return fmt.Sprintf(
		"transfer of %s from %s to %s started as %s",
		format.Amount(m.AmountInCents),
		ident.Format(m.FromAccountId),
		ident.Format(m.ToAccountId),
		ident.Format(m.TransferId),
	)
}

// IsCommand marks the message as a command.
func (*MarkTransferApproved) IsCommand() {}

// MessageDescription returns a human-readable description of the message.
func (m *MarkTransferApproved) MessageDescription() string {
	return fmt.Sprintf(
		"mark transfer %s as approved",
		ident.Format(m.TransferId),
	)
}

// IsEvent marks the message as an event.
func (*TransferApproved) IsEvent() {}

// MessageDescription returns a human-readable description of the message.
func (m *TransferApproved) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s approved",
		ident.Format(m.TransferId),
	)
}
