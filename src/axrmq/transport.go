package axrmq

import (
	"context"
	"fmt"
	"runtime"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/streadway/amqp"
)

// DefaultConcurrency is the default number of messages to handle concurrently.
var DefaultConcurrency = runtime.NumCPU() * 2

// Transport is an implementation of bus.Transport that uses RabbitMQ to
// communicate messages between endpoints.
type Transport struct {
	Conn        *amqp.Connection
	Concurrency int

	ep    string
	ch    *amqp.Channel
	pub   *Publisher
	msgs  <-chan amqp.Delivery
	close chan *amqp.Error
}

// Initialize sets up the transport to communicate as an endpoint named ep.
func (t *Transport) Initialize(ctx context.Context, ep string) error {
	ch, err := t.Conn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		// close the channel if it has not been "captured" by the transport for
		// continued use.
		if t.ch != ch {
			_ = ch.Close()
		}
	}()

	err = setupTopology(ch, ep)
	if err != nil {
		return err
	}

	n := t.Concurrency
	if n == 0 {
		n = DefaultConcurrency
	}

	err = ch.Qos(n, 0, false)
	if err != nil {
		return err
	}

	queue, _ := queueNames(ep)

	t.msgs, err = ch.Consume(
		queue,
		ep,    // consumer tag
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}

	t.ep = ep
	t.ch = ch
	t.pub = NewPublisher(t.Conn, n)
	t.close = make(chan *amqp.Error)
	ch.NotifyClose(t.close)

	return nil
}

// Subscribe instructs the transport to listen to multicast messages of the
// given type.
func (t *Transport) Subscribe(ctx context.Context, mt ax.MessageTypeSet) error {
	return setupMulticastBindings(t.ch, t.ep, mt)
}

// SendMessage sends a message.
func (t *Transport) SendMessage(ctx context.Context, m bus.OutboundEnvelope) error {
	var pub amqp.Publishing

	if err := marshalMessage(t.ep, m, &pub); err != nil {
		fmt.Println(err)
		return err
	}

	switch m.Operation {
	case bus.OpSendUnicast:
		return t.sendUnicast(ctx, m.DestinationEndpoint, pub)
	case bus.OpSendMulticast:
		return t.sendMulticast(ctx, pub)
	default:
		panic(fmt.Sprintf("unrecognized outbound operation: %d", m.Operation))
	}
}

// ReceiveMessage returns the next message that has been delivered to the
// endpoint.
func (t *Transport) ReceiveMessage(ctx context.Context) (bus.InboundEnvelope, error) {
	for {
		select {
		case del := <-t.msgs:
			m, ok, err := t.receive(ctx, del)
			if ok || err != nil {
				return m, err
			}
		case err := <-t.close:
			return bus.InboundEnvelope{}, err
		case <-ctx.Done():
			return bus.InboundEnvelope{}, ctx.Err()
		}
	}
}

func (t *Transport) receive(
	ctx context.Context,
	del amqp.Delivery,
) (bus.InboundEnvelope, bool, error) {
	m := bus.InboundEnvelope{
		Done: func(_ context.Context, op bus.InboundOperation) error {
			switch op {
			case bus.OpAck:
				return del.Ack(false) // false = single message
			case bus.OpRetry:
				return del.Reject(true) // true = requeue
			case bus.OpReject:
				return del.Reject(false) // false = don't requeue
			default:
				panic(fmt.Sprintf("unrecognized inbound operation: %d", op))
			}
		},
	}

	if err := unmarshalMessage(del, &m); err != nil {
		// TODO: sentry, etc?
		return bus.InboundEnvelope{}, false, del.Reject(false)
	}

	return m, true, nil
}

// sendUnicast sends a unicast message directly to a specific endpoint.
func (t *Transport) sendUnicast(
	ctx context.Context,
	ep string,
	pub amqp.Publishing,
) error {
	return t.pub.Publish(
		ctx,
		unicastExchange,
		ep,   // routing key
		true, // mandatory
		pub,
	)
}

// sendMulticast sends a multicast message to the its subscribers.
func (t *Transport) sendMulticast(
	ctx context.Context,
	pub amqp.Publishing,
) error {
	return t.pub.Publish(
		ctx,
		multicastExchange,
		pub.Type, // routing key
		false,    // mandatory
		pub,
	)
}
