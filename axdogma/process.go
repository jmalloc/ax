package axdogma

import (
	"context"
	"sync"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/endpoint"
	"github.com/jmalloc/ax/persistence"
	"github.com/jmalloc/ax/saga"
)

// ProcessAdaptor adapts a dogma.ProcessMessageHandler to Ax's saga.Saga
// interface.
type ProcessAdaptor struct {
	saga.IgnoreNotFound

	Name       string
	EventTypes ax.MessageTypeSet
	Handler    dogma.ProcessMessageHandler

	m         sync.Mutex
	completed map[endpoint.AttemptID]struct{}
}

var _ saga.Saga = &ProcessAdaptor{}

// PersistenceKey returns a unique identifier for the saga.
//
// The persistence key is used to relate persisted data with the saga
// implementation that owns it. Persistence keys should not be changed once
// a saga has active instances.
func (a *ProcessAdaptor) PersistenceKey() string {
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
func (a *ProcessAdaptor) MessageTypes() (tr ax.MessageTypeSet, mt ax.MessageTypeSet) {
	return a.EventTypes, ax.MessageTypeSet{}
}

// NewData returns a pointer to a new zero-value instance of the
// saga's data type.
func (a *ProcessAdaptor) NewData() saga.Data {
	return a.Handler.New().(saga.Data)
}

// HandleMessage handles a message for a particular saga instance.
func (a *ProcessAdaptor) HandleMessage(
	ctx context.Context,
	s ax.Sender,
	mctx ax.MessageContext,
	i saga.Instance,
) (err error) {
	defer unwrap(&err)

	sc := &processScope{
		ctx:      ctx,
		sender:   s,
		mctx:     mctx,
		instance: i,
		exists:   i.Revision > 0,
	}

	err = a.Handler.HandleEvent(
		ctx,
		sc,
		mctx.Envelope.Message,
	)
	if err != nil {
		return
	}

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
func (a *ProcessAdaptor) IsInstanceComplete(ctx context.Context, i saga.Instance) (bool, error) {
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

type processScope struct {
	ctx      context.Context
	sender   ax.Sender
	mctx     ax.MessageContext
	instance saga.Instance
	exists   bool
}

func (s *processScope) InstanceID() string {
	return s.instance.InstanceID.String()
}

func (s *processScope) Begin() bool {
	if s.exists {
		return false
	}

	s.exists = true

	return true
}

func (s *processScope) End() {
	if !s.exists {
		panic("can not end non-existent instance")
	}

	s.exists = false
}

func (s *processScope) Root() dogma.ProcessRoot {
	if !s.exists {
		panic("can not access aggregate root of non-existent instance")
	}

	return s.instance.Data.(dogma.ProcessRoot)
}

func (s *processScope) ExecuteCommand(m dogma.Message) {
	if !s.exists {
		panic("can not execute command against non-existent instance")
	}

	_, err := s.sender.ExecuteCommand(s.ctx, m.(ax.Command))
	if err != nil {
		wrapAndPanic(err)
	}
}

func (s *processScope) ScheduleTimeout(m dogma.Message, t time.Time) {
	panic("<not implemented>")
	// if !s.exists {
	// 	panic("can not execute command against non-existent instance")
	// }

	// _, err := s.sender.ExecuteCommand(s.ctx, m.(ax.Command), ax.DelayUntil(t))
	// if err != nil {
	// 	wrapAndPanic(err)
	// }
}

func (s *processScope) Log(f string, v ...interface{}) {
	s.mctx.Log(f, v...)
}

// ProcessMapper is a saga.Mapper that maps messages to instances using a
// Dogma aggregate's RouteCommandToInstance() method.
type ProcessMapper struct{}

var _ saga.Mapper = &ProcessMapper{}

// MapMessageToInstance returns the ID of the saga instance that is the
// target of the given message.
//
// It returns false if the message should be ignored.
func (m *ProcessMapper) MapMessageToInstance(
	ctx context.Context,
	sg saga.Saga,
	_ persistence.Tx,
	env ax.Envelope,
) (saga.InstanceID, bool, error) {
	id, ok, err := sg.(*ProcessAdaptor).Handler.RouteEventToInstance(ctx, env.Message)
	if !ok || err != nil {
		return saga.InstanceID{}, false, err
	}

	return saga.MustParseInstanceID(id), true, nil
}

// UpdateMapping notifies the mapper that an instance has been modified,
// allowing it to update it's mapping information, if necessary.
func (m *ProcessMapper) UpdateMapping(
	context.Context,
	saga.Saga,
	persistence.Tx,
	saga.Instance,
) error {
	return nil
}

// DeleteMapping notifies the mapper that an instance has been completed,
// allowing it to remove it's mapping information, if necessary.
func (m *ProcessMapper) DeleteMapping(
	context.Context,
	saga.Saga,
	persistence.Tx,
	saga.Instance,
) error {
	return nil
}
