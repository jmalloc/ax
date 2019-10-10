package axrmq

import (
	"context"
	"time"

	"github.com/streadway/amqp"
)

// Acknowledger is an implementation of bus.Acknowledger that acknowledges AMQP messages.
type Acknowledger struct {
	ep  string
	pub *publisher
	con *consumer
	del amqp.Delivery
}

// Ack acknowledges the message, indicating that is was handled successfully
// and does not need to be retried.
func (a *Acknowledger) Ack(_ context.Context) error {
	return a.del.Ack(false) // false = single message
}

// Retry requeues the message so that it is retried at some point in the
// future.
//
// d is a hint as to how long the transport should wait before retrying
// this message.
func (a *Acknowledger) Retry(ctx context.Context, _ error, d time.Duration) error {
	if d >= 0 {
		if err := a.delay(ctx, d); err != nil {
			return err
		}
	}

	// Rejecting the message causes it to be requeued in the *same queue* via
	// the DLX configuration. This allows us to use the DLX x-death header to
	// get an attempt count.
	return a.del.Reject(false) // true = don't requeue
}

// Reject indicates that the message could not be handled and should not be
// retried. Depending on the transport, this may move the message to some form
// of error queue or otherwise drop the message completely.
func (a *Acknowledger) Reject(ctx context.Context, _ error) error {
	// When rejecting a message, we need to manually shovel it to the error queue
	// and then acknowledge the original message so that it is not requeued by the
	// DLX configuration. This may result on duplicate messages on the error queue.
	if err := a.pub.RepublishAsError(ctx, a.del); err != nil {
		return err
	}

	return a.del.Ack(false) // false = single message
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
