package rpc

import (
	"context"

	"github.com/catouberos/geoloc/base"
	"github.com/catouberos/geoloc/models"
	pb "github.com/catouberos/geoloc/protos"
)

type GeolocationServer struct {
	pb.UnimplementedGeolocationServer
	App *base.App
}

func (g *GeolocationServer) CreateGeolocation(_ context.Context, in *pb.GeolocationRequest) (*pb.GeolocationRequest, error) {
	_, err := models.CreateGeolocation(g.App.DB, &models.Geolocation{
		Deg:        in.Deg,
		Lat:        in.Lat,
		Lng:        in.Lng,
		Speed:      in.Speed,
		VehicleId:  in.VehicleId,
		RouteId:    in.RouteId,
		UpdateTime: in.UpdateTime.AsTime(),
	})
	return &pb.GeolocationRequest{}, err
}
