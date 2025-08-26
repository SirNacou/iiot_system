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

	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
)

var (
	logger         = watermill.NewStdLogger(false, false)
	alertsTopic    = "iiot.alerts"
	telemetryTopic = "iiot.telemetry"
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

	insertAlertsHandler := iot_application.NewIiotAlertsCommandHandler(db)

	seeds := []string{cfg.KafkaBroker}
	var wg sync.WaitGroup

	iiotAlertsConsumer, err := iot.NewIiotAlertConsumer(insertAlertsHandler, seeds, cfg.KafkaGroupID, alertsTopic)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create IIoT alerts consumer: %v\n", err)
		os.Exit(1)
	}
	defer iiotAlertsConsumer.Close()

	wg.Add(1)
	go iiotAlertsConsumer.Start(ctx, &wg)

	<-ctx.Done()
	wg.Wait()

	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               []string{cfg.KafkaBroker},
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         cfg.KafkaGroupID,
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddPlugin(plugin.SignalsHandler)

	router.AddMiddleware(
		middleware.CorrelationID,

		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          router.Logger(),
		}.Middleware,

		middleware.Recoverer,
	)

	OeeAlertsConsumer := iot.NewOeeAlertConsumer(logger)
	router.AddNoPublisherHandler(
		"save_oee_alerts",
		"iiot.alerts",
		subscriber,
		OeeAlertsConsumer.Handle,
	)

	if err := router.Run(ctx); err != nil {
		panic(err)
	}
}
