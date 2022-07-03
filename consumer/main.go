package main

import (
	"context"
	"fmt"
	"github.com/mahdidl/Database-benchmark/Mongo/Entity"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func main() {

	MongoClient, ctx, err := connect("mongodb://localhost:27017/")
	err = Ping(MongoClient, ctx)

	// connect to rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErrorRabbitmq(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErrorRabbitmq(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnErrorRabbitmq(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnErrorRabbitmq(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for msg := range msgs {

			// get mongo collection
			mongodb := MongoClient.Database("database-benchmark")
			BenchCollection := mongodb.Collection("Bench")

			// get bench id from body
			id, _ := primitive.ObjectIDFromHex(string(msg.Body))
			filter := bson.D{{"_id", id}}

			// get count bench from db
			var BenchFromDb Entity.MongoBenchEntity
			_ = BenchCollection.FindOne(context.TODO(), filter).Decode(&BenchFromDb)

			// increase count and update bench
			counter := BenchFromDb.Counter + 1
			result, queryError := BenchCollection.UpdateOne(
				context.TODO(),
				filter,
				bson.D{
					{"$set", bson.D{{"Counter", counter}}},
				},
			)
			if queryError != nil {
				// should retry query when having error
				err = retry(ch, q, msg)
				failOnErrorRabbitmq(err, "query error")
			}

			if result.ModifiedCount == 0 {
				// should retry query when query doesn't work
				err = retry(ch, q, msg)
				failOnErrorRabbitmq(err, "Failed to publish a message")
			}

			fmt.Println("done ", time.Now(), result.ModifiedCount)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// retry is a func for retry our process in rabbitmq
func retry(ch *amqp.Channel, q amqp.Queue, d amqp.Delivery) error {
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(d.Body),
		})
	return err
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

// failOnErrorRabbitmq is a func for handle rabbitmq errors
func failOnErrorRabbitmq(err error, msg string) {
	if err != nil {
		fmt.Errorf("error %s: %s", msg, err)
	}
}
