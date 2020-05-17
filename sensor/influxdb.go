package sensor

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"strconv"
	"time"
)

type Influxdb struct {
	WriteApi influxdb2.WriteApiBlocking
}

func (i *Influxdb) Write(ctx context.Context, meta *Meta, msg *Message) error {
	t, err := strconv.ParseInt(msg.Data.Timestamp, 10, 64)
	if err != nil {
		return ErrBadTimeStamp
	}

	if err := i.WriteApi.WritePoint(ctx, influxdb2.NewPoint(
		"sensor",
		map[string]string{"user_id": meta.UserId, "controller_id": meta.ControllerId},
		map[string]interface{}{
			"temperature":   msg.Data.Temperature,
			"humidity":      msg.Data.Humidity,
			"light":         msg.Data.Light,
			"soil_moisture": msg.Data.SoilMoisture,
			"water_level":   msg.Data.WaterLevel,
		},
		time.Unix(t, 0),
	)); err != nil {
		return err
	}

	return nil
}
