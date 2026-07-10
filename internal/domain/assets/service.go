package assets

import "context"

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) Search(ctx context.Context, query SearchQuery) (SearchResult, error) {
	return s.repository.Search(ctx, query)
}

func (s *Service) GetByID(ctx context.Context, id string) (Asset, error) {
	return s.repository.GetByID(ctx, id)
}
