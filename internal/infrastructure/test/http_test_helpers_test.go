package infrastructure_test

import (
	"context"
	"net/http"

	"new-website-lelang/internal/domain/assets"
	reference "new-website-lelang/internal/domain/masterdata"
	"new-website-lelang/internal/infrastructure/memory"
	"new-website-lelang/internal/interfaces/httpapi"
)

func newTestRouter() http.Handler {
	referenceService := reference.NewService(memory.NewReferenceRepository())
	return httpapi.NewRouter(
		reference.NewReferenceHandler(referenceService),
		assets.NewAssetHandler(assets.NewService(testAssetRepository{})),
	)
}

type testAssetRepository struct{}

func (testAssetRepository) Search(context.Context, assets.SearchQuery) (assets.SearchResult, error) {
	return assets.SearchResult{Total: 1, Assets: []assets.Asset{{ID: "uuid-101", Code: "AG-JKT-001", Name: "Rumah Strategis Tebet"}}}, nil
}

func (testAssetRepository) GetByID(_ context.Context, id string) (assets.Asset, error) {
	if id != "uuid-101" {
		return assets.Asset{}, assets.ErrNotFound
	}
	return assets.Asset{ID: "uuid-101", Code: "AG-JKT-001", Name: "Rumah Strategis Tebet"}, nil
}
