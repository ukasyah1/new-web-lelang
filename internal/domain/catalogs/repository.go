package catalogs

import "context"

type Repository interface {
	GetLatest(context.Context) (Catalog, error)
	GetActive(context.Context) (MonthlyCatalog, error)
}
