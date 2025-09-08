package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"iiot_system/backend/internal/application/iot"

	"github.com/twmb/franz-go/pkg/kgo"
)

// AlertPayload represents the JSON structure of the `payload` field.
type ProductionPayload struct {
	Timestamp int64  `json:"timestamp"`
	DeviceID  string `db:"device_id" `
	Data      struct {
		ProductionType string `db:"production_type" `
		ProductSku     string `db:"product_sku" `
		UnitCount      int32  `db:"unit_count" `
		BatchID        string `db:"batch_id" `
		QualityStatus  string `db:"quality_status" `
	} `json:"data"`
}

type IiotProductionConsumer struct {
	client  *kgo.Client
	handler *iot.InsertProductionCommandHandler
}

func NewIiotProductionConsumer(handler *iot.InsertProductionCommandHandler, seeds []string, group string, topics ...string) (*IiotProductionConsumer, error) {
	c := &IiotProductionConsumer{
		handler: handler,
	}
	var err error
	c.client, err = kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topics...),
		kgo.AutoCommitMarks(),
		kgo.AutoCommitInterval(5*time.Second),
		kgo.OnPartitionsRevoked(c.revoked),
		kgo.BlockRebalanceOnPoll(),
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c IiotProductionConsumer) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := c.client.Ping(ctx); err != nil {
		fmt.Printf("unable to ping kafka: %v\n", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("shutting down iiot production consumer")
			return
		default:
			fmt.Printf("polling iiot production consumer\n")
			c.poll(ctx)
		}
	}
}

func (c IiotProductionConsumer) Close() error {
	c.client.Close()
	return nil
}

func (c IiotProductionConsumer) revoked(ctx context.Context, cl *kgo.Client, _ map[string][]int32) {
	if err := cl.CommitMarkedOffsets(ctx); err != nil {
		fmt.Printf("Failed to commit offsets: %v\n", err)
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

	commands := make([]iot.InsertProductionCommand, 0, len(r))
	for _, rec := range r {
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

		commands = append(commands, iot.InsertProductionCommand{
			Time:         time.Unix(production.Timestamp, 0).UTC(),
			DeviceID:     production.DeviceID,       // Replace with actual device ID from rec.Key or rec.Value
			ProductionType:    production.Data.ProductionType, // Replace with actual alert type from rec.Value
			ProductSku:     production.Data.ProductSku,  // Replace with actual severity from rec.Value
			UnitCount:      production.Data.UnitCount,
			BatchID: production.Data.BatchID,
			QualityStatus: production.Data.QualityStatus,
		})
	}
	if err := c.handler.Handle(ctx, commands...); err != nil {
		fmt.Printf("error handling insert alerts command: %v\n", err)
		return
	}

	c.client.MarkCommitRecords(r...)
}
