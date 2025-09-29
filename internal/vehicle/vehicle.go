package vehicle

import (
	"context"

	"github.com/catouberos/transit-radar/internal/models"
	"github.com/redis/go-redis/v9"
)

type VehicleService interface {
	Create(context.Context, CreateParams) (Vehicle, error)
	Update(context.Context, UpdateParams) (Vehicle, error)
	Get(context.Context, GetParams) (Vehicle, error)
	List(context.Context, ListParams) ([]Vehicle, error)
}

type Vehicle struct {
	ID           int64  `redis:"-"`
	LicensePlate string `redis:"license_plate"`
	Type         *int64 `redis:"type"`
}

type CreateParams struct {
	LicensePlate string
	Type         *int64
}

type UpdateParams struct {
	ID           int64
	LicensePlate *string
	Type         *int64
}

type GetParams struct {
	ID           *int64
	LicensePlate *string
}

type ListParams struct {
}

func NewVehicleService(query *models.Queries, redis *redis.Client) VehicleService {
	return VehicleServiceImpl{
		query: query,
		redis: redis,
	}
}

type VehicleServiceImpl struct {
	query *models.Queries
	redis *redis.Client
}

func (s VehicleServiceImpl) Create(ctx context.Context, params CreateParams) (Vehicle, error) {
	result, err := s.query.CreateVehicle(ctx, models.CreateVehicleParams{
		LicensePlate: params.LicensePlate,
		Type:         params.Type,
	})
	if err != nil {
		return Vehicle{}, err
	}
	stop := buildVehicle(result)

	return stop, nil
}

func (s VehicleServiceImpl) Update(ctx context.Context, params UpdateParams) (Vehicle, error) {
	result, err := s.query.UpdateVehicle(ctx, models.UpdateVehicleParams{
		ID:           params.ID,
		LicensePlate: params.LicensePlate,
		Type:         params.Type,
	})
	if err != nil {
		return Vehicle{}, err
	}
	stop := buildVehicle(result)

	return stop, nil
}

func (s VehicleServiceImpl) Get(ctx context.Context, params GetParams) (Vehicle, error) {
	result, err := s.query.GetVehicle(ctx, models.GetVehicleParams{
		ID:           params.ID,
		LicensePlate: params.LicensePlate,
	})
	if err != nil {
		return Vehicle{}, err
	}
	stop := buildVehicle(result)

	return stop, nil
}

func (s VehicleServiceImpl) List(ctx context.Context, params ListParams) ([]Vehicle, error) {
	return nil, nil
}

func buildVehicle(model models.Vehicle) Vehicle {
	return Vehicle{
		ID:           model.ID,
		LicensePlate: model.LicensePlate,
		Type:         model.Type,
	}
}
