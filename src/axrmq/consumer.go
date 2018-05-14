package axrmq

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/streadway/amqp"
)

// consumer receives messages from the broker
type consumer struct {
	ep    string
	ch    *amqp.Channel
	msgs  <-chan amqp.Delivery
	close chan *amqp.Error
}

func newConsumer(
	conn *amqp.Connection,
	ep string,
	excl bool,
	preFetch int,
) (*consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer func() {
		if ch != nil {
			_ = ch.Close()
		}
	}()

	err = declareQueues(ch, ep)
	if err != nil {
		return nil, err
	}

	err = ch.Qos(preFetch, 0, false)
	if err != nil {
		return nil, err
	}

	queue, _ := queueNames(ep)

	msgs, err := ch.Consume(
		queue,
		ep,    // consumer tag
		false, // autoAck
		excl,  // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	con := &consumer{
		ep,
		ch,
		msgs,
		make(chan *amqp.Error),
	}

	ch.NotifyClose(con.close)
	ch = nil

	return con, nil
}

func (c *consumer) BindUnicast(mt ax.MessageTypeSet) error {
	return declareUnicastBindings(c.ch, c.ep, mt)
}

func (c *consumer) BindMulticast(mt ax.MessageTypeSet) error {
	return declareMulticastBindings(c.ch, c.ep, mt)
}

func (c *consumer) Receive(ctx context.Context) (amqp.Delivery, error) {
	select {
	case del := <-c.msgs:
		return del, nil
	case err := <-c.close:
		return amqp.Delivery{}, err
	case <-ctx.Done():
		return amqp.Delivery{}, ctx.Err()
	}
}
