package axrmq

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/streadway/amqp"
)

const unicastExchange = "ax/unicast"
const multicastExchange = "ax/multicast"

// queueNames returns the name of the pending and error queue to use for the
// endpoint named ep.
func queueNames(ep string) (string, string) {
	return ep + "/pending", ep + "/error"
}

// setupTopology declares all exchanges and queues for the endpoint named ep.
func setupTopology(ch *amqp.Channel, ep string) error {
	if err := declareExchanges(ch); err != nil {
		return err
	}

	if err := declareQueues(ch, ep); err != nil {
		return err
	}

	return nil
}

// declareExchanges declares the unicast and multicast message exchanges.
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
			"x-dead-letter-routing-key": errors,
		},
	); err != nil {
		return err
	}

	if err := ch.QueueBind(
		pending,
		ep,
		unicastExchange,
		false, // noWait
		nil,   // args
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

// setupMulticastBindings sets up bindings such that the endpoint named ep
// receives published messages of the types in t.
func setupMulticastBindings(ch *amqp.Channel, ep string, t ax.MessageTypeSet) error {
	pending, _ := queueNames(ep)

	for _, mt := range t.Members() {
		if err := ch.QueueBind(
			pending,
			mt.Name,
			multicastExchange,
			false, // noWait
			nil,   // args
		); err != nil {
			return err
		}
	}

	return nil
}
