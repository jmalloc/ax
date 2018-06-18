package domain

import (
	"fmt"

	"github.com/jmalloc/ax/examples/banking/format"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/ident"
)

// OpenAccount opens a new account.
func (a *Account) OpenAccount(m *messages.OpenAccount, rec ax.EventRecorder) {
	if !a.IsOpen {
		rec(&messages.AccountOpened{
			AccountId: m.AccountId,
			Name:      m.Name,
		})
	}
}

// CreditAccount credits funds to the account.
func (a *Account) CreditAccount(m *messages.CreditAccount, rec ax.EventRecorder) {
	rec(&messages.AccountCredited{
		AccountId:     m.AccountId,
		AmountInCents: m.AmountInCents,
		TransferId:    m.TransferId,
	})
}

// DebitAccount debits funds from the account.
func (a *Account) DebitAccount(m *messages.DebitAccount, rec ax.EventRecorder) {
	rec(&messages.AccountDebited{
		AccountId:     m.AccountId,
		AmountInCents: m.AmountInCents,
		TransferId:    m.TransferId,
	})
}

// WhenAccountOpened updates the account to reflect the occurance of m.
func (a *Account) WhenAccountOpened(m *messages.AccountOpened) {
	a.AccountId = m.AccountId
	a.Name = m.Name
	a.IsOpen = true
}

// WhenAccountCredited updates the account to reflect the occurance of m.
func (a *Account) WhenAccountCredited(m *messages.AccountCredited) {
	a.BalanceInCents += m.AmountInCents
}

// WhenAccountDebited updates the account to reflect the occurance of m.
func (a *Account) WhenAccountDebited(m *messages.AccountDebited) {
	a.BalanceInCents -= m.AmountInCents
}

// InstanceDescription returns a human-readable description of the aggregate
// instance.
func (a *Account) InstanceDescription() string {
	return fmt.Sprintf(
		"account %s for %s with balance of %s",
		ident.Format(a.AccountId),
		a.Name,
		format.Amount(a.BalanceInCents),
	)
}
