package dto

type VariantUpsert struct {
	Name       string `json:"name"`
	EbmsID     int64  `json:"ebmsID"`
	IsOutbound bool   `json:"isOutbound"`
	RouteID    int64  `json:"routeID"`
}

type VariantByRouteEbmsIDUpsert struct {
	Name        string `json:"name"`
	EbmsID      int64  `json:"ebmsID"`
	IsOutbound  bool   `json:"isOutbound"`
	RouteEbmsID int64  `json:"routeEbmsID"`
}
