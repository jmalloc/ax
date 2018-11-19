package endpoint

import (
	"context"
	"errors"
	"sync"

	"github.com/jmalloc/ax/src/ax"
	opentracing "github.com/opentracing/opentracing-go"
)

// Endpoint is a named source and recipient of messages.
type Endpoint struct {
	Name              string
	OutboundTransport OutboundTransport
	InboundTransport  InboundTransport
	InboundPipeline   InboundPipeline
	OutboundPipeline  OutboundPipeline
	RetryPolicy       RetryPolicy
	SenderValidators  []Validator
	Tracer            opentracing.Tracer

	initOnce sync.Once
}

// NewSender returns an ax.Sender that can be used to send messages from this endpoint.
func (ep *Endpoint) NewSender(ctx context.Context) (ax.Sender, error) {
	if err := ep.initialize(ctx); err != nil {
		return nil, err
	}

	return SinkSender{
		Sink:       ep.OutboundPipeline,
		Validators: ep.SenderValidators,
	}, nil
}

// StartReceiving processes inbound messages until an error occurrs or ctx is canceled.
func (ep *Endpoint) StartReceiving(ctx context.Context) error {
	if ep.InboundTransport == nil {
		return errors.New("can not receive on send-only endpoint, there is no inbound transport")
	}

	if ep.InboundPipeline == nil {
		return errors.New("can not receive on send-only endpoint, there is no inbound message pipeline")
	}

	if err := ep.initialize(ctx); err != nil {
		return err
	}

	recv := &receiver{
		Transport:   ep.InboundTransport,
		In:          ep.InboundPipeline,
		Out:         ep.OutboundPipeline,
		RetryPolicy: ep.RetryPolicy,
		Tracer:      ep.Tracer,
	}

	return recv.Run(ctx)
}

func (ep *Endpoint) initialize(ctx context.Context) (err error) {
	ep.initOnce.Do(func() {
		err = ep.OutboundTransport.Initialize(ctx, ep.Name)
		if err != nil {
			return
		}

		err = ep.OutboundPipeline.Initialize(ctx, ep)
		if err != nil {
			return
		}

		if ep.InboundTransport != nil && ep.InboundPipeline != nil {
			err = ep.InboundTransport.Initialize(ctx, ep.Name)
			if err != nil {
				return
			}

			err = ep.InboundPipeline.Initialize(ctx, ep)
			if err != nil {
				return
			}
		}
	})

	return
}
