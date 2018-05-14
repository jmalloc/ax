package axrmq

import (
	"context"
	"fmt"
	"runtime"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/streadway/amqp"
)

// DefaultReceiveConcurrency is the default number of messages to process concurrently.
var DefaultReceiveConcurrency = runtime.NumCPU() * 2

// DefaultSendConcurrency is the default number of messages to send concurrently.
var DefaultSendConcurrency = runtime.NumCPU() * 10

// Transport is an implementation of bus.Transport that uses RabbitMQ to
// communicate messages between endpoints.
type Transport struct {
	Conn               *amqp.Connection
	Exclusive          bool
	SendConcurrency    int
	ReceiveConcurrency int

	ep  string
	pub *publisher
	con *consumer
}

// Initialize sets up the transport to communicate as an endpoint named ep.
func (t *Transport) Initialize(ctx context.Context, ep string) error {
	ch, err := t.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := declareExchanges(ch); err != nil {
		return err
	}

	t.ep = ep

	poolSize := t.SendConcurrency
	if poolSize == 0 {
		poolSize = DefaultSendConcurrency
	}

	t.pub = newPublisher(
		t.Conn,
		poolSize,
	)

	return nil
}

// Subscribe instructs the transport to listen to multicast messages of the
// given type.
func (t *Transport) Subscribe(ctx context.Context, op bus.Operation, mt ax.MessageTypeSet) error {
	if err := t.startConsumer(); err != nil {
		return err
	}

	switch op {
	case bus.OpSendUnicast:
		return t.con.BindUnicast(mt)
	case bus.OpSendMulticast:
		return t.con.BindMulticast(mt)
	default:
		panic(fmt.Sprintf("unrecognized outbound operation: %d", op))
	}
}

// Send sends env via the transport.
func (t *Transport) Send(ctx context.Context, env bus.OutboundEnvelope) error {
	var pub amqp.Publishing

	pub, err := marshalMessage(t.ep, env)
	if err != nil {
		return err
	}

	switch env.Operation {
	case bus.OpSendUnicast:
		return t.pub.PublishUnicast(ctx, pub, env.DestinationEndpoint)
	case bus.OpSendMulticast:
		return t.pub.PublishMulticast(ctx, pub)
	default:
		panic(fmt.Sprintf("unrecognized outbound operation: %d", env.Operation))
	}
}

// Receive returns the next message sent to this endpoint.
// It blocks until a message is available, or ctx is canceled.
func (t *Transport) Receive(ctx context.Context) (env bus.InboundEnvelope, ack bus.Acknowledger, err error) {
	err = t.startConsumer()
	if err != nil {
		return
	}

	var del amqp.Delivery

	for {
		del, err = t.con.Receive(ctx)
		if err != nil {
			return
		}

		env, err = unmarshalMessage(del)
		if err == nil {
			ack = &Acknowledger{del}
			return
		}

		// TODO: log / sentry / etc
		err = del.Reject(false) // false = don't requeue
		if err != nil {
			return
		}
	}
}

func (t *Transport) startConsumer() error {
	if t.con != nil {
		return nil
	}

	preFetch := t.ReceiveConcurrency
	if preFetch == 0 {
		preFetch = DefaultReceiveConcurrency
	}

	con, err := newConsumer(t.Conn, t.ep, t.Exclusive, preFetch)
	if err != nil {
		return err
	}

	t.con = con

	return nil
}
