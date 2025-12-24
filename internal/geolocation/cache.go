package geolocation

import "context"

type GeolocationCache interface {
	Put(context.Context, Geolocation) error
	Get(ctx context.Context, vehicleID string) (Geolocation, error)
	ListByBounding(context.Context, ListByBoundingParams) ([]Geolocation, error)
}
