package events

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	connection *amqp.Connection
	queueName  string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(connection *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		connection: connection,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil

}

func (consumer *Consumer) setup() error {
	channel, err := consumer.connection.Channel()
	if err != nil {
		log.Println("Failed to open a channel, err: ", err)
		return err
	}
	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.connection.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := declareQueue(ch)

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs-topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err

	}

	listenForMessages := make(chan struct{})

	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)

			var payload Payload
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				log.Println("Failed to unmarshal message: ", err)
			} else {
				go handelPayload(payload)
			}
		}
	}()

	log.Println("Listening for messages...")
	<-listenForMessages

	return nil

}

func handelPayload(payload Payload) {
	switch payload.Name {
	case "log":
		log.Println("Log message: ", payload.Data)
	default:
		log.Println("Unknown message type: ", payload.Name)
	}
}
