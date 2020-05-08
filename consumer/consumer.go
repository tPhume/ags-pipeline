// Consumer reads data from a channel and passes it to an interface that connects to the data source
package consumer

import (
	"context"
	"errors"
	"log"
	"time"
)

var (
	ErrNilStream = errors.New("stream is nil")

	// If DataSource returns this error Consumer will exit
	ErrSourceFatal = errors.New("fatal error from DataSource")
)

const Fatal = iota

// DataSource is an interface that Consumer uses when it receives something
type DataSource interface {
	Do(ctx context.Context, msg interface{}) error
}

// Consumer gets consumes message from a channel
type Consumer struct {
	Stream     <-chan interface{}
	DataSource DataSource
}

// Listen starts a blocking call that waits for messages
func (c *Consumer) Listen() error {
	// context for this function
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	quit := make(chan int)

	if c.Stream == nil {
		return ErrNilStream
	}

	log.Println("Waiting for messages")
	go func() {
		for msg := range c.Stream {
			log.Println("Received a message")
			ctx, cancel := context.WithTimeout(mainCtx, time.Second*25)
			go func() {
				defer cancel()

				if err := c.DataSource.Do(ctx, msg); err != nil {
					quit <- Fatal
				}
			}()
		}
	}()

	// wait for exit
	switch <-quit {
	case Fatal:
		return ErrSourceFatal
	}

	return nil
}
