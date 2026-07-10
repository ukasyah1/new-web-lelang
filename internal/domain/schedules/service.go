package schedules

import "context"

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }
func (s *Service) Search(ctx context.Context, q Query) (Result, error) {
	return s.repository.Search(ctx, q)
}
