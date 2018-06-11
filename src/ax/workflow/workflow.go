package workflow

import (
	"github.com/jmalloc/ax/src/ax/ident"
	"github.com/jmalloc/ax/src/ax/saga"
)

// ID uniquely identifies a workflow instance.
type ID struct{ ident.ID }

// Workflow is an interface for application-defined process managers.
//
// Workflows are a specialization of sagas (stateful message handlers) that
// handle events and produce commands.
//
// For each event type to be handled, the aggregate must implement a "handler"
// method that adheres to the following signature:
//
//     func (ev *<T>, do Executor)
//
// Where T is a struct type that implements ax.Event.
//
// Handler methods are responsible for mutating the state of the workflow and
// producing new commands, based on the event being handled. They may inspect
// the current state of the workflow, and then return zero or more commands
// to be executed.
//
// The names of handler methods are meaningful to the workflow system. If an
// event is meant to trigger a new workflow instance, its handler method's name
// must begin with "StartWhen". Other handler methods must begin with "When". By
// convention these prefixes are followed by the message name, such as:
//
//     func (*BankTransferWorkflow) StartWhenTransferStarted(*messages.TransferStarted)
//     func (*BankTransferWorkflow) WhenAccountDebited(*messages.AccountDebited)
type Workflow interface {
	saga.Data
}
