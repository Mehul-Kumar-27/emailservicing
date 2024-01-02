package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/cmd/api/handellers"
	"logger-service/cmd/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	port     = "8080"
	rpcPort  = "5001"
	mongoURI = "mongodb://localhost:27017"
	gRPCPort = "50001"
)

var client = mongo.Client{}

func main() {
	// connect to mongo
	mongoClient, err := connectTOmongo()
	if err != nil {
		log.Panic("Error while connecting to mongo: ", err)
	}
	client = *mongoClient

	app := handellers.LoggerService{
		Modles: data.New(&client),
	}

	go startServer(app)

	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancle()

	// close the connection to mongo
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Panic("Error while disconnecting from mongo: ", err)
		}
	}()

	/// start the serer

}

func startServer(app handellers.LoggerService) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.Routes(),
	}

	log.Println("Starting server on port: ", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Panic("Error while starting server: ", err)
	}
}

func connectTOmongo() (*mongo.Client, error) {
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
			panic(err)
		}
	}()

	return connect, nil
}
