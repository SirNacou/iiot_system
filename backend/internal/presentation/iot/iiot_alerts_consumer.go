package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"iiot_system/backend/internal/application/iot"

	"github.com/aarondl/opt/null"
	"github.com/twmb/franz-go/pkg/kgo"
)

type OuterData struct {
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
	handler *iot.InsertAlertsCommandHandler
}

func NewIiotAlertConsumer(handler *iot.InsertAlertsCommandHandler, seeds []string, group string, topics ...string) (*IiotAlertsConsumer, error) {
	c := &IiotAlertsConsumer{
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

func (c IiotAlertsConsumer) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := c.client.Ping(ctx); err != nil {
		fmt.Printf("unable to ping kafka: %v\n", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("shutting down iiot alerts consumer")
			return
		default:
			fmt.Printf("polling iiot alerts consumer\n")
			c.poll(ctx)
		}
	}
}

func (c IiotAlertsConsumer) Close() error {
	c.client.Close()
	return nil
}

func (c IiotAlertsConsumer) revoked(ctx context.Context, cl *kgo.Client, _ map[string][]int32) {
	if err := cl.CommitMarkedOffsets(ctx); err != nil {
		fmt.Printf("Failed to commit offsets: %v\n", err)
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

	commands := make([]iot.InsertAlertsCommand, 0, len(r))
	for _, rec := range r {
		var outer OuterData
		if err := json.Unmarshal(rec.Value, &outer); err != nil {
			fmt.Printf("error unmarshaling outer JSON: %v\n", err)
			continue
		}
		
		var alert AlertPayload
		if err := json.Unmarshal([]byte(outer.Payload), &alert); err != nil {
			fmt.Printf("error unmarshaling alert JSON: %v\n", err)
			continue
		}

		commands = append(commands, iot.InsertAlertsCommand{
			Time:      time.Unix(alert.Timestamp, 0).UTC(),
			DeviceID:  alert.DeviceID,               // Replace with actual device ID from rec.Key or rec.Value
			AlertType: alert.Data.AlertType, // Replace with actual alert type from rec.Value
			Severity:  alert.Data.Severity,                // Replace with actual severity from rec.Value
			Message:   alert.Data.Message,
			CurrentValue: null.From(alert.Data.CurrentValue),
		})
	}
	if err := c.handler.Handle(ctx, commands...); err != nil {
		fmt.Printf("error handling insert alerts command: %v\n", err)
		return
	}

	c.client.MarkCommitRecords(r...)
}
