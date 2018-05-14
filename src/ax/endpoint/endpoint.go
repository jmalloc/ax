package endpoint

import (
	"context"
	"sync"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
)

// RetryPolicy returns true if the message should be retried.
type RetryPolicy func(bus.InboundEnvelope, error) bool

// DefaultRetryPolicy is a RetryPolicy that rejects a message after it has been
// attempted three (3) times.
func DefaultRetryPolicy(env bus.InboundEnvelope, _ error) bool {
	return env.DeliveryCount < 3
}

// Endpoint is a named source and recipient of messages.
type Endpoint struct {
	Name        string
	Transport   bus.Transport
	In          bus.InboundPipeline
	Out         bus.OutboundPipeline
	RetryPolicy RetryPolicy

	initOnce sync.Once
}

// NewSender returns an ax.Sender that can be used to send messages from this endpoint.
func (ep *Endpoint) NewSender(ctx context.Context) (ax.Sender, error) {
	if err := ep.initialize(ctx); err != nil {
		return nil, err
	}

	return bus.SinkSender{Sink: ep.Out}, nil
}

// StartReceiving processes inbound messages until an error occurrs or ctx is canceled.
func (ep *Endpoint) StartReceiving(ctx context.Context) error {
	if err := ep.initialize(ctx); err != nil {
		return err
	}

	recv := &receiver{
		Transport:   ep.Transport,
		In:          ep.In,
		Out:         ep.Out,
		RetryPolicy: ep.RetryPolicy,
	}

	return recv.Run(ctx)
}

func (ep *Endpoint) initialize(ctx context.Context) (err error) {
	ep.initOnce.Do(func() {
		err = ep.Transport.Initialize(ctx, ep.Name)
		if err != nil {
			return
		}

		err = ep.In.Initialize(ctx, ep.Transport)
		if err != nil {
			return
		}

		err = ep.Out.Initialize(ctx, ep.Transport)
		if err != nil {
			return
		}
	})

	return
}
