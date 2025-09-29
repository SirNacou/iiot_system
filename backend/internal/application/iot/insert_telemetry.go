package application_iot

import (
	"context"
	"time"

	"iiot_system/backend/gen/models"

	"github.com/aarondl/opt/null"
	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
)

type InsertTelemetryCommand struct {
	Time               time.Time
	DeviceID           string
	TemperatureCelcius decimal.Decimal
	HumidityPercent    decimal.Decimal
	VibrationHZ        decimal.Decimal
	MotorRPM           int32
	CurrentAmps        decimal.Decimal
	MachineStatus      string
	ErrorCode          null.Val[string]
}

type InsertTelemetryCommandHandler struct {
	db bob.DB
}

func NewInsertTelemetryCommandHandler(db bob.DB) *InsertTelemetryCommandHandler {
	return &InsertTelemetryCommandHandler{
		db: db,
	}
}

func (h InsertTelemetryCommandHandler) Handle(ctx context.Context, command ...InsertTelemetryCommand) error {
	if len(command) == 0 {
		return nil
	}

	t, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := psql.Insert(
		im.Into(models.IotTelemetryEvents.Name()),
	)

	for _, v := range command {
		value := im.Values(
			psql.Arg(v.Time,
				v.DeviceID,
				v.TemperatureCelcius,
				v.HumidityPercent,
				v.VibrationHZ,
				v.MotorRPM,
				v.CurrentAmps,
				v.MachineStatus,
				v.ErrorCode,
			),
		)

		q.Apply(value)
	}

	query, args, err := q.Build(ctx)
	if err != nil {
		return err
	}

	_, err = t.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	err = t.Commit(ctx)
	if err != nil {
		return err
	}

	return err
}
