// main.go
package main

import (
	handellers "broker/cmd/handellers"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	webPort = "8080"
)

type ConnectionStruct struct {
	Connection *amqp.Connection
	Error      error
}

func main() {
	///// Connect to RabbitMQ
	ch := make(chan ConnectionStruct)
	go connectToRabbitMQ(ch)
	connectionStruct := <-ch
	if connectionStruct.Error != nil {
		log.Fatal(connectionStruct.Error)
	} else {
		log.Println("Connected to RabbitMQ")
	}
	defer connectionStruct.Connection.Close()

	logger := log.New(os.Stdout, "API", log.Lshortfile)
	logger.Printf("Starting API server %s", webPort)

	app := handellers.NewServerModel(
		connectionStruct.Connection,
	)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort), // Fix the format specifier
		Handler: app.Routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
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