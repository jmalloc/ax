package endpoint

import (
	"context"

	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/internal/servicegroup"
)

// Receiver receives a message from a transport, forwards it to the inbound
// pipeline, then acknowledges the message.
type Receiver struct {
	Transport   bus.Transport
	In          bus.InboundPipeline
	Out         bus.OutboundPipeline
	RetryPolicy RetryPolicy

	wg *servicegroup.Group
}

// Run processes inbound messages until ctx is canceled or an error occurrs.
func (r *Receiver) Run(ctx context.Context) error {
	if r.RetryPolicy == nil {
		r.RetryPolicy = DefaultRetryPolicy
	}

	r.wg = servicegroup.NewGroup(ctx)

	if err := r.wg.Go(r.receive); err != nil {
		return err
	}

	return r.wg.Wait()
}

// receive starts a new goroutine to process each inbound message.
func (r *Receiver) receive(ctx context.Context) error {
	for {
		env, ack, err := r.Transport.Receive(ctx)
		if err != nil {
			return err
		}

		r.wg.Go(func(ctx context.Context) error {
			return r.process(ctx, env, ack)
		})
	}
}

func (r *Receiver) process(
	ctx context.Context,
	env bus.InboundEnvelope,
	ack bus.Acknowledger,
) error {
	err := r.In.Accept(ctx, r.Out, env)

	if err == nil {
		return ack.Ack(ctx)
	}

	if r.RetryPolicy(env, err) {
		return ack.Retry(ctx, err)
	}

	return ack.Reject(ctx, err)
}
