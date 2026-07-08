package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"new-website-lelang/internal/domain/reference"
	"new-website-lelang/internal/infrastructure/memory"
)

func TestGetReferenceData(t *testing.T) {
	service := reference.NewService(memory.NewReferenceRepository())
	router := NewRouter(NewReferenceHandler(service), NewAssetHandler())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/reference-data", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response referenceResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != "success" {
		t.Fatalf("expected success status, got %q", response.Status)
	}
	if len(response.Data.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(response.Data.Categories))
	}
	if response.Data.AssetTypes[2].Name != "Mobil" {
		t.Fatalf("expected third asset type to be Mobil, got %q", response.Data.AssetTypes[2].Name)
	}
	if response.Data.KPKNLs[0].Code != "KPKNL-JKT1" {
		t.Fatalf("unexpected first KPKNL code: %q", response.Data.KPKNLs[0].Code)
	}
}

func TestReferenceDataRejectsUnsupportedMethod(t *testing.T) {
	service := reference.NewService(memory.NewReferenceRepository())
	router := NewRouter(NewReferenceHandler(service), NewAssetHandler())
	request := httptest.NewRequest(http.MethodPost, "/api/v1/reference-data", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, recorder.Code)
	}
	if recorder.Header().Get("Allow") != http.MethodGet {
		t.Fatalf("expected Allow header to be GET")
	}
}
