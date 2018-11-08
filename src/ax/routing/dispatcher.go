package routing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/observability"
	"github.com/jmalloc/twelf/src/twelf"
)

// Dispatcher is an inbound pipeline stage that routes messages to the
// appropriate MessageHandler instances according to a "handler table".
type Dispatcher struct {
	Routes HandlerTable
	Logger twelf.Logger

	validators []endpoint.Validator
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (d *Dispatcher) Initialize(ctx context.Context, ep *endpoint.Endpoint) error {
	d.validators = ep.SenderValidators

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

	if err := ep.InboundTransport.Subscribe(ctx, endpoint.OpSendUnicast, unicast); err != nil {
		return err
	}

	return ep.InboundTransport.Subscribe(ctx, endpoint.OpSendMulticast, multicast)
}

// Accept dispatches env to zero or more message handlers as per the dispatch
// table.
//
// The context passed to each handler contains the message envelope, such that
// any messages sent using s within that context are configured as children of env.
//
// Each message handler is invoked on its own goroutine.
func (d *Dispatcher) Accept(ctx context.Context, s endpoint.MessageSink, env endpoint.InboundEnvelope) error {
	sender := endpoint.SinkSender{
		Sink:       s,
		Validators: d.validators,
	}

	mctx := ax.NewMessageContext(
		env.Envelope,
		observability.NewDomainLogger(
			d.Logger,
			env.Envelope,
		),
	)

	for _, h := range d.Routes.Lookup(env.Type()) {
		if err := h.HandleMessage(ctx, sender, mctx); err != nil {
			return err
		}
	}

	return nil
}
