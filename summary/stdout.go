package summary

import (
	"context"
	"fmt"
)

type Stdout struct{}

func (s *Stdout) Write(ctx context.Context, summary map[string]*Summary) error {
	for k, v := range summary {
		fmt.Printf("---- Controller ID:%v\n%v\n\n", k, v)
	}

	return nil
}
