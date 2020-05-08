// Consumer reads data from a channel and passes it to an interface that connects to the data source
package consumer

import (
	"context"
)

// DataSource is an interface that Consumer uses when it receives something
type DataSource interface {
	Do(ctx context.Context, msg *interface{}) error
}

// Consumer gets consumes message from a channel
type Consumer struct {
	Stream     <-chan interface{}
	DataSource DataSource
}
