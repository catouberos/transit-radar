package base

import (
	"context"

	"github.com/catouberos/transit-radar/dto"
	"github.com/catouberos/transit-radar/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func (app *App) ImportStop(ctx context.Context, data *dto.StopImport) error {
	stopType, err := app.Query().GetStopTypeByName(ctx, data.TypeName)
	if err != nil {
		stopType, err = app.Query().CreateStopType(ctx, data.TypeName)
		if err != nil {
			return err
		}
	}

	err = app.Query().CreateOrUpdateStop(ctx, models.CreateOrUpdateStopParams{
		Code:      data.Code,
		Name:      data.Name,
		TypeID:    stopType.ID,
		EbmsID:    pgtype.Int8{Int64: data.EbmsID, Valid: true},
		Active:    data.Active,
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
	})
	if err != nil {
		return err
	}

	// TODO: redis

	return nil
}

func (app *App) GetStopByEbmsID(ctx context.Context, ebmsID int64) (*models.Stop, error) {
	stop, err := app.Query().GetStopByEbmsID(ctx, pgtype.Int8{Int64: ebmsID, Valid: true})
	if err != nil {
		return nil, err
	}

	return &stop, nil
}

func (app *App) ListStop(ctx context.Context) (*[]models.Stop, error) {
	stops, err := app.Query().ListStop(ctx)
	if err != nil {
		return nil, err
	}

	return &stops, nil
}
