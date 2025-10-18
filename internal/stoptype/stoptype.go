package stoptype

import (
	"context"

	"github.com/catouberos/transit-radar/internal/models"
	"github.com/redis/go-redis/v9"
)

type StopTypeService interface {
	Create(context.Context, CreateParams) (StopType, error)
	Update(context.Context, UpdateParams) (StopType, error)
	Get(context.Context, GetParams) (StopType, error)
	List(context.Context, ListParams) ([]StopType, error)
}

type StopType struct {
	ID   int64  `redis:"-"`
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

func NewStopTypeService(query *models.Queries, redis *redis.Client) StopTypeService {
	return &StopTypeServiceImpl{
		query: query,
		redis: redis,
	}
}

type StopTypeServiceImpl struct {
	query *models.Queries
	redis *redis.Client
}

func (s *StopTypeServiceImpl) Create(ctx context.Context, params CreateParams) (StopType, error) {
	result, err := s.query.CreateStopType(ctx, params.Name)
	if err != nil {
		return StopType{}, err
	}
	stop := buildStopType(result)

	return stop, nil
}

func (s *StopTypeServiceImpl) Update(ctx context.Context, params UpdateParams) (StopType, error) {
	result, err := s.query.UpdateStopType(ctx, models.UpdateStopTypeParams{
		ID:   params.ID,
		Name: params.Name,
	})
	if err != nil {
		return StopType{}, err
	}
	stop := buildStopType(result)

	return stop, nil
}

func (s *StopTypeServiceImpl) Get(ctx context.Context, params GetParams) (StopType, error) {
	result, err := s.query.GetStopType(ctx, models.GetStopTypeParams{
		ID: params.ID,
	})
	if err != nil {
		return StopType{}, err
	}
	stop := buildStopType(result)

	return stop, nil
}

func (s *StopTypeServiceImpl) List(ctx context.Context, params ListParams) ([]StopType, error) {
	return nil, nil
}

func buildStopType(model models.StopType) StopType {
	return StopType{
		ID:   model.ID,
		Name: model.Name,
	}
}
