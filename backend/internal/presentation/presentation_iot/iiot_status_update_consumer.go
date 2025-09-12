package presentation_iot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"iiot_system/backend/internal/application/iot"
	"iiot_system/backend/internal/infrastructure/topics"

	"github.com/twmb/franz-go/pkg/kgo"
)

// AlertPayload represents the JSON structure of the `payload` field.
type StatusUpdatePayload struct {
	Timestamp int64  `json:"timestamp"`
	DeviceID  string `json:"device_id"`
	Data      struct {
		OldStatus string `json:"old_status"`
		NewStatus string `json:"new_status"`
		Reason    string `json:"reason"`
	} `json:"data"`
}

type IiotStatusUpdateConsumer struct {
	client  *kgo.Client
	handler *application_iot.InsertStatusUpdateCommandHandler
}

func NewIiotStatusUpdateConsumer(handler *application_iot.InsertStatusUpdateCommandHandler, client *kgo.Client) (*IiotStatusUpdateConsumer, error) {
	c := &IiotStatusUpdateConsumer{
		client:  client,
		handler: handler,
	}
	return c, nil
}

func (c IiotStatusUpdateConsumer) Start(ctx context.Context) error {
	if err := c.client.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping kafka: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("shutting down iiot status update consumer")
		default:
			fmt.Printf("polling iiot status update consumer\n")
			c.poll(ctx)
		}
	}
}

func (c IiotStatusUpdateConsumer) poll(ctx context.Context) {
	fetches := c.client.PollFetches(ctx)
	if fetches.IsClientClosed() {
		return
	}

	if fetches.Empty() {
		return
	}

	fetches.EachError(func(_ string, _ int32, err error) {
		panic(err)
	})

	r := fetches.Records()

	commands := make([]application_iot.InsertStatusUpdateCommand, 0, len(r))
	for _, rec := range r {
		if rec.Topic == topics.StatusUpdateTopic {
			var outer outerData
			if err := json.Unmarshal(rec.Value, &outer); err != nil {
				fmt.Printf("error unmarshaling outer JSON: %v\n", err)
				continue
			}

			var statusUpdate StatusUpdatePayload
			if err := json.Unmarshal([]byte(outer.Payload), &statusUpdate); err != nil {
				fmt.Printf("error unmarshaling status update JSON: %v\n", err)
				continue
			}

			commands = append(commands, application_iot.InsertStatusUpdateCommand{
				Time:      time.Unix(statusUpdate.Timestamp, 0).UTC(),
				DeviceID:  statusUpdate.DeviceID, // Replace with actual device ID from rec.Key or rec.Value
				OldStatus: statusUpdate.Data.OldStatus,
				NewStatus: statusUpdate.Data.NewStatus,
				Reason:    statusUpdate.Data.Reason,
			})
		}
	}
	if err := c.handler.Handle(ctx, commands...); err != nil {
		fmt.Printf("error handling insert status update command: %v\n", err)
		return
	}

	c.client.MarkCommitRecords(r...)
}
