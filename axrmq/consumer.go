package axrmq

import (
	"context"
	"sync"

	"github.com/jmalloc/ax"
	"github.com/streadway/amqp"
)

// consumer receives messages from the broker
type consumer struct {
	ep          string
	minPreFetch int
	maxPreFetch int

	m        sync.Mutex
	ch       *amqp.Channel
	preFetch int

	msgs  <-chan amqp.Delivery
	close chan *amqp.Error
}

func newConsumer(
	conn *amqp.Connection,
	ep string,
	excl bool,
	minPreFetch, maxPreFetch int,
) (*consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer func() {
		if ch != nil {
			// Error ignored: a more significant error has already occurred if
			// we are exiting without having created a consumer.
			_ = ch.Close()
		}
	}()

	err = declareQueues(ch, ep)
	if err != nil {
		return nil, err
	}

	err = ch.Qos(minPreFetch, 0, false)
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
		ep:          ep,
		preFetch:    minPreFetch,
		minPreFetch: minPreFetch,
		maxPreFetch: maxPreFetch,
		ch:          ch,
		msgs:        msgs,
		close:       make(chan *amqp.Error),
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

func (c *consumer) IncreasePreFetch() error {
	return c.updatePreFetch(+1)
}

func (c *consumer) DecreasePreFetch() error {
	return c.updatePreFetch(-1)
}

func (c *consumer) updatePreFetch(delta int) error {
	c.m.Lock()
	defer c.m.Unlock()

	pf := c.preFetch + delta

	if pf < c.minPreFetch {
		pf = c.minPreFetch
	} else if pf > c.maxPreFetch {
		pf = c.maxPreFetch
	}

	if pf == c.preFetch {
		return nil
	}

	if err := c.ch.Qos(pf, 0, false); err != nil {
		return err
	}

	c.preFetch = pf

	return nil
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
