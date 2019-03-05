package axdogma

import (
	"context"
	"sync"

	"github.com/dogmatiq/dogma"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/saga"
)

// AggregateAdaptor adapts a dogma.AggregateMessageHandler to Ax's saga.Saga
// interface.
type AggregateAdaptor struct {
	saga.IgnoreNotFound

	Name         string
	CommandTypes ax.MessageTypeSet
	Handler      dogma.AggregateMessageHandler

	m         sync.Mutex
	completed map[endpoint.AttemptID]struct{}
}

var _ saga.Saga = &AggregateAdaptor{}

// PersistenceKey returns a unique identifier for the saga.
//
// The persistence key is used to relate persisted data with the saga
// implementation that owns it. Persistence keys should not be changed once
// a saga has active instances.
func (a *AggregateAdaptor) PersistenceKey() string {
	return a.Name
}

// MessageTypes returns the set of messages that are routed to this saga.
//
// tr is the set of "trigger" messages. If they can not be routed to an
// existing saga instance a new instance is created.
//
// mt is the set of messages that are only routed to existing instances. If
// they can not be routed to an existing instance, the HandleNotFound()
// method is called instead.
func (a *AggregateAdaptor) MessageTypes() (tr ax.MessageTypeSet, mt ax.MessageTypeSet) {
	return a.CommandTypes, ax.MessageTypeSet{}
}

// NewData returns a pointer to a new zero-value instance of the
// saga's data type.
func (a *AggregateAdaptor) NewData() saga.Data {
	return a.Handler.New().(saga.Data)
}

// HandleMessage handles a message for a particular saga instance.
func (a *AggregateAdaptor) HandleMessage(
	ctx context.Context,
	s ax.Sender,
	mctx ax.MessageContext,
	i saga.Instance,
) (err error) {
	defer unwrap(&err)

	sc := &aggregateCommandScope{
		ctx:      ctx,
		sender:   s,
		mctx:     mctx,
		instance: i,
		exists:   i.Revision > 0,
	}

	a.Handler.HandleCommand(
		sc,
		mctx.Envelope.Message,
	)

	if !sc.exists {
		env, ok := endpoint.GetEnvelope(ctx)
		if !ok {
			panic("context does not contain an inbound envelope")
		}

		a.m.Lock()
		defer a.m.Unlock()

		if a.completed == nil {
			a.completed = map[endpoint.AttemptID]struct{}{}
		}

		a.completed[env.AttemptID] = struct{}{}
	}

	return
}

// IsInstanceComplete returns true if the given instance is complete.
func (a *AggregateAdaptor) IsInstanceComplete(ctx context.Context, i saga.Instance) (bool, error) {
	env, ok := endpoint.GetEnvelope(ctx)
	if !ok {
		panic("context does not contain an inbound envelope")
	}

	a.m.Lock()
	defer a.m.Unlock()

	_, ok = a.completed[env.AttemptID]
	delete(a.completed, env.AttemptID)

	return ok, nil
}

// ApplyEvent updates d to reflect the fact that an event has occurred.
//
// It may panic if env.Message does not implement ax.Event.
func (a *AggregateAdaptor) ApplyEvent(d saga.Data, env ax.Envelope) {
	d.(dogma.AggregateRoot).ApplyEvent(env.Message)
}

type aggregateCommandScope struct {
	ctx       context.Context
	sender    ax.Sender
	mctx      ax.MessageContext
	instance  saga.Instance
	exists    bool
	created   bool
	destroyed bool
}

func (s *aggregateCommandScope) InstanceID() string {
	return s.instance.InstanceID.String()
}

func (s *aggregateCommandScope) Create() bool {
	if s.exists {
		return false
	}

	s.exists = true
	s.created = true

	return true
}

func (s *aggregateCommandScope) Destroy() {
	if !s.exists {
		panic("can not destroy non-existent instance")
	}

	s.exists = false
	s.destroyed = true
}

func (s *aggregateCommandScope) Root() dogma.AggregateRoot {
	if !s.exists {
		panic("can not access aggregate root of non-existent instance")
	}

	return s.instance.Data.(dogma.AggregateRoot)
}

func (s *aggregateCommandScope) RecordEvent(m dogma.Message) {
	if !s.exists {
		panic("can not record event against non-existent instance")
	}

	_, err := s.sender.PublishEvent(s.ctx, m.(ax.Event))
	if err != nil {
		wrapAndPanic(err)
	}
}

func (s *aggregateCommandScope) Log(f string, v ...interface{}) {
	s.mctx.Log(f, v...)
}
