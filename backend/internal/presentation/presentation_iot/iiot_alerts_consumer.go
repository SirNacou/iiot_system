package presentation_iot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"iiot_system/backend/internal/application/iot"
	"iiot_system/backend/internal/infrastructure/topics"

	"github.com/aarondl/opt/null"
	"github.com/twmb/franz-go/pkg/kgo"
)

type outerData struct {
	Payload string `json:"payload"`
}

// AlertPayload represents the JSON structure of the `payload` field.
type AlertPayload struct {
	Timestamp   int64  `json:"timestamp"`
	DeviceID    string `json:"device_id"`
	PayloadType string `json:"payload_type"`
	Data        struct {
		AlertType    string  `json:"alert_type"`
		Severity     string  `json:"severity"`
		Message      string  `json:"message"`
		CurrentValue float64 `json:"current_value"`
	} `json:"data"`
}

type IiotAlertsConsumer struct {
	client  *kgo.Client
	handler *application_iot.InsertAlertsCommandHandler
}

func NewIiotAlertConsumer(handler *application_iot.InsertAlertsCommandHandler, client *kgo.Client) (*IiotAlertsConsumer, error) {
	c := &IiotAlertsConsumer{
		handler: handler,
		client:  client,
	}

	return c, nil
}

func (c IiotAlertsConsumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("shutting down iiot alerts consumer")
		default:
			fmt.Printf("polling iiot alerts consumer\n")
			c.poll(ctx)
		}
	}
}

func (c IiotAlertsConsumer) poll(ctx context.Context) {
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

	var recordsToCommit []*kgo.Record
	commands := make([]application_iot.InsertAlertsCommand, 0, len(r))
	for _, rec := range r {
		if rec.Topic == topics.AlertsTopic {
			var outer outerData
			if err := json.Unmarshal(rec.Value, &outer); err != nil {
				fmt.Printf("error unmarshaling outer JSON: %v\n", err)
				continue
			}

			var alert AlertPayload
			if err := json.Unmarshal([]byte(outer.Payload), &alert); err != nil {
				fmt.Printf("error unmarshaling alert JSON: %v\n", err)
				continue
			}

			commands = append(commands, application_iot.InsertAlertsCommand{
				Time:         time.Unix(alert.Timestamp, 0).UTC(),
				DeviceID:     alert.DeviceID,       // Replace with actual device ID from rec.Key or rec.Value
				AlertType:    alert.Data.AlertType, // Replace with actual alert type from rec.Value
				Severity:     alert.Data.Severity,  // Replace with actual severity from rec.Value
				Message:      alert.Data.Message,
				CurrentValue: null.From(alert.Data.CurrentValue),
			})
			recordsToCommit = append(recordsToCommit, rec)
		}
	}
	if err := c.handler.Handle(ctx, commands...); err != nil {
		fmt.Printf("error handling insert alerts command: %v\n", err)
		return
	}

	c.client.MarkCommitRecords(recordsToCommit...)
}
