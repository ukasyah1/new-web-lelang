package banner

import "context"

type Repository interface {
	GetActive(ctx context.Context) ([]Banner, error)
}
