package account

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/ident"
	"github.com/jmalloc/ax/src/ax/saga"
)

// InstanceDescription returns a human-readable description of the saga instance.
func (a *Account) InstanceDescription() string {
	return fmt.Sprintf(
		"account %s for %s, balance of %d",
		ident.Format(a.AccountId),
		a.Name,
		a.Balance,
	)
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

func (aggregateRoot) GenerateInstanceID(ctx context.Context, env ax.Envelope) (id saga.InstanceID, err error) {
	err = id.Parse(env.Message.(*messages.OpenAccount).AccountId)
	return
}

func (aggregateRoot) NewData() saga.Data {
	return &Account{}
}

func (aggregateRoot) MappingKeyForMessage(ctx context.Context, env ax.Envelope) (string, error) {
	type hasAccountID interface {
		GetAccountId() string
	}

	return env.Message.(hasAccountID).GetAccountId(), nil
}

func (aggregateRoot) MappingKeysForInstance(
	ctx context.Context,
	i saga.Instance,
) (saga.KeySet, error) {
	return saga.NewKeySet(
		i.Data.(*Account).AccountId,
	), nil
}

func (aggregateRoot) HandleMessage(
	ctx context.Context,
	s ax.Sender,
	env ax.Envelope,
	i saga.Instance,
) (err error) {
	acct := i.Data.(*Account)

	switch m := env.Message.(type) {
	case *messages.OpenAccount:
		if acct.IsOpen {
			return
		}

		_, err = s.PublishEvent(ctx, &messages.AccountOpened{
			AccountId: m.AccountId,
			Name:      m.Name,
		})

	case *messages.CreditAccount:
		_, err = s.PublishEvent(ctx, &messages.AccountCredited{
			AccountId: m.AccountId,
			Cents:     m.Cents,
		})

	case *messages.DebitAccount:
		_, err = s.PublishEvent(ctx, &messages.AccountDebited{
			AccountId: m.AccountId,
			Cents:     m.Cents,
		})
	}

	return
}

// ApplyEvent updates the data to reflect the fact that ev has occurred.
func (aggregateRoot) ApplyEvent(d saga.Data, env ax.Envelope) {
	acct := d.(*Account)

	switch ev := env.Message.(type) {
	case *messages.AccountOpened:
		acct.AccountId = ev.AccountId
		acct.Name = ev.Name
		acct.IsOpen = true
	case *messages.AccountCredited:
		acct.Balance += ev.Cents
	case *messages.AccountDebited:
		acct.Balance -= ev.Cents
	}
}
