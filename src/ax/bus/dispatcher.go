package bus

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
	"golang.org/x/sync/errgroup"
)

// Dispatcher is an inbound pipeline stage that routes messages to the
// appropriate MessageHandler instances according to a "dispatch table".
type Dispatcher struct {
	Routes DispatchTable
}

// Initialize configures the subscribes t to events that the message handlers
// intend to handle.
func (d *Dispatcher) Initialize(ctx context.Context, t Transport) error {
	var events ax.MessageTypeSet

	for mt := range d.Routes {
		if mt.IsEvent() {
			events = events.Add(mt)
		}
	}

	return t.Subscribe(ctx, events)
}

// Accept dispatches env to zero or more message handlers as per the dispatch
// table.
//
// The context passed to each handler contains the message envelope, such that
// any messages sent using s within that context are configured as children of env.
//
// Each message handler is invoked on its own goroutine.
func (d *Dispatcher) Accept(ctx context.Context, s MessageSink, env InboundEnvelope) error {
	ctx = WithEnvelope(ctx, env.Envelope)
	wg, ctx := errgroup.WithContext(ctx)

	for _, h := range d.Routes.Lookup(env.Type()) {
		func(h MessageHandler) {
			wg.Go(func() error {
				return h.HandleMessage(
					ctx,
					SinkSender{s},
					env.Envelope,
				)
			})
		}(h)
	}

	return wg.Wait()
}

// DispatchTable is a set of rules that determines which handlers receive a
// message of a specific type.
type DispatchTable map[ax.MessageType][]MessageHandler

// NewDispatchTable returns a dispatch table that locates message handlers
// based on the message types that they handle.
func NewDispatchTable(handlers ...MessageHandler) (DispatchTable, error) {
	dt := DispatchTable{}

	for _, h := range handlers {
		for _, mt := range h.MessageTypes().Members() {
			x := dt[mt]

			if mt.IsCommand() && len(x) != 0 {
				return nil, fmt.Errorf(
					"can not build dispatch table, multiple message handlers are defined for the '%s' command",
					mt.Name,
				)
			}

			dt[mt] = append(x, h)
		}
	}

	return dt, nil
}

// Lookup returns the message handlers that handle mt.
func (dt DispatchTable) Lookup(mt ax.MessageType) []MessageHandler {
	return dt[mt]
}
