package memory

import (
	"context"

	"new-website-lelang/internal/domain/reference"
)

// ReferenceRepository is an in-memory adapter. Replace this adapter with a
// database-backed implementation without changing the domain service.
type ReferenceRepository struct{}

func NewReferenceRepository() *ReferenceRepository {
	return &ReferenceRepository{}
}

func (r *ReferenceRepository) GetAll(_ context.Context) (reference.Data, error) {
	return reference.Data{
		Categories: []reference.Category{
			{ID: "uuid-1", Name: "Properti"},
			{ID: "uuid-2", Name: "Kendaraan"},
		},
		AssetTypes: []reference.AssetType{
			{ID: "uuid-1", CategoryID: "uuid-1", Name: "Rumah"},
			{ID: "uuid-2", CategoryID: "uuid-1", Name: "Ruko"},
			{ID: "uuid-3", CategoryID: "uuid-2", Name: "Mobil"},
		},
		Provinces: []reference.Province{
			{ID: "uuid-1", Name: "DKI Jakarta"},
			{ID: "uuid-2", Name: "Jawa Barat"},
		},
		SalesMethods: []reference.SalesMethod{
			{ID: "uuid-1", Name: "Lelang"},
			{ID: "uuid-2", Name: "Jual Damai"},
		},
		KPKNLs: []reference.KPKNL{
			{ID: "uuid-1", Code: "KPKNL-JKT1", Name: "KPKNL Jakarta I"},
			{ID: "uuid-2", Code: "KPKNL-JKT2", Name: "KPKNL Jakarta II"},
		},
	}, nil
}

func (r *ReferenceRepository) GetCitiesByProvinceID(_ context.Context, provinceID string) ([]reference.City, error) {
	cities := []reference.City{
		{ID: "uuid-1", ProvinceID: "uuid-1", Name: "Jakarta Pusat"},
		{ID: "uuid-2", ProvinceID: "uuid-1", Name: "Jakarta Selatan"},
		{ID: "uuid-3", ProvinceID: "uuid-2", Name: "Bandung"},
	}

	result := make([]reference.City, 0, len(cities))
	for _, city := range cities {
		if city.ProvinceID == provinceID {
			result = append(result, city)
		}
	}
	return result, nil
}
