package iot

import (
	"context"
	"time"

	"iiot_system/backend/gen/models"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
)

type InsertProductionCommand struct {
	Time           time.Time
	DeviceID       string
	ProductionType string
	ProductSku     string
	UnitCount      int32
	BatchID        string
	QualityStatus  string
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
		value := im.Values(
			psql.Arg(v.Time,
				v.DeviceID,
				v.ProductionType,
				v.ProductSku,
				v.UnitCount,
				v.BatchID,
				v.QualityStatus),
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

	return t.Commit(ctx)
}
