package assets

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	reference "new-website-lelang/internal/domain/masterdata"
	"new-website-lelang/internal/infrastructure/memory"
	"new-website-lelang/internal/interfaces/httpapi"
)

func newTestRouter() http.Handler {
	referenceService := reference.NewService(memory.NewReferenceRepository())
	assetService := NewService(testAssetRepository{})
	return httpapi.NewRouter(reference.NewReferenceHandler(referenceService), NewAssetHandler(assetService))
}

type testAssetRepository struct{}

func (testAssetRepository) Search(context.Context, SearchQuery) (SearchResult, error) {
	return SearchResult{Total: 1, Assets: []Asset{{ID: "uuid-101", Code: "AG-JKT-001", Name: "Rumah Strategis Tebet"}}}, nil
}

func (testAssetRepository) GetByID(_ context.Context, id string) (Asset, error) {
	if id != "uuid-101" {
		return Asset{}, ErrNotFound
	}
	return Asset{
		ID: "uuid-101", Code: "AG-JKT-001", Name: "Rumah Strategis Tebet",
		Category:  NamedReference{ID: "uuid-1", Name: "Properti"},
		AssetType: NamedReference{ID: "uuid-1", Name: "Rumah"},
		Province:  NamedReference{ID: "uuid-1", Name: "DKI Jakarta"},
		City:      NamedReference{ID: "uuid-1", Name: "Jakarta Selatan"},
		Address:   "Jl. Teuku Umar No. 12", Coordinates: "-6.2291,106.8524",
		Description: "Rumah strategis", ImageURLs: []string{"https://example.com/image.jpg"},
		Facilities: []string{"Carport"},
	}, nil
}

func TestGetAssets(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/assets?search=rumah&metode_penjualan_id[]=1&metode_penjualan_id[]=2&page=1&limit=10",
		nil,
	)
	recorder := httptest.NewRecorder()

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response assetListResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != "success" {
		t.Fatalf("expected success status, got %q", response.Status)
	}
	if response.Meta.TotalData != 1 || response.Meta.TotalPages != 1 {
		t.Fatalf("unexpected meta: %+v", response.Meta)
	}
	if len(response.Data) != 1 || response.Data[0].CollateralCode != "AG-JKT-001" {
		t.Fatalf("unexpected asset data: %+v", response.Data)
	}
}

func TestGetAssetsRejectsInvalidPriceRange(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/assets?harga_min=500&harga_max=100", nil)
	recorder := httptest.NewRecorder()

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestGetAssetByID(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/assets/uuid-101", nil)
	recorder := httptest.NewRecorder()

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response assetDetailAPIResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != "success" || response.Data.ID != "uuid-101" || response.Data.Coordinates != "-6.2291,106.8524" {
		t.Fatalf("unexpected asset detail response: %+v", response)
	}
}

func TestGetAssetByIDReturnsNotFound(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/assets/missing", nil)
	recorder := httptest.NewRecorder()

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}
}
