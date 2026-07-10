package masterdata_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"new-website-lelang/internal/domain/assets"
	reference "new-website-lelang/internal/domain/masterdata"
	"new-website-lelang/internal/infrastructure/memory"
	"new-website-lelang/internal/interfaces/httpapi"
)

func TestGetReferenceData(t *testing.T) {
	service := reference.NewService(memory.NewReferenceRepository())
	router := httpapi.NewRouter(reference.NewReferenceHandler(service), assets.NewAssetHandler())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/reference-data", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

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

func TestGetMasterData(t *testing.T) {
	service := reference.NewService(memory.NewReferenceRepository())
	router := httpapi.NewRouter(reference.NewReferenceHandler(service), assets.NewAssetHandler())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/master-data", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Status string `json:"status"`
		Data   struct {
			Categories []struct {
				Name string `json:"nama_kategori"`
			} `json:"kategori"`
			AssetTypes []struct {
				CategoryID string `json:"kategori_id"`
				Name       string `json:"nama_tipe"`
			} `json:"tipe_aset"`
			Provinces []struct {
				Name string `json:"nama_provinsi"`
			} `json:"provinsi"`
			SalesMethods []struct {
				Name string `json:"nama_metode"`
			} `json:"metode_penjualan"`
			KPKNLs []struct {
				Code string `json:"kode_kpknl"`
				Name string `json:"nama_kpknl"`
			} `json:"kpknl"`
		} `json:"data"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != "success" {
		t.Fatalf("expected success status, got %q", response.Status)
	}
	if response.Data.Categories[0].Name != "Properti" ||
		response.Data.AssetTypes[0].CategoryID != "uuid-1" ||
		response.Data.Provinces[0].Name != "DKI Jakarta" ||
		response.Data.SalesMethods[0].Name != "Lelang" ||
		response.Data.KPKNLs[0].Code != "KPKNL-JKT1" ||
		response.Data.KPKNLs[0].Name != "KPKNL Jakarta I" {
		t.Fatalf("unexpected master data response: %+v", response.Data)
	}
}

func TestGetMasterDataFiltersByValue(t *testing.T) {
	service := reference.NewService(memory.NewReferenceRepository())
	router := httpapi.NewRouter(reference.NewReferenceHandler(service), assets.NewAssetHandler())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/master-data?value=menteng", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Data struct {
			Districts []struct {
				Name string `json:"nama_kecamatan"`
			} `json:"kecamatan"`
		} `json:"data"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(response.Data.Districts) != 1 || response.Data.Districts[0].Name != "Menteng" {
		t.Fatalf("unexpected filtered districts: %+v", response.Data.Districts)
	}
}

func TestGetCitiesByProvince(t *testing.T) {
	service := reference.NewService(memory.NewReferenceRepository())
	router := httpapi.NewRouter(reference.NewReferenceHandler(service), assets.NewAssetHandler())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/master-data/kota?provinsi_id=uuid-1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Status string `json:"status"`
		Data   []struct {
			ID         string `json:"id"`
			ProvinceID string `json:"provinsi_id"`
			Name       string `json:"nama_kota"`
		} `json:"data"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != "success" || len(response.Data) != 2 {
		t.Fatalf("unexpected city response: %+v", response)
	}
	if response.Data[0].ID != "uuid-1" ||
		response.Data[0].ProvinceID != "uuid-1" ||
		response.Data[0].Name != "Jakarta Pusat" {
		t.Fatalf("unexpected first city: %+v", response.Data[0])
	}
}

func TestGetCitiesByProvinceFiltersByValue(t *testing.T) {
	service := reference.NewService(memory.NewReferenceRepository())
	router := httpapi.NewRouter(reference.NewReferenceHandler(service), assets.NewAssetHandler())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/master-data/kota?provinsi_id=uuid-1&value=selatan", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Data []struct {
			Name string `json:"nama_kota"`
		} `json:"data"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(response.Data) != 1 || response.Data[0].Name != "Jakarta Selatan" {
		t.Fatalf("unexpected filtered cities: %+v", response.Data)
	}
}

func TestGetCitiesByProvinceRequiresProvinceID(t *testing.T) {
	service := reference.NewService(memory.NewReferenceRepository())
	router := httpapi.NewRouter(reference.NewReferenceHandler(service), assets.NewAssetHandler())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/master-data/kota", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestReferenceDataRejectsUnsupportedMethod(t *testing.T) {
	service := reference.NewService(memory.NewReferenceRepository())
	router := httpapi.NewRouter(reference.NewReferenceHandler(service), assets.NewAssetHandler())
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
