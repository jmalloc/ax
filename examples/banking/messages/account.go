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
