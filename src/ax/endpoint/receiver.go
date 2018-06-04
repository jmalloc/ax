package endpoint

import (
	"context"

	"github.com/jmalloc/ax/src/internal/servicegroup"
)

// receiver receives a message from a transport, forwards it to the inbound
// pipeline, then acknowledges the message.
type receiver struct {
	Transport   Transport
	In          InboundPipeline
	Out         OutboundPipeline
	RetryPolicy RetryPolicy

	wg *servicegroup.Group
}

// Run processes inbound messages until ctx is canceled or an error occurrs.
func (r *receiver) Run(ctx context.Context) error {
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
func (r *receiver) receive(ctx context.Context) error {
	for {
		env, ack, err := r.Transport.Receive(ctx)
		if err != nil {
			return err
		}

		if err := r.wg.Go(func(ctx context.Context) error {
			return r.process(ctx, env, ack)
		}); err != nil {
			return err
		}
	}
}

func (r *receiver) process(
	ctx context.Context,
	env InboundEnvelope,
	ack Acknowledger,
) error {
	err := r.In.Accept(ctx, r.Out, env)

	if err == nil {
		return ack.Ack(ctx)
	}

	if d, ok := r.RetryPolicy(env, err); ok {
		return ack.Retry(ctx, err, d)
	}

	return ack.Reject(ctx, err)
}
