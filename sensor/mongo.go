package sensor

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongodb struct {
	Col *mongo.Collection
}

func (m *Mongodb) Get(ctx context.Context, token string, meta *Meta) error {
	panic("implement me")
}
