// Consumer reads data from a channel and writes to data source
// It defines interface for these operations and a struct that holds the channel
package consumer

import (
	"context"
	"github.com/tPhume/ags-pipeline/stream"
)

type Source interface {
	InsertData()
}

// InsertData interface is used by Consumer to insert new data to data source
type InsertData interface {
	Insert(ctx context.Context, msg stream.Message) error
}

// Consumer holds a channel created by stream package
type Consumer struct {
	Stream <-chan stream.Message
	Source Source
}
