package route

import (
	"context"
	"fmt"

	"github.com/catouberos/transit-radar/internal/models"
	"github.com/redis/go-redis/v9"
)

const (
	RouteCacheKey = "route-%d"
)

type RouteService interface {
	Create(context.Context, CreateParams) (Route, error)
	Update(context.Context, UpdateParams) (Route, error)
	Get(context.Context, GetParams) (Route, error)
	GetByEbmsID(context.Context, GetByEbmsIDParams) (Route, error)
	List(context.Context, ListParams) ([]Route, error)
}

type Route struct {
	ID            int64   `redis:"-"`
	Number        string  `redis:"number"`
	Name          string  `redis:"name"`
	EbmsID        *int64  `redis:"ebms_id"`
	Active        bool    `redis:"active"`
	OperationTime *string `redis:"operation_time"`
	Organization  *string `redis:"organization"`
	Ticketing     *string `redis:"ticketing"`
	RouteType     *string `redis:"route_type"`
}

type CreateParams struct {
	Number        string
	Name          string
	EbmsID        *int64
	Active        bool
	OperationTime *string
	Organization  *string
	Ticketing     *string
	RouteType     *string
}

type UpdateParams struct {
	ID            int64
	Number        *string
	Name          *string
	EbmsID        *int64
	Active        *bool
	OperationTime *string
	Organization  *string
	Ticketing     *string
	RouteType     *string
}

type GetParams struct {
	ID int64
}

type GetByEbmsIDParams struct {
	EbmsID int64
}

type ListParams struct {
}

// IMPLEMENTATION

func NewRouteService(query *models.Queries, redis *redis.Client) RouteService {
	return RouteServiceImpl{
		query: query,
		redis: redis,
	}
}

type RouteServiceImpl struct {
	query *models.Queries
	redis *redis.Client
}

func (s RouteServiceImpl) Create(ctx context.Context, params CreateParams) (Route, error) {
	result, err := s.query.CreateRoute(ctx, models.CreateRouteParams{
		Number:        params.Number,
		Name:          params.Name,
		EbmsID:        params.EbmsID,
		OperationTime: params.OperationTime,
		Organization:  params.Organization,
		Ticketing:     params.Ticketing,
		RouteType:     params.RouteType,
	})
	if err != nil {
		return Route{}, err
	}
	route := buildRoute(result)

	err = s.cachePut(ctx, route)
	if err != nil {
		return route, err
	}

	return route, nil
}

func (s RouteServiceImpl) Update(ctx context.Context, params UpdateParams) (Route, error) {
	result, err := s.query.UpdateRoute(ctx, models.UpdateRouteParams{
		Number:        params.Number,
		Name:          params.Name,
		EbmsID:        params.EbmsID,
		OperationTime: params.OperationTime,
		Organization:  params.Organization,
		Ticketing:     params.Ticketing,
		RouteType:     params.RouteType,
		ID:            params.ID,
	})

	if err != nil {
		return Route{}, err
	}
	route := buildRoute(result)

	err = s.cachePut(ctx, route)
	if err != nil {
		return route, err
	}

	return route, nil
}

func (s RouteServiceImpl) Get(ctx context.Context, params GetParams) (Route, error) {
	route, err := s.cacheGet(ctx, params.ID)
	if err == nil {
		return route, err
	}

	result, err := s.query.GetRoute(ctx, models.GetRouteParams{
		ID: &params.ID,
	})
	if err != nil {
		return Route{}, err
	}
	route = buildRoute(result)

	return route, nil
}

func (s RouteServiceImpl) GetByEbmsID(ctx context.Context, params GetByEbmsIDParams) (Route, error) {
	result, err := s.query.GetRoute(ctx, models.GetRouteParams{
		EbmsID: &params.EbmsID,
	})
	if err != nil {
		return Route{}, err
	}

	return buildRoute(result), nil

}

func (s RouteServiceImpl) List(ctx context.Context, params ListParams) ([]Route, error) {
	result, err := s.query.ListRoute(ctx)
	if err != nil {
		return nil, err
	}

	routes := make([]Route, len(result))
	for i, route := range result {
		routes[i] = buildRoute(route)
	}

	return routes, nil
}

func (s RouteServiceImpl) cachePut(ctx context.Context, route Route) error {
	key := fmt.Sprintf(RouteCacheKey, route.ID)
	_, err := s.redis.HSet(ctx, key, route).Result()
	return err
}

func (s RouteServiceImpl) cacheGet(ctx context.Context, routeID int64) (Route, error) {
	key := fmt.Sprintf(RouteCacheKey, routeID)

	route := Route{}
	err := s.redis.HGetAll(ctx, key).Scan(&route)
	if err != nil {
		return Route{}, err
	}

	return route, nil
}

func buildRoute(model models.Route) Route {
	return Route{
		ID:            model.ID,
		Number:        model.Number,
		Name:          model.Name,
		EbmsID:        model.EbmsID,
		Active:        model.Active,
		OperationTime: model.OperationTime,
		Organization:  model.Organization,
		Ticketing:     model.Ticketing,
		RouteType:     model.RouteType,
	}
}
