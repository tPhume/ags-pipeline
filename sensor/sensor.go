package sensor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/tPhume/ags-pipeline/consumer"
	"log"
)

// Message represent the data from sensors
type Message struct {
	Token string `json:"token" validate:"uuid4"`
	Data  []Data `json:"data" validate:"dive"`
}

// Data are raw data points included within message
type Data struct {
	Timestamp    string  `json:"timestamp"`
	Temperature  float32 `json:"temperature"`
	Humidity     float32 `json:"humidity" validate:"gte=0,lte=100"`
	Light        float32 `json:"light" validate:"gte=0,lte=65535"`
	SoilMoisture int     `json:"soil_moisture" validate:"gte=0,lte=1000"`
	WaterLevel   int     `json:"water_level" validate:"gte=0"`
}

// Interacts with out data source
type DataSource interface {
	Write(ctx context.Context, msg *Message) error
}

// If data is from RabbitMQ - use this
// Implements consumer.DataSource
type RabbitMQ struct {
	DataSource DataSource
}

// Handles write operation given that value is amqp.Delivery
func (r *RabbitMQ) Write(ctx context.Context) error {
	delivery, ok := ctx.Value("msg").(amqp.Delivery)
	if !ok {
		return errors.New("incorrect data type")
	}

	var msg *Message
	if err := json.Unmarshal(delivery.Body, msg); err != nil {
		return errors.New("problem decoding data for message " + delivery.MessageId)
	}

	if err := r.DataSource.Write(ctx, msg); err != nil {
		return consumer.ErrFatal
	}
	log.Println(fmt.Sprintf("written message %s to data source", delivery.MessageId))

	if err := delivery.Ack(true); err != nil {
		log.Println("cannot ack message " + delivery.MessageId)
		return consumer.ErrFatal
	}
	log.Println(fmt.Sprintf("message %s acked", delivery.MessageId))

	return nil
}
