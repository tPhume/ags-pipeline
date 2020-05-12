package sensor

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/streadway/amqp"
	"github.com/tPhume/ags-pipeline/consumer"
	"log"
)

// Message represent the data from sensors
type Message struct {
	Token string `json:"token" validate:"uuid4"`
	Data  Data   `json:"data" validate:"required"`
}

// Data are raw data points included within message
type Data struct {
	Timestamp    string  `json:"timestamp"`
	Temperature  float64 `json:"temperature"`
	Humidity     float64 `json:"humidity" validate:"gte=0,lte=100"`
	Light        float64 `json:"light" validate:"gte=0,lte=65535"`
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
	Validator  *validator.Validate
	DataSource DataSource
}

// Handles write operation given that value is amqp.Delivery
func (r *RabbitMQ) Write(ctx context.Context) error {
	delivery, ok := ctx.Value("msg").(amqp.Delivery)
	if !ok {
		return errors.New("incorrect data type")
	}

	msg := &Message{}
	if err := json.Unmarshal(delivery.Body, msg); err != nil {
		log.Printf("cannot unmarshal message[%s], err: %s\n", delivery.MessageId, err.Error())
		_ = delivery.Nack(false, false)
		return errors.New("problem decoding data for message " + delivery.MessageId)
	}

	if err := r.Validator.Struct(msg); err != nil {
		_ = delivery.Nack(false, false)
		return err
	}

	if err := r.DataSource.Write(ctx, msg); err != nil {
		log.Printf("cannot write message[%s] , err: %s\n", delivery.MessageId, err.Error())
		_ = delivery.Nack(false, true)
		return consumer.ErrFatal
	}
	log.Printf("written message %s to data source\n", delivery.MessageId)

	if err := delivery.Ack(true); err != nil {
		log.Printf("cannot ack messag[%s], err: %s", delivery.MessageId, err.Error())
		_ = delivery.Nack(false, true)
		return consumer.ErrFatal
	}
	log.Printf("message[%s] acked", delivery.MessageId)

	return nil
}
