package summary

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"math"
	"time"
)

type Influx struct {
	QueryApi influxdb2.QueryApi
}

func (i *Influx) Read(ctx context.Context, summary map[string]*Summary) error {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().Add(time.Hour * -12).Format("2006-01-02")

	queryString := fmt.Sprintf(`from(bucket: "production/autogen")
  |> range(start: %sT17:00:00Z, stop: %sT16:59:59Z)
  |> filter(fn: (r) => r._measurement == "sensor")
  |> mean()
  |> duplicate(column: "_stop", as: "_time")`, today, yesterday)

	// Query data
	result, err := i.QueryApi.Query(context.Background(), queryString)
	if err != nil {
		return err
	}

	for result.Next() {
		record := result.Record()

		// Get user_id and controller_id
		userId := record.ValueByKey("user_id").(string)
		controllerId := record.ValueByKey("controller_id").(string)

		// Add value to map
		add(summary, userId, controllerId, record.Field(), record.ValueByKey("_value"))
	}

	return nil
}

func add(m map[string]*Summary, userId string, controllerId string, field string, value interface{}) {
	if _, exist := m[controllerId]; !exist {
		m[controllerId] = &Summary{UserId: userId, ControllerId: controllerId, Data: make(Data)}
	}

	m[controllerId].Data[field] = roundFloat(value.(float64))
}

func roundFloat(x float64) float64 {
	return math.Round(x*100) / 100
}