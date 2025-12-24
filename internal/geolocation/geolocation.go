package geolocation

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/catouberos/transit-radar/internal/models"
	"github.com/catouberos/transit-radar/types"
	"github.com/cridenour/go-postgis"
	"github.com/redis/go-redis/v9"
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
	Limit  int32
	Offset int32
}

type ListByBoundingParams struct {
	Latitude, Longitude, Width, Height float64
	Unit                               string
}

var _ GeolocationService = (*geolocationService)(nil)

func NewGeolocationService(query *models.Queries, redis *redis.Client) GeolocationService {
	return &geolocationService{
		query: query,
		redis: redis,
	}
}

type geolocationService struct {
	query *models.Queries
	redis *redis.Client
}

func (s *geolocationService) Create(ctx context.Context, params CreateParams) (Geolocation, error) {
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

	if err := s.cachePut(ctx, geolocation); err != nil {
		slog.WarnContext(ctx, "cannot cache geolocation", "error", err)
	}
	if err := s.geoPut(ctx, geolocation); err != nil {
		slog.WarnContext(ctx, "cannot cache geolocation", "error", err)
	}

	return geolocation, nil
}

func (s *geolocationService) Get(ctx context.Context, params GetParams) (Geolocation, error) {
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

func (s *geolocationService) List(ctx context.Context, params ListParams) ([]Geolocation, error) {
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

func (s *geolocationService) ListByBounding(ctx context.Context, params ListByBoundingParams) ([]Geolocation, error) {
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

func buildGeolocation(model models.Geolocation) Geolocation {
	return Geolocation{
		Degree: model.Degree,
		Location: types.LatLng{
			// in lng/lat order
			Longitude: model.Location.X,
			Latitude:  model.Location.Y,
		},
		Speed:     model.Speed,
		VehicleID: model.VehicleID,
		VariantID: model.VariantID,
		Timestamp: model.Timestamp,
	}
}
