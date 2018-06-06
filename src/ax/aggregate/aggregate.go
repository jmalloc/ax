package aggregate

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/ident"
	"github.com/jmalloc/ax/src/ax/saga"
)

// ID uniquely identifies an aggregate instance.
type ID struct{ ident.ID }

// Aggregate is an interface for application-defined domain aggregates.
//
// Aggregates are a specialization of sagas (stateful message handlers) that
// handle commands and produce events.
//
// For each command type to be handled, the aggregate must implement a "handler"
// method that adheres to the following signature:
//
//     func (cmd *<T>, rec Recorder)
//
// Where T is a struct type that implements ax.Command.
//
// Handler methods are responsible for producing new events based on the cmd
// being handled. They may inspect the current state of the aggregate, and then
// record zero or more events using rec. Handlers should never mutate the
// aggregate state.
//
// For each of the event types passed to rec, the aggregate must implement an
// "applier" method that adheres to the following signature:
//
//     func (ev *T)
//
// Where T is a struct type that implements ax.Event.
//
// Applier methods are responsible for mutating the aggregate state. The applier
// is called every time an event is recorded, *and* when loading an event-sourced
// aggregate from the message store.
//
// The names of handler and applier methods are not meaningful to the aggregate
// system. By convention, command handlers are named after their commands, and
// event appliers are prefixed with the word "When", such as:
//
//     func (*BankAccount) CreditAccount(*messages.CreditAccount, Recorder)
//     func (*BankAccount) WhenAccountCredited(*messages.AccountCredited)
type Aggregate interface {
	saga.Data
}

// Recorder is a function that records the events produced by an aggregate.
type Recorder func(ax.Event)
