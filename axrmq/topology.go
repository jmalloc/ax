package axrmq

import (
	"github.com/jmalloc/ax"
	"github.com/streadway/amqp"
)

const unicastExchange = "ax/unicast"
const multicastExchange = "ax/multicast"

// queueNames returns the name of the pending and error queue to use for the
// endpoint named ep.
func queueNames(ep string) (string, string) {
	return ep + "/pending", ep + "/error"
}

// declareExchanges declares the unicast and multicast message exchanges.
//
// Both are 'direct' exchanges, meaning that messages are routed based on exact
// match of the 'routing key'.
func declareExchanges(ch *amqp.Channel) error {
	if err := ch.ExchangeDeclare(
		unicastExchange,
		"direct",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait,
		nil,   // args,
	); err != nil {
		return err
	}

	if err := ch.ExchangeDeclare(
		multicastExchange,
		"direct",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait,
		nil,   // args,
	); err != nil {
		return err
	}

	return nil
}

// declareQueues declares the pending and error queues for the endpoint named ep.
func declareQueues(ch *amqp.Channel, ep string) error {
	pending, errors := queueNames(ep)

	if _, err := ch.QueueDeclare(
		pending,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": pending, // route dead-lettered messages back to the pending queue
		},
	); err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(
		errors,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	); err != nil {
		return err
	}

	return nil
}

// unicastRoutingKey returns the routing key to use when a unicast message is
// being sent to ep.
//
// The routing key is composed from both the message type and destination
// endpoint name. This ensures that firstly, endpoints to do not receive any
// messages they do not intend to handle, but more importantly that publishers
// receive an error when attempting to send a unicast message to an endpoint
// that will not handle the message. The latter is enforced by using the
// 'mandatory' flag in the publisher.
func unicastRoutingKey(mt string, ep string) string {
	return mt + ":" + ep
}

// multicastRoutingKey returns the routing key to use when a multicast message
// is sent.
func multicastRoutingKey(mt string) string {
	return mt
}

// declareUnicastBindings sets up bindings such that the endpoint named ep
// receives published messages of the types in t.
func declareUnicastBindings(ch *amqp.Channel, ep string, t ax.MessageTypeSet) error {
	pending, _ := queueNames(ep)

	for _, mt := range t.Members() {
		if err := ch.QueueBind(
			pending,
			unicastRoutingKey(mt.Name, ep),
			unicastExchange,
			false, // noWait
			nil,   // args
		); err != nil {
			return err
		}
	}

	return nil
}

// declareMulticastBindings sets up bindings such that the endpoint named ep
// receives published messages of the types in t.
func declareMulticastBindings(ch *amqp.Channel, ep string, t ax.MessageTypeSet) error {
	pending, _ := queueNames(ep)

	for _, mt := range t.Members() {
		if err := ch.QueueBind(
			pending,
			multicastRoutingKey(mt.Name),
			multicastExchange,
			false, // noWait
			nil,   // args
		); err != nil {
			return err
		}
	}

	return nil
}
