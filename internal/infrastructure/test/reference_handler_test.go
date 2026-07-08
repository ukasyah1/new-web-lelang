package infrastructure_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReferenceData(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/reference-data", nil)
	recorder := httptest.NewRecorder()
	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Status string `json:"status"`
		Data   struct {
			Categories []struct {
				ID string `json:"id"`
			} `json:"kategori"`
			AssetTypes []struct {
				Name string `json:"nama_tipe"`
			} `json:"tipe_aset"`
			KPKNLs []struct {
				Code string `json:"kode_kpknl"`
			} `json:"kpknl"`
		} `json:"data"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != "success" || len(response.Data.Categories) != 2 {
		t.Fatalf("unexpected response: %+v", response)
	}
	if response.Data.AssetTypes[2].Name != "Mobil" || response.Data.KPKNLs[0].Code != "KPKNL-JKT1" {
		t.Fatalf("unexpected reference values: %+v", response.Data)
	}
}

func TestReferenceDataRejectsUnsupportedMethod(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/v1/reference-data", nil)
	recorder := httptest.NewRecorder()
	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed || recorder.Header().Get("Allow") != http.MethodGet {
		t.Fatalf("unexpected method response: status=%d allow=%q", recorder.Code, recorder.Header().Get("Allow"))
	}
}
