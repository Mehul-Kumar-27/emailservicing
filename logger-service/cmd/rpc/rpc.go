package rpc

import (
	"context"
	"log"
	data "logger-service/cmd/data"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

type RPCServer struct {
}

func NewRpcServer(mongo *mongo.Client) *RPCServer {
	client = mongo

	return &RPCServer{}
}

type RPCPayload struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
}

func (app *RPCServer) LogInfo(payload RPCPayload, response *string) error {
	log.Println("RPC call received")
	collection := client.Database("logger").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreateAt:  time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	*response = "Log entry created using RPC call" + payload.Name
	return nil
}

func (app *RPCServer) Nothing() {
	log.Println("Nothing")
}
