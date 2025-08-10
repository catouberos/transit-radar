package base

import (
	"context"
	"log/slog"

	"github.com/catouberos/transit-radar/dto"
	"github.com/catouberos/transit-radar/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

// Import should be called on one-to-many relationship between `variant` and `stop`, for example import/replace/update all stops of one variant
func (app *App) ImportVariantStops(ctx context.Context, data *[]dto.VariantStopByEbmsIDImport) error {
	tx, err := app.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := app.Query().WithTx(tx)

	for _, record := range *data {
		variant, err := qtx.GetVariantByRouteEbmsID(ctx,
			models.GetVariantByRouteEbmsIDParams{
				EbmsID:   pgtype.Int8{Int64: record.VariantEbmsID, Valid: true},
				EbmsID_2: pgtype.Int8{Int64: record.RouteEbmsID, Valid: true},
			},
		)
		if err != nil {
			slog.Error("Variant not found", "record", record)
			return err
		}

		stop, err := qtx.GetStopByEbmsID(ctx, pgtype.Int8{Int64: record.StopEbmsID, Valid: true})
		if err != nil {
			slog.Error("Stop not found", "record", record)
			return err
		}

		if _, err = qtx.CreateVariantStop(ctx, models.CreateVariantStopParams{
			VariantID:  variant.ID,
			StopID:     stop.ID,
			OrderScore: int32(record.OrderScore),
		}); err != nil {
			// ignore duplicate
			slog.Error("Error creating variant stop", "record", record)
			continue
		}
	}

	// TODO: redis

	return tx.Commit(ctx)
}
