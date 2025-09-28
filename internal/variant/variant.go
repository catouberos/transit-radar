package variant

import (
	"context"

	"github.com/catouberos/transit-radar/internal/models"
	"github.com/redis/go-redis/v9"
)

type VariantService interface {
	Create(context.Context, CreateParams) (Variant, error)
	Update(context.Context, UpdateParams) (Variant, error)
	Get(context.Context, GetParams) (Variant, error)
	GetByEbmsID(context.Context, GetByEbmsIDParams) (Variant, error)
	List(context.Context, ListParams) ([]Variant, error)
}

type Variant struct {
	ID            int64    `redis:"-"`
	Name          string   `redis:"name"`
	EbmsID        *int64   `redis:"ebms_id"`
	IsOutbound    bool     `redis:"is_outbound"`
	RouteID       int64    `redis:"route_id"`
	Description   *string  `redis:"description"`
	ShortName     *string  `redis:"short_name"`
	Distance      *float32 `redis:"distance"`
	Duration      *int32   `redis:"duration"`
	StartStopName *string  `redis:"start_stop_name"`
	EndStopName   *string  `redis:"end_stop_name"`
}

type CreateParams struct {
	Name          string
	EbmsID        *int64
	IsOutbound    bool
	RouteID       int64
	Description   *string
	ShortName     *string
	Distance      *float32
	Duration      *int32
	StartStopName *string
	EndStopName   *string
}

type UpdateParams struct {
	ID            int64
	Name          *string
	EbmsID        *int64
	IsOutbound    *bool
	RouteID       *int64
	Description   *string
	ShortName     *string
	Distance      *float32
	Duration      *int32
	StartStopName *string
	EndStopName   *string
}

type GetParams struct {
	ID int64
}

type GetByEbmsIDParams struct {
	EbmsID int64
}

type ListParams struct {
}

func NewVariantService(query *models.Queries, redis *redis.Client) VariantService {
	return VariantServiceImpl{
		query: query,
		redis: redis,
	}
}

type VariantServiceImpl struct {
	query *models.Queries
	redis *redis.Client
}

func (s VariantServiceImpl) Create(ctx context.Context, params CreateParams) (Variant, error) {
	result, err := s.query.CreateVariant(ctx, models.CreateVariantParams{
		Name:          params.Name,
		EbmsID:        params.EbmsID,
		IsOutbound:    params.IsOutbound,
		RouteID:       params.RouteID,
		Description:   params.Description,
		ShortName:     params.ShortName,
		Distance:      params.Distance,
		Duration:      params.Duration,
		StartStopName: params.StartStopName,
		EndStopName:   params.EndStopName,
	})
	if err != nil {
		return Variant{}, err
	}
	stop := buildVariant(result)

	return stop, nil
}

func (s VariantServiceImpl) Update(ctx context.Context, params UpdateParams) (Variant, error) {
	result, err := s.query.UpdateVariant(ctx, models.UpdateVariantParams{
		Name:          params.Name,
		EbmsID:        params.EbmsID,
		IsOutbound:    params.IsOutbound,
		RouteID:       params.RouteID,
		Description:   params.Description,
		ShortName:     params.ShortName,
		Distance:      params.Distance,
		Duration:      params.Duration,
		StartStopName: params.StartStopName,
		EndStopName:   params.EndStopName,
	})
	if err != nil {
		return Variant{}, err
	}
	stop := buildVariant(result)

	return stop, nil
}

func (s VariantServiceImpl) Get(ctx context.Context, params GetParams) (Variant, error) {
	result, err := s.query.GetVariant(ctx, models.GetVariantParams{
		ID: &params.ID,
	})
	if err != nil {
		return Variant{}, err
	}
	stop := buildVariant(result)

	return stop, nil
}

func (s VariantServiceImpl) GetByEbmsID(ctx context.Context, params GetByEbmsIDParams) (Variant, error) {
	result, err := s.query.GetVariant(ctx, models.GetVariantParams{
		EbmsID: &params.EbmsID,
	})
	if err != nil {
		return Variant{}, err
	}
	stop := buildVariant(result)

	return stop, nil
}

func (s VariantServiceImpl) List(ctx context.Context, params ListParams) ([]Variant, error) {
	return nil, nil
}

func buildVariant(model models.Variant) Variant {
	return Variant{
		ID:            model.ID,
		Name:          model.Name,
		EbmsID:        model.EbmsID,
		IsOutbound:    model.IsOutbound,
		RouteID:       model.ID,
		Description:   model.Description,
		ShortName:     model.ShortName,
		Distance:      model.Distance,
		Duration:      model.Duration,
		StartStopName: model.StartStopName,
		EndStopName:   model.EndStopName,
	}
}
