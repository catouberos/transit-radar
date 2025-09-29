package stop

import (
	"context"

	"github.com/catouberos/transit-radar/internal/models"
	"github.com/redis/go-redis/v9"
)

type StopService interface {
	Create(context.Context, CreateParams) (Stop, error)
	// Update(context.Context, UpdateParams) (Stop, error)
	Get(context.Context, GetParams) (Stop, error)
	GetByEbmsID(context.Context, GetByEbmsIDParams) (Stop, error)
	List(context.Context, ListParams) ([]Stop, error)
}

type Stop struct {
	ID        int64   `redis:"-"`
	Code      string  `redis:"code"`
	Name      string  `redis:"name"`
	TypeID    int64   `redis:"type_id"`
	EbmsID    *int64  `redis:"ebms_id"`
	Active    bool    `redis:"active"`
	Latitude  float32 `redis:"latitude"`
	Longitude float32 `redis:"longitude"`
}

type CreateParams struct {
	Code      string
	Name      string
	TypeID    int64
	EbmsID    *int64
	Active    bool
	Latitude  float32
	Longitude float32
}

type UpdateParams struct {
	ID        int64
	Code      *string
	Name      *string
	TypeID    *int64
	EbmsID    *int64
	Active    *bool
	Latitude  *float32
	Longitude *float32
}

type GetParams struct {
	ID int64
}

type GetByEbmsIDParams struct {
	EbmsID int64
}

type ListParams struct {
}

func NewStopService(query *models.Queries, redis *redis.Client) StopService {
	return StopServiceImpl{
		query: query,
		redis: redis,
	}
}

type StopServiceImpl struct {
	query *models.Queries
	redis *redis.Client
}

func (s StopServiceImpl) Create(ctx context.Context, params CreateParams) (Stop, error) {
	result, err := s.query.CreateStop(ctx, models.CreateStopParams{
		Code:      params.Code,
		Name:      params.Name,
		TypeID:    params.TypeID,
		EbmsID:    params.EbmsID,
		Active:    params.Active,
		Latitude:  params.Latitude,
		Longitude: params.Longitude,
	})
	if err != nil {
		return Stop{}, err
	}
	stop := buildStop(result)

	return stop, nil
}

func (s StopServiceImpl) Get(ctx context.Context, params GetParams) (Stop, error) {
	result, err := s.query.GetStop(ctx, models.GetStopParams{
		ID: &params.ID,
	})
	if err != nil {
		return Stop{}, err
	}
	stop := buildStop(result)

	return stop, nil
}

func (s StopServiceImpl) GetByEbmsID(ctx context.Context, params GetByEbmsIDParams) (Stop, error) {
	result, err := s.query.GetStop(ctx, models.GetStopParams{
		EbmsID: &params.EbmsID,
	})
	if err != nil {
		return Stop{}, err
	}
	stop := buildStop(result)

	return stop, nil
}

func (s StopServiceImpl) List(ctx context.Context, params ListParams) ([]Stop, error) {
	return nil, nil
}

func buildStop(model models.Stop) Stop {
	return Stop{
		ID:        model.ID,
		Code:      model.Code,
		Name:      model.Name,
		TypeID:    model.TypeID,
		EbmsID:    model.EbmsID,
		Active:    model.Active,
		Latitude:  model.Latitude,
		Longitude: model.Longitude,
	}
}
