package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/spf13/viper"
	"github.com/tPhume/ags-pipeline/sensor"
	"log"
	"math"
	"strings"
	"time"
)

func main() {
	var err error

	// Setup configuration
	viper.SetConfigFile("influxdb.env")

	err = viper.ReadInConfig()
	failOnError("could not read influxdb.env", err)

	// Load environment variables
	influxURL := viper.GetString("INFLUX_URL")
	influxToken := viper.GetString("INFLUX_TOKEN")
	influxDatabase := viper.GetString("INFLUX_DATABASE")

	failOnEmpty(influxURL, influxToken, influxDatabase)

	// Create Influxdb client
	client := influxdb2.NewClient(influxURL, influxToken)
	defer client.Close()

	// Create an influxdb writer
	influxWriter := client.WriteApiBlocking("", influxDatabase)

	// Mock data
	mock := sensor.Message{
		Token: "test-token",
		Data: sensor.Data{
			Temperature:  25.4,
			Humidity:     50.5,
			Light:        1000,
			SoilMoisture: 300,
			WaterLevel:   511,
		},
	}

	// Write new data point to influxdb
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = influxWriter.WritePoint(ctx, influxdb2.NewPoint(
		"sensor",
		map[string]string{"user_id": "test-user-id", "controller_id": "test-controller-id"},
		map[string]interface{}{
			"temperature":   roundFloat(mock.Data.Temperature),
			"humidity":      roundFloat(mock.Data.Humidity),
			"light":         roundFloat(mock.Data.Light),
			"soil_moisture": mock.Data.SoilMoisture,
			"water_level":   mock.Data.WaterLevel,
		},
		time.Now(),
	))
	failOnError("could not write point", err)
}

func roundFloat(x float64) float64 {
	return math.Round(x*100) / 100
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
