package messages

import (
	"fmt"

	"github.com/jmalloc/ax/examples/banking/format"
	"github.com/jmalloc/ax/src/ax/ident"
)

// IsCommand marks the message as a command.
func (*OpenAccount) IsCommand() {}

// MessageDescription returns a human-readable description of the message.
func (m *OpenAccount) MessageDescription() string {
	return fmt.Sprintf(
		"open account %s for %s",
		ident.Format(m.AccountId),
		m.Name,
	)
}

// IsEvent marks the message as an event.
func (*AccountOpened) IsEvent() {}

// MessageDescription returns a human-readable description of the message.
func (m *AccountOpened) MessageDescription() string {
	return fmt.Sprintf(
		"account %s opened for %s",
		ident.Format(m.AccountId),
		m.Name,
	)
}

// IsCommand marks the message as a command.
func (*CreditAccount) IsCommand() {}

// MessageDescription returns a human-readable description of the message.
func (m *CreditAccount) MessageDescription() string {
	return fmt.Sprintf(
		"credit %s to account %s",
		format.Amount(m.AmountInCents),
		ident.Format(m.AccountId),
	)
}

// IsEvent marks the message as an event.
func (*AccountCredited) IsEvent() {}

// MessageDescription returns a human-readable description of the message.
func (m *AccountCredited) MessageDescription() string {
	return fmt.Sprintf(
		"credited %s to account %s",
		format.Amount(m.AmountInCents),
		ident.Format(m.AccountId),
	)
}

// IsCommand marks the message as a command.
func (*DebitAccount) IsCommand() {}

// MessageDescription returns a human-readable description of the message.
func (m *DebitAccount) MessageDescription() string {
	return fmt.Sprintf(
		"debit %s from account %s",
		format.Amount(m.AmountInCents),
		ident.Format(m.AccountId),
	)
}

// IsEvent marks the message as an event.
func (*AccountDebited) IsEvent() {}

// MessageDescription returns a human-readable description of the message.
func (m *AccountDebited) MessageDescription() string {
	return fmt.Sprintf(
		"debited %s from account %s",
		format.Amount(m.AmountInCents),
		ident.Format(m.AccountId),
	)
}
