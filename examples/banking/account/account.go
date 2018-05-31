package account

import (
	"fmt"

	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax/aggregate"
	"github.com/jmalloc/ax/src/ax/ident"
)

func (a *Account) OpenAccount(m *messages.OpenAccount, rec aggregate.Recorder) {
	if !a.IsOpen {
		rec(&messages.AccountOpened{
			AccountId: m.AccountId,
			Name:      m.Name,
		})
	}
}

func (a *Account) CreditAccount(m *messages.CreditAccount, rec aggregate.Recorder) {
	rec(&messages.AccountCredited{
		AccountId: m.AccountId,
		Cents:     m.Cents,
	})
}

func (a *Account) DebitAccount(m *messages.DebitAccount, rec aggregate.Recorder) {
	rec(&messages.AccountDebited{
		AccountId: m.AccountId,
		Cents:     m.Cents,
	})
}

func (a *Account) WhenAccountOpened(m *messages.AccountOpened) {
	a.AccountId = m.AccountId
	a.Name = m.Name
	a.IsOpen = true
}

func (a *Account) WhenAccountCredited(m *messages.AccountCredited) {
	a.Balance += m.Cents
}

func (a *Account) WhenAccountDebited(m *messages.AccountDebited) {
	a.Balance -= m.Cents
}

// InstanceDescription returns a human-readable description of the aggregate
// instance.
func (a *Account) InstanceDescription() string {
	return fmt.Sprintf(
		"account %s for %s, balance of %d",
		ident.Format(a.AccountId),
		a.Name,
		a.Balance,
	)
}

// AggregateRoot is a saga that implements the Account aggregate.
var AggregateRoot = aggregate.New(
	&Account{},
	aggregate.IdentifyByField("AccountId"),
)
