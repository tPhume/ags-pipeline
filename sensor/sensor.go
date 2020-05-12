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

// Meta describes the information about the controller
type Meta struct {
	ControllerId string `bson:"_id"`
	UserId       string `bson:"user_id"`
}

// Interacts with out data source
type Storage interface {
	Write(ctx context.Context, meta *Meta, msg *Message) error
}

// Meta Storage is where we retrieve metadata base on the controller token
type MetaStorage interface {
	Get(ctx context.Context, token string, meta *Meta) error
}

var ErrBadToken = errors.New("token not found in storage")

// If data is from RabbitMQ - use this
// Implements consumer.DataSource
type RabbitMQ struct {
	Validator   *validator.Validate
	Storage     Storage
	MetaStorage MetaStorage
}

// Handles write operation given that value is amqp.Delivery
func (r *RabbitMQ) Write(ctx context.Context) error {
	// Get value from context
	delivery, ok := ctx.Value("msg").(amqp.Delivery)
	if !ok {
		return errors.New("incorrect data type")
	}

	// Get data from message
	msg := &Message{}
	if err := json.Unmarshal(delivery.Body, msg); err != nil {
		log.Printf("cannot unmarshal message[%s], err: %s\n", delivery.MessageId, err.Error())
		_ = delivery.Nack(false, false)
		return errors.New("problem decoding data for message " + delivery.MessageId)
	}

	// Validate the data
	if err := r.Validator.Struct(msg); err != nil {
		_ = delivery.Nack(false, false)
		return err
	}

	// Get metadata
	meta := &Meta{}
	if err := r.MetaStorage.Get(ctx, msg.Token, meta); err != nil {
		if err == ErrBadToken {
			log.Printf("bad token for message[%s]\n", delivery.MessageId)
			_ = delivery.Nack(false, false)
			return err
		}

		log.Printf("problem occured with meta storage\n")
		_ = delivery.Nack(false, true)
		return err
	}

	// Write data to storage
	if err := r.Storage.Write(ctx, meta, msg); err != nil {
		log.Printf("cannot write message[%s] , err: %s\n", delivery.MessageId, err.Error())
		_ = delivery.Nack(false, true)
		return consumer.ErrFatal
	}
	log.Printf("written message %s to data source\n", delivery.MessageId)

	// Ack the message
	if err := delivery.Ack(true); err != nil {
		log.Printf("cannot ack messag[%s], err: %s", delivery.MessageId, err.Error())
		_ = delivery.Nack(false, true)
		return consumer.ErrFatal
	}
	log.Printf("message[%s] acked", delivery.MessageId)

	return nil
}
