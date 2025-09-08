package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	iot_application "iiot_system/backend/internal/application/iot"
	"iiot_system/backend/internal/infrastructure/configs"
	"iiot_system/backend/internal/presentation/iot"

	"github.com/ThreeDotsLabs/watermill"
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

	seeds := []string{cfg.KafkaBroker}
	insertAlertsHandler := iot_application.NewIiotAlertsCommandHandler(db)
	iiotAlertsConsumer, err := iot.NewIiotAlertConsumer(insertAlertsHandler, seeds, cfg.KafkaGroupID, alertsTopic)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create IIoT alerts consumer: %v\n", err)
		os.Exit(1)
	}
	defer iiotAlertsConsumer.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go iiotAlertsConsumer.Start(ctx, &wg)

	<-ctx.Done()
	wg.Wait()
}
