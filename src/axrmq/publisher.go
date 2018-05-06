package axrmq

import (
	"context"
	"errors"

	"github.com/streadway/amqp"
)

// Publisher publishes AMQP messages reliably using AMQP "publisher confirms".
// https://www.rabbitmq.com/confirms.html#publisher-confirms
//
// It maintains a capped-size pool of AMQP channels which are placed into
// "confirm mode" when they are first created.
type Publisher struct {
	conn     *amqp.Connection
	channels chan *channel
}

// NewPublisher returns a new publisher with a channel pool
func NewPublisher(conn *amqp.Connection, n int) *Publisher {
	return &Publisher{
		conn,
		make(chan *channel, n),
	}
}

// channel is a handle to an AMQP channel that has been placed into "confirm
// mode", along with (Go) channels used to signal when confirmations have been
// received from the broker.
type channel struct {
	Channel *amqp.Channel
	Close   chan *amqp.Error
	Return  chan amqp.Return
	Confirm chan amqp.Confirmation
}

// Publish sends a message to the broker, and blocks until a confirmation is
// received.
//
// It returns an error if the broker does not acknowledge publication of the
// message. Otherwise it has the same behavior as amqp.Channel.Publish().
func (p *Publisher) Publish(
	ctx context.Context,
	exchange string,
	key string,
	mandatory bool,
	msg amqp.Publishing,
) error {
	msg.DeliveryMode = 2 // persistent

	ch, err := p.acquire()
	if err != nil {
		return err
	}

	if err := ch.Channel.Publish(
		exchange,
		key,
		mandatory,
		false, // immediate
		msg,
	); err != nil {
		return err
	}

	select {
	case r := <-ch.Return:
		go p.confirmThenRelease(ch)
		return errors.New("broker could not route message, " + r.ReplyText)

	case c := <-ch.Confirm:
		// wait for a confirmation from the broker, once we receive one
		// (regardless of whether it's an ack or a nack) we can return the
		// channel to the pool
		p.release(ch)

		if c.Ack {
			return nil
		}

		// there's no more meaningful error to be returned here. The RMQ docs
		// simply say: "basic.nack will only be delivered if an internal error
		// occurs in the Erlang process responsible for a queue."
		return errors.New("broker did not confirm message publication")

	case err := <-ch.Close:
		// if the channel is closed before we receive the confirmation, we do
		// not return the channel to the pool
		return err

	case <-ctx.Done():
		// if our context is canceled before we receive the confirmation, return
		// the channel the pool only after our confirmation has been consumed.
		go p.confirmThenRelease(ch)
		return ctx.Err()
	}
}

// acquire gets a channel from the pool, or opens a new channel and places it
// into "confirm mode" if the pool is empty.
func (p *Publisher) acquire() (*channel, error) {
	select {
	case ch := <-p.channels:
		return ch, nil
	default:
	}

	c, err := p.conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := c.Confirm(false); err != nil { // false = noWait
		return nil, err
	}

	ch := &channel{
		Channel: c,
		Close:   make(chan *amqp.Error),
		Return:  make(chan amqp.Return, 1),
		Confirm: make(chan amqp.Confirmation, 1),
	}

	c.NotifyClose(ch.Close)
	c.NotifyReturn(ch.Return)
	c.NotifyPublish(ch.Confirm)

	return ch, nil
}

// release returns a channel to the pool, or closes it if the pool is full.
func (p *Publisher) release(ch *channel) {
	select {
	case p.channels <- ch:
	default:
		_ = ch.Channel.Close()
	}
}

// confirmThenRelease waits for the next confirm on ch before returning it to
// the pool. This ensures that some future publisher doesn't see a previous
// call's confirmation message as its own.
func (p *Publisher) confirmThenRelease(ch *channel) {
	select {
	case <-ch.Confirm:
		p.release(ch)
	case <-ch.Close:
	}
}
