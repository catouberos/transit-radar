package geolocation

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

const (
	cacheKey = "geolocation:vehicle:%d"
)

var _ GeolocationCache = (*redisCache)(nil)

// using VehicleID as cache
type redisCache struct {
	client *redis.Client
}

func (r *redisCache) Put(ctx context.Context, geolocation Geolocation) error {
	cmd := r.client.HSet(ctx, fmt.Sprintf(cacheKey, geolocation.VehicleID), geolocation)
	if err := cmd.Err(); err != nil {
		return err
	}

	cmd = r.client.GeoAdd(ctx, "geolocations", &redis.GeoLocation{
		Name:      fmt.Sprintf(cacheKey, geolocation.VehicleID),
		Latitude:  geolocation.Location.Latitude,
		Longitude: geolocation.Location.Longitude,
	})
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *redisCache) Get(ctx context.Context, vehicleID int64) (Geolocation, error) {
	return r.get(ctx, fmt.Sprintf(cacheKey, vehicleID))
}

func (r *redisCache) ListByBounding(ctx context.Context, params ListByBoundingParams) ([]Geolocation, error) {
	cmd := r.client.GeoSearchLocation(ctx, "geolocations", &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Latitude:  params.Latitude,
			Longitude: params.Longitude,

			BoxWidth:  params.Width,
			BoxHeight: params.Height,
			BoxUnit:   params.Unit,
		},
		WithCoord: true,
	})

	locations, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	geolocations := make([]Geolocation, len(locations))

	for i, location := range locations {
		values := strings.Split(location.Name, ":")
		rawVehicleID := values[len(values)-1]
		vehicleID, err := strconv.ParseInt(rawVehicleID, 10, 64)
		if err != nil {
			// TODO: log
			continue
		}

		geolocation, err := s.Get(ctx, GetParams{
			VehicleID: vehicleID,
		})
		if err != nil {
			// TODO: log
			continue
		}

		geolocations[i] = geolocation
	}

	return geolocations, nil
}

func (r *redisCache) get(ctx context.Context, key string) (Geolocation, error) {
	cmd := r.client.HGetAll(ctx, key)
	if err := cmd.Err(); err != nil {
		return Geolocation{}, err
	}

	var geolocation Geolocation
	if err := cmd.Scan(&geolocation); err != nil {
		return Geolocation{}, err
	}

	return geolocation, nil
}
