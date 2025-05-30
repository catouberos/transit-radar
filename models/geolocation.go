package models

import (
	"database/sql"
	"time"
)

type Geolocation struct {
	Deg        float32
	Lat        float32
	Lng        float32
	Speed      float32
	VehicleId  int32
	RouteId    int32
	UpdateTime time.Time
}

func CreateGeolocation(db *sql.DB, geolocation *Geolocation) (*Geolocation, error) {
	_, err := db.Query(
		"INSERT INTO public.geolocation(deg, lat, lng, speed, vehicle_id, route_id, update_time) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		geolocation.Deg,
		geolocation.Lat,
		geolocation.Lng,
		geolocation.Speed,
		geolocation.VehicleId,
		geolocation.RouteId,
		geolocation.UpdateTime,
	)
	return nil, err
}
