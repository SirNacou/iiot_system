package presentation_iot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	application_iot "iiot_system/backend/internal/application/iot"
	"iiot_system/backend/internal/infrastructure/topics"

	"github.com/twmb/franz-go/pkg/kgo"
)

// AlertPayload represents the JSON structure of the `payload` field.
type ProductionPayload struct {
	Timestamp int64  `json:"timestamp"`
	DeviceID  string `json:"device_id" `
	Data      struct {
		ProductionType string `json:"production_type" `
		ProductSku     string `json:"product_sku" `
		UnitCount      int32  `json:"unit_count" `
		BatchID        string `json:"batch_id" `
		QualityStatus  string `json:"quality_status" `
	} `json:"data"`
}

type IiotProductionConsumer struct {
	client  *kgo.Client
	handler *application_iot.InsertProductionCommandHandler
}

func NewIiotProductionConsumer(handler *application_iot.InsertProductionCommandHandler, client *kgo.Client) (*IiotProductionConsumer, error) {
	c := &IiotProductionConsumer{
		client:  client,
		handler: handler,
	}
	return c, nil
}

func (c IiotProductionConsumer) Start(ctx context.Context) error {
	if err := c.client.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping kafka: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("shutting down iiot production consumer")
		default:
			fmt.Printf("polling iiot production consumer\n")
			c.poll(ctx)
		}
	}
}

func (c IiotProductionConsumer) poll(ctx context.Context) {
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

	commands := make([]application_iot.InsertProductionCommand, 0, len(r))
	for _, rec := range r {
		if rec.Topic == topics.ProductionTopic {
			var outer outerData
			if err := json.Unmarshal(rec.Value, &outer); err != nil {
				fmt.Printf("error unmarshaling outer JSON: %v\n", err)
				continue
			}

			var production ProductionPayload
			if err := json.Unmarshal([]byte(outer.Payload), &production); err != nil {
				fmt.Printf("error unmarshaling alert JSON: %v\n", err)
				continue
			}

			commands = append(commands, application_iot.InsertProductionCommand{
				Time:           time.Unix(production.Timestamp, 0).UTC(),
				DeviceID:       production.DeviceID,            // Replace with actual device ID from rec.Key or rec.Value
				ProductionType: production.Data.ProductionType, // Replace with actual alert type from rec.Value
				ProductSku:     production.Data.ProductSku,     // Replace with actual severity from rec.Value
				UnitCount:      production.Data.UnitCount,
				BatchID:        production.Data.BatchID,
				QualityStatus:  production.Data.QualityStatus,
			})
		}
	}
	if err := c.handler.Handle(ctx, commands...); err != nil {
		fmt.Printf("error handling insert telemetry command: %v\n", err)
		return
	}

	c.client.MarkCommitRecords(r...)
}
