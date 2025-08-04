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
		Name:       data.Name,
		EbmsID:     pgtype.Int8{Int64: data.EbmsID, Valid: true},
		IsOutbound: data.IsOutbound,
		RouteID:    routeId.ID,
	})

	if err != nil {
		return nil, err
	}

	// todo: replace in redis

	return &result, nil
}
