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

func (m *Mongo) WriteMean(ctx context.Context, summary map[string]*Summary) error {
	today := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// Populate for bulk writing
	models := make([]mongo.WriteModel, 0)
	for _, v := range summary {
		// Create update query
		m := mongo.NewUpdateOneModel()

		m.SetUpsert(true)
		m.SetFilter(bson.M{"controller_id": v.ControllerId, "user_id": v.UserId, "date": today})
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

		models = append(models, m)
	}

	if _, err := m.Col.BulkWrite(ctx, models); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) WriteMedian(ctx context.Context, summary map[string]*Summary) error {
	today := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// Populate for bulk writing
	models := make([]mongo.WriteModel, 0)
	for _, v := range summary {
		// Create update query
		m := mongo.NewUpdateOneModel()

		m.SetUpsert(true)
		m.SetFilter(bson.M{"controller_id": v.ControllerId, "user_id": v.UserId, "date": today})
		m.SetUpdate(bson.M{
			"$set": bson.M{
				"median_humidity":      v.Data["humidity"],
				"median_light":         v.Data["light"],
				"median_soil_moisture": v.Data["soil_moisture"],
				"median_temperature":   v.Data["temperature"],
				"median_water_level":   v.Data["water_level"],
			},
		})

		models = append(models, m)
	}

	if _, err := m.Col.BulkWrite(ctx, models); err != nil {
		return err
	}

	return nil
}
