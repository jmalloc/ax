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
func (*MarkTransferComplete) IsCommand() {}

// MessageDescription returns a human-readable description of the message.
func (m *MarkTransferComplete) MessageDescription() string {
	return fmt.Sprintf(
		"mark transfer %s as complete",
		ident.Format(m.TransferId),
	)
}

// IsEvent marks the message as an event.
func (*TransferCompleted) IsEvent() {}

// MessageDescription returns a human-readable description of the message.
func (m *TransferCompleted) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s completed",
		ident.Format(m.TransferId),
	)
}
