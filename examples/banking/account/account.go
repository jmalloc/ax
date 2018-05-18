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
var AggregateRoot saga.Saga = &aggregateRoot{}

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

func (aggregateRoot) NewInstance(ctx context.Context, env ax.Envelope) (saga.InstanceID, saga.Data, error) {
	var id saga.InstanceID
	id.MustParse(env.Message.(*messages.OpenAccount).AccountId)

	return id, &Account{}, nil
}

func (aggregateRoot) BuildMappingTable(ctx context.Context, i saga.Instance) (map[string]string, error) {
	return map[string]string{
		"account_id": i.Data.(*Account).AccountId,
	}, nil
}

func (aggregateRoot) MapMessage(ctx context.Context, env ax.Envelope) (string, string, error) {
	type hasAccountID interface {
		GetAccountId() string
	}

	return "account_id", env.Message.(hasAccountID).GetAccountId(), nil
}

func (aggregateRoot) HandleMessage(
	ctx context.Context,
	s ax.Sender,
	env ax.Envelope,
	i saga.Instance,
) error {
	acct := i.Data.(*Account)

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
