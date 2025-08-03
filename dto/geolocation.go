package dto

import "time"

type GeolocationInsert struct {
	Degree    float32   `json:"deg"`
	Latitude  float32   `json:"lat"`
	Longitude float32   `json:"lng"`
	Speed     float32   `json:"spd"`
	VehicleID int64     `json:"vID"`
	RouteID   int64     `json:"rID"`
	Timestamp time.Time `json:"dts"`
}
