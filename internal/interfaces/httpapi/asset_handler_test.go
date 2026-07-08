package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"new-website-lelang/internal/domain/reference"
	"new-website-lelang/internal/infrastructure/memory"
)

func newTestRouter() http.Handler {
	referenceService := reference.NewService(memory.NewReferenceRepository())
	return NewRouter(NewReferenceHandler(referenceService), NewAssetHandler())
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
	if response.Meta.TotalData != 45 || response.Meta.TotalPages != 5 {
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
