package base

import (
	"context"

	"github.com/catouberos/geoloc/dto"
	"github.com/catouberos/geoloc/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func (app *App) CreateGeolocationByRouteIDAndPlateAndBound(ctx context.Context, data *dto.GeolocationByRouteIDAndPlateAndBoundInsert) (*models.Geolocation, error) {
	route, err := app.Query().GetRouteByEbmsID(ctx, pgtype.Int8{Int64: data.RouteID, Valid: true})
	if err != nil {
		return nil, err
	}

	variant, err := app.Query().GetVariantByRouteIDAndOutbound(ctx, models.GetVariantByRouteIDAndOutboundParams{
		RouteID:    route.ID,
		IsOutbound: data.IsOutbound,
	})
	if err != nil {
		return nil, err
	}

	vehicle, err := app.Query().GetVehicleByLicensePlate(ctx, data.LicensePlate)
	if err != nil {
		vehicle, err = app.Query().CreateVehicle(ctx, data.LicensePlate)
		if err != nil {
			return nil, err
		}
	}

	result, err := app.Query().CreateGeolocation(ctx, models.CreateGeolocationParams{
		Degree:    data.Degree,
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
		Speed:     data.Speed,
		VehicleID: vehicle.ID,
		VariantID: variant.ID,
		Timestamp: pgtype.Timestamptz{Time: data.Timestamp, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	// todo: replace in redis

	return &result, nil
}
