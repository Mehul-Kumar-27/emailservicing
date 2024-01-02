package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	port     = "8080"
	rpcPort  = "5001"
	mongoURI = "mongodb://mongo:27017"
	gRPCPort = "50001"
)

var mongClient = mongo.Client{}

type LoggerService struct{}

func main() {
	// connect to mongo
	mongoClient, err := connectTOmongo()
	if err != nil {
		log.Panic("Error while connecting to mongo: ", err)
	}
	mongClient = *mongoClient

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
