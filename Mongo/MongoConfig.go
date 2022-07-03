package Mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var (
	MongoName    = "database-benchmark"
	MongoContext context.Context
	MongoClient  *mongo.Client
	MongoUrl     = "mongodb://localhost:27017/"
)

// MongoConfig is a func for configuration mongodb
func MongoConfig() {
	var err error
	MongoClient, MongoContext, err = connect(MongoUrl)
	if err != nil {
		log.Println("fatal 6")
		log.Fatal(err)
	}

	Ping(MongoClient, MongoContext)
}

// connect is a func fot connect to mongodb
func connect(uri string) (*mongo.Client, context.Context, error) {
	context := context.TODO()
	client, err := mongo.Connect(context, options.Client().ApplyURI(uri))
	return client, context, err
}

// Ping is a func for make sure we connect to mongodb successfully
func Ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}
