package base

import (
	"context"

	"github.com/catouberos/transit-radar/dto"
	"github.com/catouberos/transit-radar/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func (app *App) CreateOrUpdateRoute(ctx context.Context, data *dto.RouteUpsert) (*models.Route, error) {
	result, err := app.Query().CreateOrUpdateRoute(ctx, models.CreateOrUpdateRouteParams{
		Number:        data.Number,
		Name:          data.Name,
		EbmsID:        pgtype.Int8{Int64: data.EbmsID, Valid: true},
		OperationTime: pgtype.Text{String: data.OperationTime, Valid: true},
		Organization:  pgtype.Text{String: data.Organization, Valid: true},
		Ticketing:     pgtype.Text{String: data.Ticketing, Valid: true},
		RouteType:     pgtype.Text{String: data.RouteType, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// todo: replace in redis

	return &result, nil
}

func (app *App) GetRouteByEbmsID(ctx context.Context, ebmsID int64) (*models.Route, error) {
	result, err := app.Query().GetRouteByEbmsID(ctx, pgtype.Int8{Int64: ebmsID, Valid: true})

	if err != nil {
		return nil, err
	}

	// todo: replace in redis

	return &result, nil

}
