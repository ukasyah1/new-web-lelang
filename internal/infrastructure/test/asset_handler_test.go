package infrastructure_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

	var response struct {
		Status string `json:"status"`
		Meta   struct {
			TotalData  int `json:"total_data"`
			TotalPages int `json:"total_pages"`
		} `json:"meta"`
		Data []struct {
			CollateralCode string `json:"kode_agunan"`
		} `json:"data"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != "success" || response.Meta.TotalData != 1 || response.Meta.TotalPages != 1 {
		t.Fatalf("unexpected response: %+v", response)
	}
	if len(response.Data) != 1 || response.Data[0].CollateralCode != "AG-JKT-001" {
		t.Fatalf("unexpected assets: %+v", response.Data)
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
