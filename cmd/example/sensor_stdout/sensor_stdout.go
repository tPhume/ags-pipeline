package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"github.com/tPhume/ags-pipeline/consumer"
	"github.com/tPhume/ags-pipeline/sensor"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var err error

	// Setup configuration
	viper.SetConfigFile(os.Args[1])

	err = viper.ReadInConfig()
	failOnError("could not read env file", err)

	// Load environment variables
	rabbitURI := viper.GetString("RABBIT_URI")
	queue := viper.GetString("QUEUE")

	mongoURI := viper.GetString("MONGO_URI")
	mongoDb := viper.GetString("MONGO_DB")

	failOnEmpty(rabbitURI, queue, mongoURI, mongoDb)

	// Connect and consume messages from RabbitMQ
	// Create RabbitMQ connection
	log.Println("Creating a connection to RabbitMQ")
	conn, err := amqp.Dial(rabbitURI)
	failOnError("could not connect to rabbitmq", err)

	// Create new channel from RabbitMQ connection
	log.Println("Opening a new channel")
	ch, err := conn.Channel()
	failOnError("could not open new channel", err)

	// Register influx consumer
	log.Println("Registering queue consumers")
	msg, err := ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError("could not register queue consumer", err)

	// Connect to Mongodb
	// Get Database and Collection
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	failOnError("could not create mongo client", err)

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = mongoClient.Connect(timeout)
	failOnError("could not start mongo connection", err)

	mongoDatabase := mongoClient.Database(mongoDb)
	controllerCol := mongoDatabase.Collection("controller")

	// Create a validator
	v := validator.New()

	// Create sensor.RabbitMQ with sensor.Stdout and sensor.Mongo
	stdout := &sensor.Stdout{}
	metaMongo := &sensor.Mongodb{Col: controllerCol}
	rabbitMQSensor := &sensor.RabbitMQ{Validator: v, Storage: stdout, MetaStorage: metaMongo}

	// Create consumer.Listener type
	listener := consumer.Listener{Stream: msg, Handle: rabbitMQSensor.Write}

	log.Fatal(listener.Listen())
}

func failOnEmpty(env ...string) {
	for i, e := range env {
		if strings.TrimSpace(e) == "" {
			log.Fatal("missing env at position ", i)
		}
	}
}

func failOnError(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
