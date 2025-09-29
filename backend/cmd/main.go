package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	application_iot "iiot_system/backend/internal/application/iot"
	"iiot_system/backend/internal/infrastructure/configs"
	"iiot_system/backend/internal/infrastructure/topics"
	"iiot_system/backend/internal/presentation/presentation_iot"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/stephenafamo/bob"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := configs.LoadConfig()

	var db bob.DB
	{
		dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Printf("Unable to create connection pool: %v\n", err)
			os.Exit(1)
		}
		defer dbpool.Close()
		db = bob.NewDB(stdlib.OpenDBFromPool(dbpool))
	}
	defer db.Close()

	kafka, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.KafkaBroker),
		kgo.ConsumerGroup(cfg.KafkaGroupID),
		kgo.ConsumeTopics(topics.AlertsTopic, topics.ProductionTopic, topics.StatusUpdateTopic, topics.TelemetryTopic),
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
		log.Fatalf("Unable to create kafka client: %v\n", err)
	}
	defer kafka.Close()

	if err := kafka.Ping(ctx); err != nil {
		log.Fatalf("unable to ping kafka: %v\n", err)
	}

	insertAlertsHandler := application_iot.NewInsertAlertsCommandHandler(db)
	insertProductionHandler := application_iot.NewInsertProductionCommandHandler(db)
	insertStatusUpdateHandler := application_iot.NewInsertStatusUpdateCommandHandler(db)
	insertTelemetryHandler := application_iot.NewInsertTelemetryCommandHandler(db)

	iiotAlertsConsumer, err := presentation_iot.NewIiotAlertConsumer(insertAlertsHandler, kafka)
	if err != nil {
		log.Fatalf("Unable to create IIoT alerts consumer: %v\n", err)
	}

	iiotProductionConsumer, err := presentation_iot.NewIiotProductionConsumer(insertProductionHandler, kafka)
	if err != nil {
		log.Fatalf("Unable to create IIoT production consumer: %v\n", err)
	}

	iiotStatusUpdateConsumer, err := presentation_iot.NewIiotStatusUpdateConsumer(insertStatusUpdateHandler, kafka)
	if err != nil {
		log.Fatalf("Unable to create IIoT status update consumer: %v\n", err)
	}

	iiotTelemetryConsumer, err := presentation_iot.NewIiotTelemetryConsumer(insertTelemetryHandler, kafka)
	if err != nil {
		log.Printf("Unable to create IIoT telemetry consumer: %v\n", err)
		os.Exit(1)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	var wg sync.WaitGroup

	wg.Go(func() {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port)))
	})

	wg.Go(func() {
		err := iiotAlertsConsumer.Start(ctx)
		if err != nil {
			log.Fatal("Iiot alert consumer stopped with error", err)
		}
	})
	wg.Go(func() {
		err := iiotProductionConsumer.Start(ctx)
		if err != nil {
			log.Fatal("Iiot production consumer stopped with error", err)
		}
	})
	wg.Go(func() {
		err := iiotStatusUpdateConsumer.Start(ctx)
		if err != nil {
			log.Fatal("Iiot status update consumer stopped with error", err)
		}
	})
	wg.Go(func() {
		err := iiotTelemetryConsumer.Start(ctx)
		if err != nil {
			log.Fatal("Iiot telemetry consumer stopped with error", err)
		}
	})

	log.Println("Application is running. Press Ctrl+C to stop.")
	<-ctx.Done()
	log.Println("Waiting for consumers to finish...")
	wg.Wait()
	log.Println("Application shutting down gracefully.")
}
