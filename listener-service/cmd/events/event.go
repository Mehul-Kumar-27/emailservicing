package events

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs-topic", // name
		"topic",	  // type
		true,		  // durable
		false,		  // auto-deleted
		false,		  // internal
		false,		  // no-wait
		nil,		  // arguments
	)
}

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"events", // name
		true,	  // durable lets rock the world around us
		false,	  // delete when unused
		false,	  // exclusive
		false,	  // no-wait
		nil,	  // arguments
	)
}
