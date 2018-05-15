package routing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"golang.org/x/sync/errgroup"
)

// Dispatcher is an inbound pipeline stage that routes messages to the
// appropriate MessageHandler instances according to a "handler table".
type Dispatcher struct {
	Routes HandlerTable
}

// Initialize subscribes t to events that the message handlers intend to handle.
func (d *Dispatcher) Initialize(ctx context.Context, t endpoint.Transport) error {
	var unicast, multicast ax.MessageTypeSet

	for mt := range d.Routes {
		if mt.IsCommand() {
			unicast = unicast.Add(mt)
		} else if mt.IsEvent() {
			multicast = multicast.Add(mt)
		} else {
			unicast = unicast.Add(mt)
			multicast = multicast.Add(mt)
		}
	}

	if err := t.Subscribe(ctx, endpoint.OpSendUnicast, unicast); err != nil {
		return err
	}

	return t.Subscribe(ctx, endpoint.OpSendMulticast, multicast)
}

// Accept dispatches env to zero or more message handlers as per the dispatch
// table.
//
// The context passed to each handler contains the message envelope, such that
// any messages sent using s within that context are configured as children of env.
//
// Each message handler is invoked on its own goroutine.
func (d *Dispatcher) Accept(ctx context.Context, s endpoint.MessageSink, env endpoint.InboundEnvelope) error {
	ctx = endpoint.WithEnvelope(ctx, env.Envelope)
	wg, ctx := errgroup.WithContext(ctx)

	for _, h := range d.Routes.Lookup(env.Type()) {
		func(h MessageHandler) {
			wg.Go(func() error {
				return h.HandleMessage(
					ctx,
					endpoint.SinkSender{Sink: s},
					env.Envelope,
				)
			})
		}(h)
	}

	return wg.Wait()
}
