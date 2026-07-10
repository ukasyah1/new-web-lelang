package masterdata

import "context"

// Repository is the output port for retrieving auction reference data.
type Repository interface {
	GetAll(ctx context.Context, value string) (Data, error)
	GetCitiesByProvinceID(ctx context.Context, provinceID, value string) ([]City, error)
	GetDistrictsByCityID(ctx context.Context, cityID, value string) ([]District, error)
	GetProvinces(ctx context.Context, value string) ([]Province, error)
	GetCategories(ctx context.Context, value string) ([]Category, error)
	GetAssetTypes(ctx context.Context, value string) ([]AssetType, error)
}
