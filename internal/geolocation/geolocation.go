package geolocation

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/catouberos/transit-radar/internal/models"
	"github.com/catouberos/transit-radar/types"
	"github.com/cridenour/go-postgis"
	"github.com/redis/go-redis/v9"
)

const (
	cacheKey = "geolocation:vehicle:%d"
)

type GeolocationService interface {
	Create(context.Context, CreateParams) (Geolocation, error)
	Get(context.Context, GetParams) (Geolocation, error)
	List(context.Context, ListParams) ([]Geolocation, error)
	ListByBounding(context.Context, ListByBoundingParams) ([]Geolocation, error)
}

type Geolocation struct {
	Degree    float32      `redis:"degree"`
	Location  types.LatLng `redis:"-"`
	Speed     float32      `redis:"speed"`
	VehicleID int64        `redis:"vehicle_id"`
	VariantID int64        `redis:"variant_id"`
	Timestamp time.Time    `redis:"timestamp"`
}

type CreateParams struct {
	Degree    float32
	Location  types.LatLng
	Speed     float32
	VehicleID int64
	VariantID int64
	Timestamp time.Time
}

type GetParams struct {
	VehicleID int64
}

type ListParams struct {
	Limit int32
}

type ListByBoundingParams struct {
	Latitude, Longitude, Width, Height float64
	Unit                               string
}

var _ GeolocationService = (*GeolocationServiceImpl)(nil)

func NewGeolocationService(query *models.Queries, redis *redis.Client) GeolocationService {
	return &GeolocationServiceImpl{
		query: query,
		redis: redis,
	}
}

type GeolocationServiceImpl struct {
	query *models.Queries
	redis *redis.Client
}

func (s *GeolocationServiceImpl) Create(ctx context.Context, params CreateParams) (Geolocation, error) {
	result, err := s.query.CreateGeolocation(ctx, models.CreateGeolocationParams{
		Degree: params.Degree,
		Location: postgis.Point{
			X: params.Location.Longitude,
			Y: params.Location.Latitude,
		},
		Speed:     params.Speed,
		VehicleID: params.VehicleID,
		VariantID: params.VariantID,
		Timestamp: params.Timestamp,
	})
	if err != nil {
		return Geolocation{}, err
	}
	geolocation := buildGeolocation(result)

	// TODO: error handling
	s.cachePut(ctx, geolocation)
	s.geoPut(ctx, geolocation)

	return geolocation, nil
}

func (s *GeolocationServiceImpl) Get(ctx context.Context, params GetParams) (Geolocation, error) {
	cached, err := s.cacheGet(ctx, params.VehicleID)
	if err == nil {
		return cached, nil
	}

	result, err := s.query.GetGeolocation(ctx, &params.VehicleID)
	if err != nil {
		return Geolocation{}, err
	}
	geolocation := buildGeolocation(result)

	return geolocation, nil
}

func (s *GeolocationServiceImpl) List(ctx context.Context, params ListParams) ([]Geolocation, error) {
	result, err := s.query.ListGeolocation(ctx, models.ListGeolocationParams{
		Limit: params.Limit,
	})
	if err != nil {
		return nil, err
	}

	geolocations := make([]Geolocation, len(result))
	for i, geolocation := range result {
		geolocations[i] = buildGeolocation(geolocation)
	}

	return geolocations, nil
}

func (s *GeolocationServiceImpl) ListByBounding(ctx context.Context, params ListByBoundingParams) ([]Geolocation, error) {
	cmd := s.redis.GeoSearchLocation(ctx, "geolocations", &redis.GeoSearchLocationQuery{
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

func (s *GeolocationServiceImpl) geoPut(ctx context.Context, geolocation Geolocation) error {
	cmd := s.redis.GeoAdd(ctx, "geolocations", &redis.GeoLocation{
		Name:      fmt.Sprintf(cacheKey, geolocation.VehicleID),
		Latitude:  geolocation.Latitude,
		Longitude: geolocation.Longitude,
	})
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (s *GeolocationServiceImpl) cachePut(ctx context.Context, geolocation Geolocation) error {
	cmd := s.redis.HSet(ctx, fmt.Sprintf(cacheKey, geolocation.VehicleID), geolocation)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (s *GeolocationServiceImpl) cacheGet(ctx context.Context, vehicleID int64) (Geolocation, error) {
	return s.cacheGetRaw(ctx, fmt.Sprintf(cacheKey, vehicleID))
}

func (s *GeolocationServiceImpl) cacheGetRaw(ctx context.Context, key string) (Geolocation, error) {
	cmd := s.redis.HGetAll(ctx, key)
	if err := cmd.Err(); err != nil {
		return Geolocation{}, err
	}

	var geolocation Geolocation
	if err := cmd.Scan(&geolocation); err != nil {
		return Geolocation{}, err
	}

	return geolocation, nil
}

func buildGeolocation(model models.Geolocation) Geolocation {
	return Geolocation{
		Degree:    model.Degree,
		Latitude:  model.Latitude,
		Longitude: model.Longitude,
		Speed:     model.Speed,
		VehicleID: model.VehicleID,
		VariantID: model.VariantID,
		Timestamp: model.Timestamp,
	}
}
