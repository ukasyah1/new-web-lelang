package award

import "context"

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll(ctx context.Context) ([]Award, error) {
	return s.repository.GetAll(ctx)
}
