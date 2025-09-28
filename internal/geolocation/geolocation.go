package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/catouberos/transit-radar/dto"
	"github.com/catouberos/transit-radar/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

func (a *App) CreateGeolocationByRouteIDAndPlateAndBound(ctx context.Context, data *dto.GeolocationByRouteIDAndPlateAndBoundInsert) (*models.Geolocation, error) {
	route, err := a.Query().GetRouteByEbmsID(ctx, pgtype.Int8{Int64: data.RouteID, Valid: true})
	if err != nil {
		return nil, err
	}

	variant, err := a.Query().GetVariantByRouteIDAndOutbound(ctx, models.GetVariantByRouteIDAndOutboundParams{
		RouteID:    route.ID,
		IsOutbound: data.IsOutbound,
	})
	if err != nil {
		return nil, err
	}

	vehicle, err := a.Query().GetVehicleByLicensePlate(ctx, data.LicensePlate)
	if err != nil {
		vehicle, err = a.Query().CreateVehicle(ctx, data.LicensePlate)
		if err != nil {
			return nil, err
		}
	}

	result, err := a.Query().CreateGeolocation(ctx, models.CreateGeolocationParams{
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

	cmd := a.Redis().GeoAdd(ctx, "geolocations", &redis.GeoLocation{
		Name:      fmt.Sprintf("geolocation:%d", vehicle.ID),
		Latitude:  float64(result.Latitude),
		Longitude: float64(result.Longitude),
	})
	if err := cmd.Err(); err != nil {
		slog.Warn("Cannot add geolocation to Redis", "error", err)
	}

	cmd = a.Redis().HSet(ctx, fmt.Sprintf("geolocation:%d", vehicle.ID), &dto.Geolocation{
		Degree:    result.Degree,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		Speed:     result.Speed,
		VehicleID: result.VehicleID,
		VariantID: result.VariantID,
		Timestamp: result.Timestamp.Time,
	})
	if err := cmd.Err(); err != nil {
		slog.Warn("Cannot add geolocation details to Redis", "error", err)
	}

	return &result, nil
}

func (a *App) ListGeolocationByBounding(ctx context.Context, lat, lng float32, width, height float32) ([]*dto.Geolocation, error) {
	results := []*dto.Geolocation{}

	cmd := a.Redis().GeoSearchLocation(ctx, "geolocations", &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Latitude:  float64(lat),
			Longitude: float64(lng),

			BoxWidth:  float64(width),
			BoxHeight: float64(height),
			BoxUnit:   "m",
		},
		WithCoord: true,
	})

	locations, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	for _, location := range locations {
		cmd := a.Redis().HGetAll(ctx, location.Name)
		if err := cmd.Err(); err != nil {
			slog.Warn("Location not in cache, getting from database...")
			// TODO: get location from database
			continue
		}

		result := &dto.Geolocation{}
		if err := cmd.Scan(result); err != nil {
			slog.Error("Cannot parse from Redis")
			// TODO: get from database also
			continue
		}

		results = aend(results, result)
	}

	return results, nil
}
