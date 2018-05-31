package domain

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/examples/banking/messages"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/ident"
	"github.com/jmalloc/ax/src/ax/saga"
)

// InstanceDescription returns a human-readable description of the aggregate
// instance.
func (w *TransferWorkflow) InstanceDescription() string {
	return fmt.Sprintf(
		"transfer workflow for %s",
		ident.Format(w.TransferId),
	)
}

type transferWorkflowSaga struct {
	saga.ErrorIfNotFound
}

func (transferWorkflowSaga) SagaName() string {
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

func (transferWorkflowSaga) GenerateInstanceID(ctx context.Context, env ax.Envelope) (id saga.InstanceID, err error) {
	id.GenerateUUID()
	return
}

func (transferWorkflowSaga) NewData() saga.Data {
	return &TransferWorkflow{}
}

func (transferWorkflowSaga) MappingKeyForMessage(ctx context.Context, env ax.Envelope) (k string, ok bool, err error) {
	type hasTransferID interface {
		GetTransferId() string
	}

	transferID := env.Message.(hasTransferID).GetTransferId()
	return transferID, transferID != "", nil
}

func (transferWorkflowSaga) MappingKeysForInstance(_ context.Context, i saga.Instance) ([]string, error) {
	return []string{
		i.Data.(*TransferWorkflow).TransferId, // map based on the transfer ID
	}, nil
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
		_, err = s.ExecuteCommand(ctx, &messages.MarkTransferApproved{
			TransferId: wf.TransferId,
		})

		// TODO: mark done, depends on https://github.com/jmalloc/ax/issues/16
	}

	return err
}

// TransferWorkflowSaga is a saga that implements the business process for a
// funds transfer.
var TransferWorkflowSaga saga.Saga = transferWorkflowSaga{}
