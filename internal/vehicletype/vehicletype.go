package vehicletype

import (
	"context"

	"github.com/catouberos/transit-radar/internal/models"
	"github.com/redis/go-redis/v9"
)

type VehicleTypeService interface {
	Create(context.Context, CreateParams) (VehicleType, error)
	Update(context.Context, UpdateParams) (VehicleType, error)
	Get(context.Context, GetParams) (VehicleType, error)
	List(context.Context, ListParams) ([]VehicleType, error)
}

type VehicleType struct {
	ID   int64  `redis:"id"`
	Name string `redis:"name"`
}

type CreateParams struct {
	Name string
}

type UpdateParams struct {
	ID   int64
	Name *string
}

type GetParams struct {
	ID *int64
}

type ListParams struct {
}

func NewVehicleTypeService(query *models.Queries, redis *redis.Client) VehicleTypeService {
	return VehicleTypeServiceImpl{
		query: query,
		redis: redis,
	}
}

type VehicleTypeServiceImpl struct {
	query *models.Queries
	redis *redis.Client
}

func (s VehicleTypeServiceImpl) Create(ctx context.Context, params CreateParams) (VehicleType, error) {
	result, err := s.query.CreateVehicleType(ctx, params.Name)
	if err != nil {
		return VehicleType{}, err
	}
	stop := buildVehicleType(result)

	return stop, nil
}

func (s VehicleTypeServiceImpl) Update(ctx context.Context, params UpdateParams) (VehicleType, error) {
	result, err := s.query.UpdateVehicleType(ctx, models.UpdateVehicleTypeParams{
		ID:   params.ID,
		Name: params.Name,
	})
	if err != nil {
		return VehicleType{}, err
	}
	stop := buildVehicleType(result)

	return stop, nil
}

func (s VehicleTypeServiceImpl) Get(ctx context.Context, params GetParams) (VehicleType, error) {
	result, err := s.query.GetVehicleType(ctx, models.GetVehicleTypeParams{
		ID: params.ID,
	})
	if err != nil {
		return VehicleType{}, err
	}
	stop := buildVehicleType(result)

	return stop, nil
}

func (s VehicleTypeServiceImpl) List(ctx context.Context, params ListParams) ([]VehicleType, error) {
	return nil, nil
}

func buildVehicleType(model models.VehicleType) VehicleType {
	return VehicleType{
		ID:   model.ID,
		Name: model.Name,
	}
}
