package dto

import "time"

type GeolocationInsert struct {
	Degree      float32   `json:"degree"`
	Latitude    float32   `json:"latitude"`
	Longitude   float32   `json:"longitude"`
	Speed       float32   `json:"speed"`
	VehicleID   int64     `json:"vehicleID"`
	RouteID     int64     `json:"routeID"`
	VariationID int64     `json:"variationID"`
	Timestamp   time.Time `json:"timestamp"`
}

type GeolocationByPlateAndBoundInsert struct {
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
