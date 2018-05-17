package messages

import (
	"fmt"

	"github.com/jmalloc/ax/src/ax/ident"
)

// IsCommand marks the message as a command.
func (*OpenAccount) IsCommand() {}

// Description returns a human-readable description of the message.
func (m *OpenAccount) Description() string {
	return fmt.Sprintf(
		"open account %s for %s",
		ident.Format(m.AccountId),
		m.Name,
	)
}

// IsEvent marks the message as an event.
func (*AccountOpened) IsEvent() {}

// Description returns a human-readable description of the message.
func (m *AccountOpened) Description() string {
	return fmt.Sprintf(
		"account %s opened for %s",
		ident.Format(m.AccountId),
		m.Name,
	)
}

// IsCommand marks the message as a command.
func (*CreditAccount) IsCommand() {}

// Description returns a human-readable description of the message.
func (m *CreditAccount) Description() string {
	return fmt.Sprintf(
		"credit %d to account %s",
		m.Cents,
		ident.Format(m.AccountId),
	)
}

// IsEvent marks the message as an event.
func (*AccountCredited) IsEvent() {}

// Description returns a human-readable description of the message.
func (m *AccountCredited) Description() string {
	return fmt.Sprintf(
		"credited %d to account %s",
		m.Cents,
		ident.Format(m.AccountId),
	)
}

// IsCommand marks the message as a command.
func (*DebitAccount) IsCommand() {}

// Description returns a human-readable description of the message.
func (m *DebitAccount) Description() string {
	return fmt.Sprintf(
		"debit %d from account %s",
		m.Cents,
		ident.Format(m.AccountId),
	)
}

// IsEvent marks the message as an event.
func (*AccountDebited) IsEvent() {}

// Description returns a human-readable description of the message.
func (m *AccountDebited) Description() string {
	return fmt.Sprintf(
		"debited %s from account %s",
		m.Cents,
		ident.Format(m.AccountId),
	)
}
