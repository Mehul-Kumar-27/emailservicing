package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name,omitempty" json:"name,omitempty"`
	Data      string    `bson:"data,omitempty" json:"data,omitempty"`
	CreateAt  time.Time `bson:"create_at,omitempty" json:"create_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func (l *LogEntry) Create() error {
	collection := client.Database("logger").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      l.Name,
		Data:      l.Data,
		CreateAt:  time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error while inserting log entry: ", err)
		return err
	}

	return nil
}

func (l *LogEntry) GetAll() ([]*LogEntry, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancle()

	collection := client.Database("logger").Collection("logs")

	ops := options.Find()

	ops.SetSort(bson.D{{Key: "create_at", Value: -1}})

	cursor, err := collection.Find(ctx, bson.D{}, ops)

	if err != nil {
		log.Println("Error while getting all log entry: ", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var log LogEntry
		err := cursor.Decode(&log)
		if err != nil {
			fmt.Println("Error while decoding log entry: ", err)
			return nil, err
		}

		logs = append(logs, &log)

	}

	return logs, nil

}
