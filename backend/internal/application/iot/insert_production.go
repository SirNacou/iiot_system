package iot

import (
	"context"
	"fmt"
	"time"

	"iiot_system/backend/gen/models"

	"github.com/aarondl/opt/null"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/im"
)

type InsertProductionCommand struct {
	Time         time.Time
	DeviceID     string
	AlertType    string
	Severity     string
	Message      string
	CurrentValue null.Val[float64]
}

type InsertProductionCommandHandler struct {
	db bob.DB
}

func NewIiotProductionCommandHandler(db bob.DB) *InsertProductionCommandHandler {
	return &InsertProductionCommandHandler{
		db: db,
	}
}

func (h InsertProductionCommandHandler) Handle(ctx context.Context, command ...InsertProductionCommand) error {
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

	r, err := t.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	t.Commit(ctx)

	row, err := r.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("Executed query: %s with args: %v, row: %v\n", query, args, row)

	return err
}
