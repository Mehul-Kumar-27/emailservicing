package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"logger-service/cmd/api/handellers"
	"logger-service/cmd/data"
	"net/http"
)

const (
	port     = "8080"
	mongoURI = "mongodb://mongo:27017"
)

var client = mongo.Client{}
var wg sync.WaitGroup

func main() {
	// Connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic("Error while connecting to mongo: ", err)
	} else {
		log.Println("Connected to mongo")
	}
	client = *mongoClient

	// Create logger service
	app := handellers.LoggerService{
		Modles: data.New(&client),
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.Routes(),
	}

	// Create a channel to receive OS signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Increment the WaitGroup counter
	wg.Add(1)

	// Start a goroutine to listen for OS signals
	go func() {
		defer wg.Done() // Decrement the WaitGroup counter when done

		sig := <-signalCh
		log.Printf("Received signal: %v. Shutting down...\n", sig)

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Shutdown the HTTP server
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Error during server shutdown:", err)
		}

		// Disconnect from MongoDB
		log.Println("Disconnecting from mongo")
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Error while disconnecting from mongo:", err)
		}
	}()

	// Start the HTTP server
	log.Println("Starting server on port: ", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Panic("Error while starting server: ", err)
	}

	// Wait for all goroutines to finish before exiting
	wg.Wait()
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "12345",
	})

	connect, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error while connecting to mongo: ", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		if err := connect.Disconnect(ctx); err != nil {
			log.Println("Error while disconnecting from mongo:", err)
		}
	}()

	return connect, nil
}
