package summary

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Mongo struct {
	Col *mongo.Collection
}

func (m *Mongo) Write(ctx context.Context, summary map[string]*Summary) error {
	today := time.Now().Add(time.Hour * -12).Format("2006-01-02")

	// Populate for bulk writing
	models := make([]mongo.WriteModel, 0)
	for _, v := range summary {
		// Create update query
		m := mongo.NewUpdateOneModel()

		m.SetUpsert(true)
		m.SetFilter(bson.M{"controller_id": v.Id, "user_id": v.UserId})
		m.SetUpdate(bson.M{
			"$set": bson.M{
				"date":               today,
				"mean_humidity":      v.Data["humidity"],
				"mean_light":         v.Data["light"],
				"mean_soil_moisture": v.Data["soil_moisture"],
				"mean_temperature":   v.Data["temperature"],
				"mean_water_level":   v.Data["water_level"],
			},
		})
	}

	if _, err := m.Col.BulkWrite(ctx, models); err != nil {
		return err
	}

	return nil
}
