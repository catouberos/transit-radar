package dto

type StopImport struct {
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	TypeName        string  `json:"typeId"` // should automatically search and replace with ID
	EbmsID          int64   `json:"ebmsId"`
	Active          bool    `json:"active"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	AddressNumber   string  `json:"addressNumber"`
	AddressStreet   string  `json:"addressStreet"`
	AddressWard     string  `json:"addressWard"`
	AddressDistrict string  `json:"addressDistrict"`
}
