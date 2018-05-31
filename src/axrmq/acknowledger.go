package axrmq

import (
	"context"
	"time"

	"github.com/streadway/amqp"
)

// Acknowledger is an implementation of bus.Acknowledger that acknowledges AMQP messages.
type Acknowledger struct {
	con *consumer
	del amqp.Delivery
}

// Ack acknowledges the message, indicating that is was handled successfully
// and does not need to be redelivered.
func (a *Acknowledger) Ack(_ context.Context) error {
	return a.del.Ack(false) // false = single message
}

// Retry requeues the message so that it is redelivered at some point in the
// future.
//
// d is a hint as to how long the transport should wait before redelivering
// this message.
func (a *Acknowledger) Retry(ctx context.Context, _ error, d time.Duration) error {
	if d >= 0 {
		if err := a.delay(ctx, d); err != nil {
			return err
		}
	}

	return a.del.Reject(true) // true = requeue
}

// Reject indicates that the message could not be handled and should not be
// retried. Depending on the transport, this may move the message to some form
// of error queue or otherwise drop the message completely.
func (a *Acknowledger) Reject(_ context.Context, _ error) error {
	return a.del.Reject(false) // false = don't requeue
}

func (a *Acknowledger) delay(ctx context.Context, d time.Duration) error {
	if err := a.con.IncreasePreFetch(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
	case <-time.After(d):
	}

	if err := a.con.DecreasePreFetch(); err != nil {
		return err
	}

	return ctx.Err()
}
