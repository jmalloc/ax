package domain

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/ident"
	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/ax/saga/mapping/keyset"
)

// InstanceDescription returns a human-readable description of the aggregate
// instance.
func (w *TransferWorkflow) InstanceDescription() string {
	return fmt.Sprintf(
		"transfer workflow for %s",
		ident.Format(w.TransferId),
	)
}

var (
	// TransferWorkflowSaga is a saga that implements the business process for a
	// funds transfer.
	TransferWorkflowSaga saga.Saga = transferWorkflowSaga{}

	// TransferWorkflowResolver maps messages to the transfer workflow.
	TransferWorkflowResolver keyset.Resolver = transferWorkflowResolver{}
)

// IsInstanceComplete returns true if the transfer has completed processing.
func (w *TransferWorkflow) IsInstanceComplete() bool {
	return w.IsApproved
}

type transferWorkflowSaga struct {
	saga.ErrorIfNotFound
	saga.CompletableByData
}

func (transferWorkflowSaga) PersistenceKey() string {
	return proto.MessageName(&TransferWorkflow{})
}

func (transferWorkflowSaga) MessageTypes() (tr ax.MessageTypeSet, mt ax.MessageTypeSet) {
	return ax.TypesOf(
			&messages.TransferStarted{},
		), ax.TypesOf(
			&messages.AccountCredited{},
			&messages.AccountDebited{},
		)
}

func (transferWorkflowSaga) NewData() saga.Data {
	return &TransferWorkflow{}
}

func (transferWorkflowSaga) HandleMessage(ctx context.Context, s ax.Sender, env ax.Envelope, i saga.Instance) error {
	var err error
	wf := i.Data.(*TransferWorkflow)

	switch m := env.Message.(type) {
	case *messages.TransferStarted:
		wf.TransferId = m.TransferId
		wf.FromAccountId = m.FromAccountId
		wf.ToAccountId = m.ToAccountId
		wf.AmountInCents = m.AmountInCents

		_, err = s.ExecuteCommand(ctx, &messages.DebitAccount{
			AccountId:     wf.FromAccountId,
			AmountInCents: wf.AmountInCents,
			TransferId:    wf.TransferId,
		})

	case *messages.AccountDebited:
		_, err = s.ExecuteCommand(ctx, &messages.CreditAccount{
			AccountId:     wf.ToAccountId,
			AmountInCents: wf.AmountInCents,
			TransferId:    wf.TransferId,
		})

	case *messages.AccountCredited:
		wf.IsApproved = true

		_, err = s.ExecuteCommand(ctx, &messages.MarkTransferApproved{
			TransferId: wf.TransferId,
		})
	}

	return err
}

type transferWorkflowResolver struct{}

func (transferWorkflowResolver) GenerateInstanceID(ctx context.Context, env ax.Envelope) (id saga.InstanceID, err error) {
	id.GenerateUUID()
	return
}

func (transferWorkflowResolver) MappingKeyForMessage(ctx context.Context, env ax.Envelope) (k string, ok bool, err error) {
	type hasTransferID interface {
		GetTransferId() string
	}

	transferID := env.Message.(hasTransferID).GetTransferId()
	return transferID, transferID != "", nil
}

func (transferWorkflowResolver) MappingKeysForInstance(_ context.Context, i saga.Instance) ([]string, error) {
	return []string{
		i.Data.(*TransferWorkflow).TransferId, // map based on the transfer ID
	}, nil
}
