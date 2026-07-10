package memory

import (
	"context"
	"strings"

	reference "new-website-lelang/internal/domain/masterdata"
)

// ReferenceRepository is an in-memory adapter. Replace this adapter with a
// database-backed implementation without changing the domain service.
type ReferenceRepository struct{}

func NewReferenceRepository() *ReferenceRepository {
	return &ReferenceRepository{}
}

func (r *ReferenceRepository) GetAll(_ context.Context, value string) (reference.Data, error) {
	data := reference.Data{
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
		Districts: []reference.District{
			{ID: "uuid-1", CityID: "uuid-1", Name: "Menteng"},
			{ID: "uuid-2", CityID: "uuid-1", Name: "Tanah Abang"},
			{ID: "uuid-3", CityID: "uuid-3", Name: "Coblong"},
		},
		SalesMethods: []reference.SalesMethod{
			{ID: "uuid-1", Name: "Lelang"},
			{ID: "uuid-2", Name: "Jual Damai"},
		},
		KPKNLs: []reference.KPKNL{
			{ID: "uuid-1", Code: "KPKNL-JKT1", Name: "KPKNL Jakarta I"},
			{ID: "uuid-2", Code: "KPKNL-JKT2", Name: "KPKNL Jakarta II"},
		},
	}
	if value == "" {
		return data, nil
	}
	value = strings.ToUpper(strings.TrimSpace(value))
	data.Categories = filterCategories(data.Categories, value)
	data.AssetTypes = filterAssetTypes(data.AssetTypes, value)
	data.Provinces = filterProvinces(data.Provinces, value)
	data.Districts = filterDistricts(data.Districts, value)
	return data, nil
}

func (r *ReferenceRepository) GetCitiesByProvinceID(_ context.Context, provinceID, value string) ([]reference.City, error) {
	cities := []reference.City{
		{ID: "uuid-1", ProvinceID: "uuid-1", Name: "Jakarta Pusat"},
		{ID: "uuid-2", ProvinceID: "uuid-1", Name: "Jakarta Selatan"},
		{ID: "uuid-3", ProvinceID: "uuid-2", Name: "Bandung"},
	}

	result := make([]reference.City, 0, len(cities))
	for _, city := range cities {
		if city.ProvinceID == provinceID && (strings.TrimSpace(value) == "" || strings.Contains(strings.ToUpper(city.Name), strings.ToUpper(strings.TrimSpace(value)))) {
			result = append(result, city)
		}
	}
	return result, nil
}

func (r *ReferenceRepository) GetDistrictsByCityID(_ context.Context, cityID, value string) ([]reference.District, error) {
	districts := []reference.District{{ID: "uuid-1", CityID: "uuid-1", Name: "Menteng"}, {ID: "uuid-2", CityID: "uuid-1", Name: "Tanah Abang"}, {ID: "uuid-3", CityID: "uuid-3", Name: "Coblong"}}
	result := make([]reference.District, 0)
	for _, item := range districts {
		if item.CityID == cityID && (strings.TrimSpace(value) == "" || strings.Contains(strings.ToUpper(item.Name), strings.ToUpper(strings.TrimSpace(value)))) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *ReferenceRepository) GetProvinces(_ context.Context, value string) ([]reference.Province, error) {
	provinces := []reference.Province{{ID: "uuid-1", Name: "DKI Jakarta"}, {ID: "uuid-2", Name: "Jawa Barat"}}
	result := make([]reference.Province, 0)
	for _, item := range provinces { if strings.TrimSpace(value) == "" || strings.Contains(strings.ToUpper(item.Name), strings.ToUpper(strings.TrimSpace(value))) { result = append(result, item) } }
	return result, nil
}
func (r *ReferenceRepository) GetCategories(_ context.Context, value string) ([]reference.Category, error) { items:=[]reference.Category{{ID:"uuid-1",Name:"Properti"},{ID:"uuid-2",Name:"Kendaraan"}}; return filterCategories(items,strings.ToUpper(strings.TrimSpace(value))),nil }
func (r *ReferenceRepository) GetAssetTypes(_ context.Context, value string) ([]reference.AssetType, error) { items:=[]reference.AssetType{{ID:"uuid-1",CategoryID:"uuid-1",Name:"Rumah"},{ID:"uuid-2",CategoryID:"uuid-1",Name:"Ruko"},{ID:"uuid-3",CategoryID:"uuid-2",Name:"Mobil"}}; return filterAssetTypes(items,strings.ToUpper(strings.TrimSpace(value))),nil }

func filterCategories(items []reference.Category, value string) []reference.Category {
	result := make([]reference.Category, 0, len(items))
	for _, item := range items {
		if strings.Contains(strings.ToUpper(item.Name), value) {
			result = append(result, item)
		}
	}
	return result
}

func filterAssetTypes(items []reference.AssetType, value string) []reference.AssetType {
	result := make([]reference.AssetType, 0, len(items))
	for _, item := range items {
		if strings.Contains(strings.ToUpper(item.Name), value) {
			result = append(result, item)
		}
	}
	return result
}

func filterProvinces(items []reference.Province, value string) []reference.Province {
	result := make([]reference.Province, 0, len(items))
	for _, item := range items {
		if strings.Contains(strings.ToUpper(item.Name), value) {
			result = append(result, item)
		}
	}
	return result
}

func filterDistricts(items []reference.District, value string) []reference.District {
	result := make([]reference.District, 0, len(items))
	for _, item := range items {
		if strings.Contains(strings.ToUpper(item.Name), value) {
			result = append(result, item)
		}
	}
	return result
}
