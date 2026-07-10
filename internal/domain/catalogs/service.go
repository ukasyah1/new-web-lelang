package catalogs

import "context"

type Service struct{ repository Repository }

func NewService(repository Repository) *Service                   { return &Service{repository: repository} }
func (s *Service) GetLatest(ctx context.Context) (Catalog, error) { return s.repository.GetLatest(ctx) }
func (s *Service) GetActive(ctx context.Context) (MonthlyCatalog, error) {
	return s.repository.GetActive(ctx)
}
