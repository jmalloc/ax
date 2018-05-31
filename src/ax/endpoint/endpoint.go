package endpoint

import (
	"context"
	"errors"
	"sync"

	"github.com/jmalloc/ax/src/ax"
)

// Endpoint is a named source and recipient of messages.
type Endpoint struct {
	Name        string
	Transport   Transport
	In          InboundPipeline
	Out         OutboundPipeline
	RetryPolicy RetryPolicy

	initOnce sync.Once
}

// NewSender returns an ax.Sender that can be used to send messages from this endpoint.
func (ep *Endpoint) NewSender(
	ctx context.Context,
	validators []Validator,
) (ax.Sender, error) {
	if err := ep.initialize(ctx); err != nil {
		return nil, err
	}

	return SinkSender{
		Sink:       ep.Out,
		Validators: validators,
	}, nil
}

// StartReceiving processes inbound messages until an error occurrs or ctx is canceled.
func (ep *Endpoint) StartReceiving(ctx context.Context) error {
	if ep.In == nil {
		return errors.New("can not receive on send-only endpoint, there is no inbound message pipeline")
	}

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

		if ep.In != nil {
			err = ep.In.Initialize(ctx, ep)
			if err != nil {
				return
			}
		}

		err = ep.Out.Initialize(ctx, ep)
		if err != nil {
			return
		}
	})

	return
}
