package infrastructure_test

import (
	"net/http"

	"new-website-lelang/internal/domain/catalog"
	"new-website-lelang/internal/domain/reference"
	"new-website-lelang/internal/infrastructure/memory"
	"new-website-lelang/internal/interfaces/httpapi"
)

func newTestRouter() http.Handler {
	referenceService := reference.NewService(memory.NewReferenceRepository())
	return httpapi.NewRouter(
		reference.NewReferenceHandler(referenceService),
		catalog.NewAssetHandler(),
	)
}
