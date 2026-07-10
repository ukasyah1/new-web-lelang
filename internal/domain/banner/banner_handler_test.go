package banner_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"new-website-lelang/internal/domain/assets"
	"new-website-lelang/internal/domain/banner"
	reference "new-website-lelang/internal/domain/masterdata"
	"new-website-lelang/internal/infrastructure/memory"
	"new-website-lelang/internal/interfaces/httpapi"
)

func TestGetActiveBanners(t *testing.T) {
	referenceService := reference.NewService(memory.NewReferenceRepository())
	bannerService := banner.NewService(memory.NewBannerRepository())
	router := httpapi.NewRouter(
		reference.NewReferenceHandler(referenceService),
		assets.NewAssetHandler(),
		nil,
		nil,
		banner.NewBannerHandler(bannerService),
	)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/banners", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Status string `json:"status"`
		Data   []struct {
			ID         string  `json:"id"`
			ImageURL   string  `json:"image_url"`
			TargetURL  *string `json:"target_url"`
			OrderIndex int     `json:"order_index"`
		} `json:"data"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != "success" || len(response.Data) != 2 {
		t.Fatalf("unexpected banner response: %+v", response)
	}
	if response.Data[0].ID != "uuid-1" ||
		response.Data[0].ImageURL != "https://api.lelang.com/api/cms/images/banner1_uuid" ||
		response.Data[0].TargetURL == nil ||
		*response.Data[0].TargetURL != "/daftar-agunan" ||
		response.Data[0].OrderIndex != 1 {
		t.Fatalf("unexpected first banner: %+v", response.Data[0])
	}
	if response.Data[1].TargetURL != nil {
		t.Fatalf("expected second banner target_url to be null, got %q", *response.Data[1].TargetURL)
	}
}
