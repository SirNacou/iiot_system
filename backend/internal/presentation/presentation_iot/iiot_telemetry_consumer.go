package presentation_iot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	application_iot "iiot_system/backend/internal/application/iot"
	"iiot_system/backend/internal/infrastructure/topics"

	"github.com/aarondl/opt/null"
	"github.com/shopspring/decimal"
	"github.com/twmb/franz-go/pkg/kgo"
)

// AlertPayload represents the JSON structure of the `payload` field.
type TelemetryPayload struct {
	Timestamp int64  `json:"timestamp"`
	DeviceID  string `db:"device_id" `
	Data      struct {
		TemperatureCelcius decimal.Decimal  `json:"temperature_celcius" `
		HumidityPercent    decimal.Decimal  `json:"humidity_percent" `
		VibrationHZ        decimal.Decimal  `json:"vibration_hz" `
		MotorRPM           int32            `json:"motor_rpm" `
		CurrentAmps        decimal.Decimal  `json:"current_amps" `
		MachineStatus      string           `json:"machine_status" `
		ErrorCode          null.Val[string] `json:"error_code" `
	} `json:"data"`
}

type IiotTelemetryConsumer struct {
	client  *kgo.Client
	handler *application_iot.InsertTelemetryCommandHandler
}

func NewIiotTelemetryConsumer(handler *application_iot.InsertTelemetryCommandHandler, client *kgo.Client) (*IiotTelemetryConsumer, error) {
	c := &IiotTelemetryConsumer{
		client:  client,
		handler: handler,
	}
	return c, nil
}

func (c IiotTelemetryConsumer) Start(ctx context.Context) error {
	if err := c.client.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping kafka: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("shutting down iiot telemetry consumer")
		default:
			fmt.Printf("polling iiot telemetry consumer\n")
			c.poll(ctx)
		}
	}
}

func (c IiotTelemetryConsumer) poll(ctx context.Context) {
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

	commands := make([]application_iot.InsertTelemetryCommand, 0, len(r))
	for _, rec := range r {
		if rec.Topic == topics.TelemetryTopic {
			var outer outerData
			if err := json.Unmarshal(rec.Value, &outer); err != nil {
				fmt.Printf("error unmarshaling outer JSON: %v\n", err)
				continue
			}

			var telemetry TelemetryPayload
			if err := json.Unmarshal([]byte(outer.Payload), &telemetry); err != nil {
				fmt.Printf("error unmarshaling telemetry JSON: %v\n", err)
				continue
			}

			commands = append(commands, application_iot.InsertTelemetryCommand{
				Time:               time.Unix(telemetry.Timestamp, 0).UTC(),
				DeviceID:           telemetry.DeviceID, // Replace with actual device ID from rec.Key or rec.Value
				TemperatureCelcius: telemetry.Data.TemperatureCelcius,
				HumidityPercent:    telemetry.Data.HumidityPercent,
				VibrationHZ:        telemetry.Data.VibrationHZ,
				MotorRPM:           telemetry.Data.MotorRPM,
				CurrentAmps:        telemetry.Data.CurrentAmps,
				MachineStatus:      telemetry.Data.MachineStatus,
				ErrorCode:          telemetry.Data.ErrorCode,
			})
		}
	}
	if err := c.handler.Handle(ctx, commands...); err != nil {
		fmt.Printf("error handling insert telemetry command: %v\n", err)
		return
	}

	c.client.MarkCommitRecords(r...)
}
