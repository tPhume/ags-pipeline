package sensor

import (
	"context"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
	"time"
)

type Rethinkdb struct {
	Session *rethinkdb.Session
}

func (r *Rethinkdb) Write(ctx context.Context, meta *Meta, msg *Message) error {
	// we ignore errors
	_ = rethinkdb.Table("sensor").Insert(map[string]interface{}{
		"user_id":       meta.UserId,
		"controller_id": meta.ControllerId,
		"temperature":   msg.Data.Temperature,
		"humidity":      msg.Data.Humidity,
		"light":         msg.Data.Light,
		"soil_moisture": msg.Data.SoilMoisture,
		"water_level":   msg.Data.WaterLevel,
		"time":          time.Now().Unix(),
	}, rethinkdb.InsertOpts{Conflict: "replace"}).Exec(r.Session)

	return nil
}
