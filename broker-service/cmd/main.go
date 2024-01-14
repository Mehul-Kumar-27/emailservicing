// main.go
package main

import (
	handellers "broker/cmd/handellers"
	"fmt"
	"log"
	"net/http"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	webPort = "8080"
)

type RabbbitMQconnection struct {
	conn  *amqp.Connection
	error error
}

func main() {
	///// Connect to RabbitMQ
	connectionChan := make(chan RabbbitMQconnection)
	go connectToRabbitMQ(connectionChan)
	rabbitMQconnection := <-connectionChan
	if rabbitMQconnection.error != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", rabbitMQconnection.error)
	} else {
		log.Println("Connected to RabbitMQ")
	}
	/// Star
	defer rabbitMQconnection.conn.Close()

	logger := log.New(os.Stdout, "API", log.Lshortfile)
	logger.Printf("Starting API server %s", webPort)

	app := handellers.NewServerModel(
		rabbitMQconnection.conn,
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

func connectToRabbitMQ(connectionChan chan RabbbitMQconnection) {
	var url = "amqp://guest:guest@rabbitmq:5672/"
	count := 10

	for {
		conn, err := amqp.Dial(url)
		if err == nil {
			log.Println("RabbitMQ is ready to accept connections")
			connectionChan <- RabbbitMQconnection{conn: conn, error: nil}

		} else {
			if count <= 10 {
				log.Printf("RabbitMQ connection is probably not ready yet. error: %s", err)
			} else {
				connectionChan <- RabbbitMQconnection{conn: nil, error: err}
				break
			}
		}
	}

	log.Println("Failed to connect to RabbitMQ")
}
