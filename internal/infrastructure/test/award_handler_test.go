package infrastructure_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"new-website-lelang/internal/domain/award"
)

type fakeAwardRepository struct {
	records []award.Award
}

func (r fakeAwardRepository) GetAll(_ context.Context) ([]award.Award, error) {
	return r.records, nil
}

func TestGetAwards(t *testing.T) {
	fileName := "award.png"
	service := award.NewService(fakeAwardRepository{records: []award.Award{
		{ID: "uuid-award-1", FileName: &fileName},
	}})
	handler := award.NewAwardHandler(service)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/awards", nil)
	recorder := httptest.NewRecorder()

	router := gin.New()
	router.GET("/api/v1/awards", handler.GetAll)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response struct {
		Status string `json:"status"`
		Data   []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != "success" || len(response.Data) != 1 || response.Data[0].ID != "uuid-award-1" {
		t.Fatalf("unexpected response: %+v", response)
	}
}
