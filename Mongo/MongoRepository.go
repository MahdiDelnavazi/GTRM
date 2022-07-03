package Mongo

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mahdidl/Database-benchmark/Mongo/Entity"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"net/http"
)

var benchId primitive.ObjectID
var counter = 0

// CreateMongoBench is a func for create bench seed
func CreateMongoBench() {
	benchBson := Entity.MongoBenchEntity{Name: "mimdl", Counter: 0}
	mongodb := MongoClient.Database(MongoName)
	BenchCollection := mongodb.Collection("Bench")
	result, queryError := BenchCollection.InsertOne(context.Background(), benchBson)
	if queryError != nil {
		fmt.Println(queryError)
		return
	}
	benchId = result.InsertedID.(primitive.ObjectID)
}

// BenchIncrement is a handler for endpoint to increment one shared variable in one document
func BenchIncrement(ginContext *gin.Context) {

	// setup mongo transaction
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()

	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := MongoClient.StartSession()
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, err)
		panic(err)
		return
	}
	defer session.EndSession(context.TODO())

	// start transaction
	err = mongo.WithSession(context.TODO(), session, func(sessionContext mongo.SessionContext) error {
		err = session.StartTransaction(txnOpts)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, err)
			return err
		}

		mongodb := MongoClient.Database(MongoName)
		BenchCollection := mongodb.Collection("Bench")

		// read document for increment counter on it
		var BenchFromDb Entity.MongoBenchEntity
		readResult := BenchCollection.FindOne(sessionContext, bson.M{"_id": benchId}).Decode(&BenchFromDb)
		if readResult != nil {
			return readResult
		}
		counter = counter + 1

		// connect to rabbitmq
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnErrorRabbitmq(err, "Failed to connect to RabbitMQ", ginContext)
		defer conn.Close()

		// open channel
		ch, err := conn.Channel()
		failOnErrorRabbitmq(err, "Failed to open a channel", ginContext)
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"hello",
			false,
			false,
			false,
			false,
			nil,
		)
		failOnErrorRabbitmq(err, "Failed to declare a queue", ginContext)

		buf := &bytes.Buffer{}
		gob.NewEncoder(buf).Encode(benchId.Hex())
		body := buf.Bytes()

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnErrorRabbitmq(err, "Failed to publish a message", ginContext)

		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, err)
			return err
		}
		if err = session.CommitTransaction(sessionContext); err != nil {
			ginContext.JSON(http.StatusInternalServerError, err)
			return err
		}
		ginContext.JSON(http.StatusOK, "ok")
		return nil
	})
	if err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			ginContext.JSON(http.StatusInternalServerError, abortErr)
			panic(abortErr)
		}
		ginContext.JSON(http.StatusInternalServerError, err)
		panic(err)
	}

}

// failOnErrorRabbitmq is a func for handle rabbitmq errors
func failOnErrorRabbitmq(err error, msg string, context *gin.Context) {
	if err != nil {
		fmt.Errorf("error %s: %s", msg, err)
		context.JSON(http.StatusInternalServerError, err)
		context.Done()
	}
}
