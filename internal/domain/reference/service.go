package reference

import "context"

// Service contains the reference-data use cases.
type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll(ctx context.Context) (Data, error) {
	return s.repository.GetAll(ctx)
}

func (s *Service) GetCitiesByProvinceID(ctx context.Context, provinceID string) ([]City, error) {
	return s.repository.GetCitiesByProvinceID(ctx, provinceID)
}
