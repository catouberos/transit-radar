package dto

type RouteUpsert struct {
	Number string `json:"number"`
	Name   string `json:"name"`
	EbmsID int64  `json:"ebmsID"`
	Active bool   `json:"active"`
}
