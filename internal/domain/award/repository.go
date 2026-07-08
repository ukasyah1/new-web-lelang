package award

import "context"

type Repository interface {
	GetAll(ctx context.Context) ([]Award, error)
}
