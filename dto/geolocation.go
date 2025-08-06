package dto

import "time"

type GeolocationByRouteIDAndPlateAndBoundInsert struct {
	Degree       float32 `json:"degree"`
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
	Speed        float32 `json:"speed"`
	LicensePlate string  `json:"licensePlate"`
	RouteID      int64   `json:"routeID"`
	// direction = 0 -> IsOutbound = true
	IsOutbound bool      `json:"isOutbound"`
	Timestamp  time.Time `json:"timestamp"`
}

type Geolocation struct {
	Degree    float32   `json:"geolocation" redis:"geolocation"`
	Latitude  float32   `json:"latitude" redis:"latitude"`
	Longitude float32   `json:"longitude" redis:"longitude"`
	Speed     float32   `json:"speed" redis:"speed"`
	VehicleID int64     `json:"vehicleId" redis:"vehicle_id"`
	VariantID int64     `json:"variantId" redis:"variant_id"`
	Timestamp time.Time `json:"timestamp" redis:"timestamp"`
}
