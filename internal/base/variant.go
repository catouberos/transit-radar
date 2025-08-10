package base

import (
	"context"

	"github.com/catouberos/transit-radar/dto"
	"github.com/catouberos/transit-radar/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func (app *App) CreateOrUpdateVariantByRouteEbmsID(ctx context.Context, data *dto.VariantByRouteEbmsIDUpsert) (*models.Variant, error) {
	routeId, err := app.GetRouteByEbmsID(ctx, data.RouteEbmsID)
	if err != nil {
		return nil, err
	}

	result, err := app.Query().CreateOrUpdateVariant(ctx, models.CreateOrUpdateVariantParams{
		Name:          data.Name,
		EbmsID:        pgtype.Int8{Int64: data.EbmsID, Valid: true},
		IsOutbound:    data.IsOutbound,
		RouteID:       routeId.ID,
		Description:   pgtype.Text{String: data.Description, Valid: true},
		ShortName:     pgtype.Text{String: data.ShortName, Valid: true},
		Distance:      pgtype.Float4{Float32: data.Distance, Valid: true},
		Duration:      pgtype.Int4{Int32: data.Duration, Valid: true},
		StartStopName: pgtype.Text{String: data.StartStopName, Valid: true},
		EndStopName:   pgtype.Text{String: data.EndStopName, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// todo: replace in redis

	return &result, nil
}

func (app *App) ListVariant(ctx context.Context) ([]models.Variant, error) {
	result, err := app.Query().ListVariant(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (app *App) GetVariant(ctx context.Context, id int64) (models.Variant, error) {
	result, err := app.Query().GetVariant(ctx, id)
	if err != nil {
		return models.Variant{}, err
	}

	return result, nil
}

func (app *App) ListVariantByRouteId(ctx context.Context, routeID int64) ([]models.Variant, error) {
	result, err := app.Query().ListVariantByRouteId(ctx, routeID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
