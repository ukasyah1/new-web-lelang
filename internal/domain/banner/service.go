package banner

import "context"

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetActive(ctx context.Context) ([]Banner, error) {
	return s.repository.GetActive(ctx)
}
