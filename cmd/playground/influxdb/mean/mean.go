package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/spf13/viper"
	"github.com/tPhume/ags-pipeline/summary"
	"log"
	"math"
	"strings"
)

const queryString = `from(bucket: "production/autogen")
  |> range(start: 2020-05-17T17:00:00Z, stop: 2020-05-18T16:59:59Z)
  |> filter(fn: (r) => r._measurement == "sensor")
  |> mean()
  |> duplicate(column: "_stop", as: "_time")`

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

	// Get query client
	queryApi := client.QueryApi("")

	// Query data
	result, err := queryApi.Query(context.Background(), queryString)
	failOnError("could not query data", err)

	// Collect results
	data := make(map[string]*summary.Summary)
	for result.Next() {
		record := result.Record()

		// Get user_id and controller_id
		userId := record.ValueByKey("user_id").(string)
		controllerId := record.ValueByKey("controller_id").(string)

		// Add value to map
		_ = addToMap(data, userId, controllerId, record.Field(), record.ValueByKey("_value"))
	}

	// Print result
	for k, v := range data {
		fmt.Printf("---- Controller ID:%v\n%v\n\n", k, v)
	}
}

func addToMap(m map[string]*summary.Summary, userId string, controllerId string, field string, value interface{}) error {
	if _, exist := m[controllerId]; !exist {
		m[controllerId] = &summary.Summary{UserId: userId, ControllerId: controllerId, Data: make(summary.Data)}
	}

	m[controllerId].Data[field] = roundFloat(value.(float64))

	return nil
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
