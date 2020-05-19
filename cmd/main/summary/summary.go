package main

import (
	"context"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/spf13/viper"
	"github.com/tPhume/ags-pipeline/summary"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
)

func main() {
	var err error

	// Setup configuration
	viper.SetConfigFile("main_summary.env")

	err = viper.ReadInConfig()
	failOnError("could not read main_summary.env", err)

	// Load environment variables
	influxURL := viper.GetString("INFLUX_URL")
	influxToken := viper.GetString("INFLUX_TOKEN")
	influxDatabase := viper.GetString("INFLUX_DATABASE")

	mongoUri := viper.GetString("MONGO_URI")
	mongoDb := viper.GetString("MONGO_DB")

	failOnEmpty(influxURL, influxToken, influxDatabase, mongoUri, mongoDb)

	// Create Influxdb client
	client := influxdb2.NewClient(influxURL, influxToken)
	defer client.Close()

	// Get query client
	queryApi := client.QueryApi("")

	// Create Mongo Client
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	failOnError("fail to create mongo client", err)

	err = mongoClient.Connect(context.Background())
	failOnError("fail to connect to mongo", err)

	mongoDatabase := mongoClient.Database(mongoDb)
	summaryCol := mongoDatabase.Collection("summary")

	// Create summary.Influx reader type and summary.Stdout writer type
	reader := &summary.Influx{QueryApi: queryApi}
	writer := &summary.Mongo{Col: summaryCol}

	storage := &summary.Storage{
		Reader: reader,
		Writer: writer,
	}

	// Init Gin
	engine := gin.New()
	engine.POST("api/v1/mean", storage.HandleMean)

	log.Fatal(engine.Run("0.0.0.0:8081"))
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
