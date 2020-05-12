package sensor

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go"
)

type Influxdb struct {
	WriteApi influxdb2.WriteApiBlocking
}

func (i *Influxdb) Write(ctx context.Context, msg *Message) error {
	panic("implement me")
}
