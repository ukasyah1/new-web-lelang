package masterdata

import "context"

// Service contains the reference-data use cases.
type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll(ctx context.Context, value string) (Data, error) {
	return s.repository.GetAll(ctx, value)
}

func (s *Service) GetCitiesByProvinceID(ctx context.Context, provinceID, value string) ([]City, error) {
	return s.repository.GetCitiesByProvinceID(ctx, provinceID, value)
}

func (s *Service) GetDistrictsByCityID(ctx context.Context, cityID, value string) ([]District, error) {
	return s.repository.GetDistrictsByCityID(ctx, cityID, value)
}

func (s *Service) GetProvinces(ctx context.Context, value string) ([]Province, error) {
	return s.repository.GetProvinces(ctx, value)
}
func (s *Service) GetCategories(ctx context.Context, value string) ([]Category, error) { return s.repository.GetCategories(ctx, value) }
func (s *Service) GetAssetTypes(ctx context.Context, value string) ([]AssetType, error) { return s.repository.GetAssetTypes(ctx, value) }
