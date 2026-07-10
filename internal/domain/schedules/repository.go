package schedules

import "context"

type Repository interface {
	Search(context.Context, Query) (Result, error)
}
