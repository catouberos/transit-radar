package variantstop

import (
	"context"

	"github.com/catouberos/transit-radar/internal/models"
	"github.com/redis/go-redis/v9"
)

type VariantStopService interface {
	Create(context.Context, CreateParams) (VariantStop, error)
	Update(context.Context, UpdateParams) (VariantStop, error)
	List(context.Context, ListParams) ([]VariantStop, error)
	Delete(context.Context, DeleteParams) error
}

type VariantStop struct {
	VariantID  int64 `redis:"-"`
	StopID     int64 `redis:"-"`
	OrderScore int32 `redis:"-"`
}

type CreateParams struct {
	VariantID  int64
	StopID     int64
	OrderScore int32
}

type UpdateParams struct {
	VariantID  int64
	StopID     int64
	OrderScore int32
}

type ListParams struct {
	VariantID *int64
	StopID    *int64
}

type DeleteParams struct {
	VariantID  int64
	OrderScore int32
}

func NewVariantStopService(query *models.Queries, redis *redis.Client) VariantStopService {
	return &VariantStopServiceImpl{
		query: query,
		redis: redis,
	}
}

type VariantStopServiceImpl struct {
	query *models.Queries
	redis *redis.Client
}

func (s *VariantStopServiceImpl) Create(ctx context.Context, params CreateParams) (VariantStop, error) {
	result, err := s.query.CreateVariantStop(ctx, models.CreateVariantStopParams{
		VariantID:  params.VariantID,
		StopID:     params.StopID,
		OrderScore: params.OrderScore,
	})
	if err != nil {
		return VariantStop{}, err
	}
	stop := buildVariantStop(result)

	return stop, nil
}

func (s *VariantStopServiceImpl) Update(ctx context.Context, params UpdateParams) (VariantStop, error) {
	result, err := s.query.UpdateVariantStop(ctx, models.UpdateVariantStopParams{
		VariantID:  params.VariantID,
		StopID:     params.StopID,
		OrderScore: params.OrderScore,
	})
	if err != nil {
		return VariantStop{}, err
	}
	stop := buildVariantStop(result)

	return stop, nil
}

func (s *VariantStopServiceImpl) List(ctx context.Context, params ListParams) ([]VariantStop, error) {
	results, err := s.query.ListVariantStop(ctx, models.ListVariantStopParams{
		VariantID: params.VariantID,
		StopID:    params.StopID,
	})
	if err != nil {
		return nil, err
	}

	variantsStops := make([]VariantStop, 0, len(results))
	for _, result := range results {
		variantsStops = append(variantsStops, buildVariantStop(result))
	}

	return variantsStops, nil
}

func (s *VariantStopServiceImpl) Delete(ctx context.Context, params DeleteParams) error {
	err := s.query.DeleteVariantStop(ctx, models.DeleteVariantStopParams{
		VariantID:  params.VariantID,
		OrderScore: params.OrderScore,
	})
	if err != nil {
		return err
	}

	return nil
}

func buildVariantStop(model models.VariantsStop) VariantStop {
	return VariantStop{
		VariantID:  model.VariantID,
		StopID:     model.StopID,
		OrderScore: model.OrderScore,
	}
}
