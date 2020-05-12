package sensor

import (
	"context"
	"encoding/json"
	"fmt"
)

type Stdout struct{}

func (s *Stdout) Write(ctx context.Context, meta *Meta, msg *Message) error {
	msgJson, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("user_id: %s\ncontroller_id: %s\ndata: %s", meta.UserId, meta.ControllerId, string(msgJson))
	return nil
}
