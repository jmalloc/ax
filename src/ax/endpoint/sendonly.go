package endpoint

import (
	"context"
	"sync"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
)

// SendOnlyEndpoint is an endpoint that can produce messages, but does not receive
// messages.
type SendOnlyEndpoint struct {
	Name      string
	Transport bus.Transport
	Out       bus.OutboundPipeline

	initOnce sync.Once
}

// NewSender returns an ax.Sender that can be used to send messages from this endpoint.
func (ep *SendOnlyEndpoint) NewSender(ctx context.Context) (ax.Sender, error) {
	if err := ep.initialize(ctx); err != nil {
		return nil, err
	}

	return bus.SinkSender{Sink: ep.Out}, nil
}

func (ep *SendOnlyEndpoint) initialize(ctx context.Context) (err error) {
	ep.initOnce.Do(func() {
		err = ep.Transport.Initialize(ctx, ep.Name)
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
