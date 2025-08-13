package main

import (
	"context"
	"fmt"
	"os"
	"time"

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

var logger = watermill.NewStdLogger(false, false)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	cfg := configs.LoadConfig()

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	db := bob.NewDB(stdlib.OpenDBFromPool(dbpool))
	defer db.Close()

	// cl, err := kgo.NewClient()
	// cl.LeaveGroup()
	// defer cl.Close()

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
		"oee.alerts",
		subscriber,
		OeeAlertsConsumer.Handle,
	)

	if err := router.Run(ctx); err != nil {
		panic(err)
	}
}
