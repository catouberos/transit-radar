package dto

type VariantStopByEbmsIDImport struct {
	VariantEbmsID int64 `json:"variantEbmsId"`
	StopEbmsID    int64 `json:"stopEbmsdD"`
	OrderScore    int32 `json:"orderScore"`
}
