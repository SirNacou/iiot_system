package application_iot

import (
	"context"
	"time"

	"iiot_system/backend/gen/models"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
)

type InsertStatusUpdateCommand struct {
	Time      time.Time
	DeviceID  string
	OldStatus string
	NewStatus string
	Reason    string
}

type InsertStatusUpdateCommandHandler struct {
	db bob.DB
}

func NewInsertStatusUpdateCommandHandler(db bob.DB) *InsertStatusUpdateCommandHandler {
	return &InsertStatusUpdateCommandHandler{
		db: db,
	}
}

func (h InsertStatusUpdateCommandHandler) Handle(ctx context.Context, command ...InsertStatusUpdateCommand) error {
	if len(command) == 0 {
		return nil
	}

	t, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := psql.Insert(
		im.Into(models.IotStatusEvents.Name()),
	)

	for _, v := range command {
		value := im.Values(
			psql.Arg(v.Time,
				v.DeviceID,
				v.OldStatus,
				v.NewStatus,
				v.Reason,
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
