package dto

type VariantStopByEbmsIDImport struct {
	// route ID here is neccessary because variant is not unique, but route+variant ebms ID is
	RouteEbmsID   int64 `json:"routeEbmsId"`
	VariantEbmsID int64 `json:"variantEbmsId"`
	StopEbmsID    int64 `json:"stopEbmsId"`
	OrderScore    int32 `json:"orderScore"`
}
