package reference

import "context"

// Repository is the output port for retrieving auction reference data.
type Repository interface {
	GetAll(ctx context.Context) (Data, error)
}
