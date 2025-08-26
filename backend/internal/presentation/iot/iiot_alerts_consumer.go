package iot

import (
	"context"
	"fmt"
	"iiot_system/backend/internal/application/iot"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type IiotAlertsConsumer struct {
	client *kgo.Client
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
	
	for _, rec := range r {
		fmt.Printf("received message: topic=%s partition=%d offset=%d key=%s value=%s\n", rec.Topic, rec.Partition, rec.Offset, string(rec.Key), string(rec.Value))
	}
	c.handler.Handle(ctx, )

	c.client.MarkCommitRecords(r...)
}
