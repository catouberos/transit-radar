package types

// LatLng stores a location latitude and longitude, based on https://pkg.go.dev/google.golang.org/genproto/googleapis/type/latlng
type LatLng struct {
	// The latitude in degrees. It must be in the range [-90.0, +90.0].
	Latitude float64 `json:"latitude,omitempty"`
	// The longitude in degrees. It must be in the range [-180.0, +180.0].
	Longitude float64 `json:"longitude,omitempty"`
}
