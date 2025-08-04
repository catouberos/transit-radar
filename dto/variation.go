package dto

type VariationUpsert struct {
	Name       string `json:"name"`
	EbmsID     int64  `json:"ebmsID"`
	IsOutbound bool   `json:"isOutbound"`
	RouteID    int64  `json:"routeID"`
}

type VariationByRouteEbmsIDUpsert struct {
	Name        string `json:"name"`
	EbmsID      int64  `json:"ebmsID"`
	IsOutbound  bool   `json:"isOutbound"`
	RouteEbmsID int64  `json:"routeEbmsID"`
}
