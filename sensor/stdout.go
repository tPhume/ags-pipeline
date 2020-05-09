package sensor

import (
	"context"
	"encoding/json"
	"fmt"
)

type Stdout struct{}

func (s *Stdout) Write(ctx context.Context, msg *Message) error {
	msgJson, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(msgJson))
	return nil
}
