// Package consume allows creation of a Consumer
// A Consumer reads data from a channel created from Stream package
// Data is then written using InsertData interface
package consume

// Message represent the data (associated with a controller) that is read from some data source
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
