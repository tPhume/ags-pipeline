package main

import (
	"github.com/spf13/viper"
	"github.com/tPhume/ags-pipeline/sensor"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
	"math"
	"strings"
	"time"
)

func main() {
	var err error

	// Setup configuration
	viper.SetConfigFile("rethinkdb.env")

	err = viper.ReadInConfig()
	failOnError("could not read rethinkdb.env", err)

	// Load environment variables
	rethinkURL := viper.GetString("RETHINK_URL")
	rethinkDatabase := viper.GetString("RETHINK_DATABASE")

	failOnEmpty(rethinkURL, rethinkDatabase)

	// Connect to RethinkDb
	session, err := rethinkdb.Connect(rethinkdb.ConnectOpts{Address: rethinkURL, Database: rethinkDatabase})
	failOnError("could not connect to rethinkdb", err)

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

	err = rethinkdb.Table("sensor").Insert(map[string]interface{}{
		"user_id":       "test-user-id",
		"controller_id": "test-controller-id",
		"temperature":   roundFloat(mock.Data.Temperature),
		"humidity":      roundFloat(mock.Data.Humidity),
		"light":         roundFloat(mock.Data.Light),
		"soil_moisture": mock.Data.SoilMoisture,
		"water_level":   mock.Data.WaterLevel,
		"time":          time.Now().Unix(),
	}, rethinkdb.InsertOpts{Conflict: "replace"}).Exec(session)
	failOnError("could not insert data", err)
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
