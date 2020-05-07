// Package stream allows creation of a Stream
// A Stream reads from a data source continuously
// It sends validates the data and sends it to another channel
package stream

import "github.com/streadway/amqp"

// Message represent the data that comes from the data source
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

// RabbitStream type holds a receive channel (amqp.Delivery) where it gets the data
// And a send channel to send data after validation
type RabbitStream struct {
	Receive <-chan amqp.Delivery
	Send    chan<- *Message
}

func (s *RabbitStream) ListenRabbitMQ() error {
	panic("implement me")
}
