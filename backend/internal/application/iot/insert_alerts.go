package application_iot

import (
	"context"
	"time"

	"iiot_system/backend/gen/models"

	"github.com/aarondl/opt/null"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/im"
)

type InsertAlertsCommand struct {
	Time         time.Time
	DeviceID     string
	AlertType    string
	Severity     string
	Message      string
	CurrentValue null.Val[float64]
}

type InsertAlertsCommandHandler struct {
	db bob.DB
}

func NewInsertAlertsCommandHandler(db bob.DB) *InsertAlertsCommandHandler {
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
		val, exist := v.CurrentValue.Get()

		var value bob.Mod[*dialect.InsertQuery]
		if exist {
			value = im.Values(
				psql.Arg(v.Time,
					v.DeviceID,
					v.AlertType,
					v.Severity,
					v.Message,
					val),
			)
		} else {
			value = im.Values(
				psql.Arg(v.Time,
					v.DeviceID,
					v.AlertType,
					v.Severity,
					v.Message,
					nil),
			)
		}

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
