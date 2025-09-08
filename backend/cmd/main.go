package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	iot_application "iiot_system/backend/internal/application/iot"
	"iiot_system/backend/internal/infrastructure/configs"
	"iiot_system/backend/internal/presentation/iot"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	telemetryTopic    = "iiot.telemetry"
	productionTopic   = "iiot.production"
	statusUpdateTopic = "iiot.status_update"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := configs.LoadConfig()

	var db bob.DB
	{
		dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
			os.Exit(1)
		}
		defer dbpool.Close()
		db = bob.NewDB(stdlib.OpenDBFromPool(dbpool))
	}
	defer db.Close()

	kafka, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.KafkaBroker),
		kgo.ConsumerGroup(cfg.KafkaGroupID),
		kgo.ConsumeTopics(iot.AlertsTopic, productionTopic, statusUpdateTopic, telemetryTopic),
		kgo.AutoCommitMarks(),
		kgo.AutoCommitInterval(5*time.Second),
		kgo.OnPartitionsRevoked(func(ctx context.Context, c *kgo.Client, m map[string][]int32) {
			if err := c.CommitMarkedOffsets(ctx); err != nil {
				fmt.Printf("Failed to commit offsets: %v\n", err)
			}
		}),
		kgo.BlockRebalanceOnPoll(),
	)
	if err != nil {
		fmt.Printf("Unable to create kafka client: %v\n", err)
		os.Exit(1)
	}
	defer kafka.Close()

	insertAlertsHandler := iot_application.NewIiotAlertsCommandHandler(db)
	iiotAlertsConsumer, err := iot.NewIiotAlertConsumer(insertAlertsHandler, kafka)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create IIoT alerts consumer: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go iiotAlertsConsumer.Start(ctx, &wg)

	<-ctx.Done()
	wg.Wait()
}
