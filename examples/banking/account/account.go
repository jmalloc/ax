package account

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/ident"
	"github.com/jmalloc/ax/src/ax/saga"
)

// SagaDescription returns a human-readable description of the saga instance.
func (a *Account) SagaDescription() string {
	return fmt.Sprintf("account %s", ident.Format(a.AccountId))
}

// AggregateRoot is a saga that implements the Account aggregate.
var AggregateRoot = &aggregateRoot{}

type aggregateRoot struct {
	saga.ErrorIfNotFound
}

func (aggregateRoot) SagaName() string {
	return "Account"
}

func (aggregateRoot) MessageTypes() (ax.MessageTypeSet, ax.MessageTypeSet) {
	return ax.TypesOf(
			&messages.OpenAccount{},
		), ax.TypesOf(
			&messages.CreditAccount{},
			&messages.DebitAccount{},
		)
}

func (aggregateRoot) MapMessage(m ax.Message) string {
	type hasAccountID interface {
		GetAccountId() string
	}

	return m.(hasAccountID).GetAccountId()
}

func (aggregateRoot) MapData(_ ax.MessageType, i saga.Data) string {
	return i.(*Account).AccountId
}

func (aggregateRoot) NewInstance(m ax.Message) (saga.InstanceID, saga.Data) {
	var id saga.InstanceID
	id.MustParse(m.(*messages.OpenAccount).AccountId)
	return id, &Account{}
}

func (aggregateRoot) HandleMessage(
	ctx context.Context,
	s ax.Sender,
	env ax.Envelope,
	i saga.Data,
) error {
	acct := i.(*Account)

	switch m := env.Message.(type) {
	case *messages.OpenAccount:
		if acct.IsOpen {
			return nil
		}

		acct.IsOpen = true
		acct.AccountId = m.AccountId
		acct.Name = m.Name

		return s.PublishEvent(ctx, &messages.AccountOpened{
			AccountId: m.AccountId,
			Name:      m.Name,
		})

	case *messages.CreditAccount:
		acct.Balance += m.Cents

		return s.PublishEvent(ctx, &messages.AccountCredited{
			AccountId: m.AccountId,
			Cents:     m.Cents,
		})

	case *messages.DebitAccount:
		acct.Balance -= m.Cents

		return s.PublishEvent(ctx, &messages.AccountDebited{
			AccountId: m.AccountId,
			Cents:     m.Cents,
		})
	}

	return nil
}
