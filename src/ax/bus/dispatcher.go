package bus

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
	"golang.org/x/sync/errgroup"
)

// Dispatcher is an inbound pipeline stage that delivers messages to
// MessageHandler implementations.
type Dispatcher struct {
	Routes DispatchTable
}

// Initialize configures the transport to subscribe to the events that the
// handlers can handle.
func (d *Dispatcher) Initialize(ctx context.Context, t Transport) error {
	var events ax.MessageTypeSet

	for mt := range d.Routes {
		if mt.IsEvent() {
			events = events.Add(mt)
		}
	}

	return t.Subscribe(ctx, events)
}

// DeliverMessage passes m to zero or more message handlers according to the
// dispatch table.
func (d *Dispatcher) DeliverMessage(
	ctx context.Context,
	s MessageSender,
	m InboundEnvelope,
) error {
	wg, ctx := errgroup.WithContext(ctx)
	mt := ax.TypeOf(m.Envelope.Message)
	mc := &MessageContext{
		Context:  ctx,
		Envelope: m.Envelope,
		Sender:   s,
	}

	for _, h := range d.Routes.Lookup(mt) {
		func(h MessageHandler) {
			wg.Go(func() error {
				return h.HandleMessage(mc, m.Envelope.Message)
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
