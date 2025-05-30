package rpc

import (
	"context"

	"github.com/catouberos/geoloc/base"
	"github.com/catouberos/geoloc/models"
	pb "github.com/catouberos/geoloc/protos"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GeolocationServer struct {
	pb.UnimplementedGeolocationServer
	App *base.App
}

func (g *GeolocationServer) CreateGeolocation(ctx context.Context, in *pb.GeolocationRequest) (*pb.GeolocationResponse, error) {
	geolocation, err := g.App.Queries.CreateGeolocation(ctx, models.CreateGeolocationParams{
		Degree:    pgtype.Float4{Float32: in.Degree, Valid: true},
		Latitude:  pgtype.Float4{Float32: in.Latitude, Valid: true},
		Longitude: pgtype.Float4{Float32: in.Longitude, Valid: true},
		Speed:     pgtype.Float4{Float32: in.Speed, Valid: true},
		VehicleID: pgtype.Int4{Int32: in.VehicleId, Valid: true},
		RouteID:   pgtype.Int4{Int32: in.RouteId, Valid: true},
		Timestamp: pgtype.Timestamptz{Time: in.Timestamp.AsTime(), Valid: true},
	})
	return &pb.GeolocationResponse{
		Degree:    geolocation.Degree.Float32,
		Latitude:  geolocation.Latitude.Float32,
		Longitude: geolocation.Longitude.Float32,
		Speed:     geolocation.Speed.Float32,
		VehicleId: geolocation.VehicleID.Int32,
		RouteId:   geolocation.RouteID.Int32,
		Timestamp: timestamppb.New(geolocation.Timestamp.Time),
	}, err
}
