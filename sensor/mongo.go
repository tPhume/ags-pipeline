package sensor

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongodb struct {
	Col     *mongo.Collection
	DataCol *mongo.Collection
}

func (m *Mongodb) Get(ctx context.Context, token string, meta *Meta) error {
	// Query for document
	res := m.Col.FindOne(ctx, bson.M{"token": token})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return ErrBadToken
		}

		return res.Err()
	}

	// Decode to meta struct
	if err := res.Decode(meta); err != nil {
		return err
	}

	return nil
}

func (m *Mongodb) Write(ctx context.Context, meta *Meta, msg *Message) error {
	_, _ = m.Col.UpdateOne(ctx, bson.M{"_id": meta.ControllerId, "user_id": meta.UserId}, msg.Data)
	return nil
}
