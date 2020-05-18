package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/spf13/viper"
	"log"
	"strings"
)

const queryString = `from(bucket: "production/autogen")
  |> range(start: 2020-05-18T00:00:00Z, stop: 2020-05-18T23:59:00Z)
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

	// Print result
	for result.Next() {
		fmt.Printf("%v\n\n", result.Record())
	}
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
