package axrmq

import (
	"context"

	"github.com/streadway/amqp"
)

// Acknowledger is an implementation of bus.Acknowledger that acknowledges AMQP messages.
type Acknowledger struct {
	Delivery amqp.Delivery
}

// Ack acknowledges the message, indicating that is was handled successfully
// and does not need to be redelivered.
func (a *Acknowledger) Ack(_ context.Context) error {
	return a.Delivery.Ack(false) // false = single message
}

// Retry requeues the message so that it is redelivered at some point in the
// future.
func (a *Acknowledger) Retry(_ context.Context, _ error) error {
	return a.Delivery.Reject(true) // true = requeue
}

// Reject indicates that the message could not be handled and should not be
// retried. Depending on the transport, this may move the message to some form
// of error queue or otherwise drop the message completely.
func (a *Acknowledger) Reject(_ context.Context, _ error) error {
	return a.Delivery.Reject(false) // false = don't requeue
}
