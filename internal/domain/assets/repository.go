package assets

import "context"

type Repository interface {
	Search(ctx context.Context, query SearchQuery) (SearchResult, error)
	GetByID(ctx context.Context, id string) (Asset, error)
}
