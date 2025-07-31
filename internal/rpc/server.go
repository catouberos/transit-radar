package rpc

import (
	"context"

	"github.com/catouberos/geoloc/base"
	"github.com/catouberos/geoloc/protos"
)

type GeolocationServer struct {
	protos.UnimplementedGeolocationServer
	App *base.App
}

var _ protos.GeolocationServer = (*GeolocationServer)(nil)

func (g *GeolocationServer) GetVehiclesByRoute(_ context.Context, in *protos.VehiclesByRouteRequest) (*protos.VehiclesByRouteResponse, error) {
	return nil, nil
}

func (g *GeolocationServer) GetVehiclesByStation(_ context.Context, in *protos.VehiclesByStationRequest) (*protos.VehiclesByStationResponse, error) {
	return nil, nil
}
