package main

import (
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/spf13/viper"
	"github.com/tPhume/ags-pipeline/summary"
	"log"
	"net/http/httptest"
	"strings"
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

	// Get query client
	queryApi := client.QueryApi("")

	// Create summary.Influx reader type and summary.Stdout writer type
	reader := &summary.Influx{QueryApi: queryApi}
	writer := &summary.Stdout{}

	storage := summary.Storage{
		Reader: reader,
		Writer: writer,
	}

	// Test handler
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	storage.Handle(c)
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
