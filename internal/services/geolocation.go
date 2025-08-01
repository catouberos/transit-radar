package services

import (
	"context"

	"github.com/catouberos/geoloc/base"
	"github.com/catouberos/geoloc/internal/events"
	"github.com/catouberos/geoloc/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func NewGeolocation(ctx context.Context, app *base.App, data *events.GeolocationInsert) (*models.Geolocation, error) {
	result, err := app.Queries.CreateGeolocation(ctx, models.CreateGeolocationParams{
		Degree:    data.Degree,
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
		Speed:     data.Speed,
		VehicleID: data.VehicleID,
		RouteID:   data.RouteID,
		Timestamp: pgtype.Timestamptz{Time: data.Timestamp, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	// todo: replace in redis

	return &result, nil
}
