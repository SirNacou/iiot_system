package iot

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

type InsertAlertsCommand struct {
	Time         time.Time
	DeviceID     string
	AlertType    string
	Severity     string
	Message      string
	CurrentValue null.Val[decimal.Decimal]
}

type InsertAlertsCommandHandler struct {
	db bob.DB
}

func NewIiotAlertsCommandHandler(db bob.DB) *InsertAlertsCommandHandler {
	return &InsertAlertsCommandHandler{
		db: db,
	}
}

func (h InsertAlertsCommandHandler) Handle(ctx context.Context, command ...InsertAlertsCommand) error {
	if len(command) == 0 {
		return nil
	}

	t, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := psql.Insert(
		im.Into(models.IotAlertEvents.Name()),
	)

	for _, v := range command {
		value := im.Values(
			psql.Quote(v.Time.String()),
			psql.Quote(v.DeviceID),
			psql.Quote(v.AlertType),
			psql.Quote(v.Severity),
			psql.Quote(v.Message),
			psql.Quote(v.CurrentValue.GetOrZero().String()),
		)

		q.Apply(value)
	}

	s, _, err := q.Build(ctx)
	if err != nil {
		return err
	}
	_, err = t.ExecContext(ctx, s)

	return err
}
