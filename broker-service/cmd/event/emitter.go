package events

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (emitter *Emitter) setup() error {
	channel, err := emitter.connection.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

func (emitter *Emitter) Push(event string, severity string) error {
	channel, err := emitter.connection.Channel()
	if err != nil {
		return err
	}

	log.Println("Pushing event: ", event, " with severity: ", severity)
	ctx := context.Background()
	error := channel.PublishWithContext(ctx, "logs-topic", severity, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(event),
	})
	if error != nil {
		return error
	}

	listner := channel.NotifyReturn(make(chan amqp.Return))
	go ReturnListner(listner)

	return nil
}

func ReturnListner(listner chan amqp.Return) {
	for ret := range listner {
		log.Println("Message returned: ", ret)
	}
}

func NewEmitter(connection *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: connection,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
