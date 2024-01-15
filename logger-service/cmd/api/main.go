package main

import (
	"context"
	"fmt"
	"log"
	r "logger-service/cmd/rpc"
	"net"
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
	"net/rpc"
)

const (
	port     = "8080"
	rpcPort  = "5001"
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
		err := mongoClient.Ping(context.TODO(), nil)
		if err != nil {
			log.Panic("Error while pinging mongo: ", err)
		} else {
			log.Println("Connected to mongo")
		}
	}
	client = *mongoClient

	// Create logger service
	app := handellers.LoggerService{
		Modles: data.New(&client),
	}

	rpcServerConnection := r.NewRpcServer(&client)

	e := rpc.Register(rpcServerConnection)
	if e != nil {
		log.Println("Error while registering RPC server:", e)
		return
	}
	// Create a channel to receive OS signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go rpcListen(signalCh)

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.Routes(),
	}

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
		Password: "password",
	})

	connect, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error while connecting to mongo: ", err)
		return nil, err
	}

	return connect, nil
}

func rpcListen(signalChan chan os.Signal) {
	log.Printf("Starting RPC server on port: %s", rpcPort)

	// Create a TCP listener
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Println("Error while creating TCP listener:", err)
		signalChan <- syscall.SIGTERM
	}
	go func() {
		<-signalChan
		listener.Close()
	}()
	defer listener.Close()

	for {
		rpcConnections, err := listener.Accept()
		if err != nil {
			log.Println("Error while accepting connection:", err)
			continue
		}
		log.Println("Connection accepted")
		go rpc.ServeConn(rpcConnections)
	}
}
