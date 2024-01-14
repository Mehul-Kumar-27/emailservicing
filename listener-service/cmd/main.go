package main

import (
	"listener-service/cmd/events"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConnectionStruct struct {
	Connection *amqp.Connection
	Error      error
}

func main() {
	//Connect to RabbitMQ

	ch := make(chan ConnectionStruct)
	go connectToRabbitMQ(ch)
	connectionStruct := <-ch
	if connectionStruct.Error != nil {
		log.Fatal(connectionStruct.Error)
	} else {
		log.Println("Connected to RabbitMQ")
	}
	defer connectionStruct.Connection.Close()

	//Let's start listning to the messages from the rabbitMQ
	consumer, err := events.NewConsumer(connectionStruct.Connection)
	if err != nil {
		log.Fatal(err)
	}

	// calling the listen function
	err = consumer.Listen([]string{"log.INFO", "log.ERROR", "log.WARNING"})
	if err != nil {
		log.Fatal(err)
	}

}

func connectToRabbitMQ(ch chan ConnectionStruct) {
	var connectionString = "amqp://guest:guest@rabitmq:5672/"
	var count = 0
	var sleepTime = 2 * time.Second

	for {
		connection, err := amqp.Dial(connectionString)
		if err != nil {
			count++
			log.Println("RabbbitMQ is probably not available yet. Recieved Errror: ", err)
		} else {
			rabbitMQConnection := ConnectionStruct{
				Connection: connection,
				Error:      nil,
			}
			ch <- rabbitMQConnection
		}

		if count > 10 {
			log.Println("RabbitMQ is not available. Giving up")
			rabbitMQConnection := ConnectionStruct{
				Connection: nil,
				Error:      err,
			}
			ch <- rabbitMQConnection
			break
		} else {
			time.Sleep(sleepTime)
		}

	}
}
