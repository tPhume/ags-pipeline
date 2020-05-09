package main

import (
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"github.com/tPhume/ags-pipeline/consumer"
	"github.com/tPhume/ags-pipeline/sensor"
	"log"
	"strings"
)

func main() {
	var err error

	// Setup configuration
	viper.SetConfigFile("example.env")

	err = viper.ReadInConfig()
	failOnError("could not read example.env", err)

	// Load environment variables
	rabbitURI := viper.GetString("RABBIT_URI")
	queueName := viper.GetString("QUEUE_NAME")

	failOnEmpty(rabbitURI, queueName)

	// Connect and consume messages from RabbitMQ
	// Create RabbitMQ connection
	log.Println("Creating a connection to RabbitMQ")
	conn, err := amqp.Dial(rabbitURI)
	failOnError("could not connect to rabbitmq", err)

	// Create new channel from RabbitMQ connection
	log.Println("Opening a new channel")
	ch, err := conn.Channel()
	failOnError("could not open new channel", err)

	// Register a consumer
	log.Println("Registering a consumer")
	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError("could not register consumer", err)

	// Create sensor.RabbitMQ with sensor.Stdout
	stdout := &sensor.Stdout{}
	rabbitMQSensor := &sensor.RabbitMQ{DataSource: stdout}

	// Create consumer.Listener type
	listener := consumer.Listener{Stream: msgs, Handle: rabbitMQSensor.Write}

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
