package faq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeFAQRepository struct {
	lang       string
	categories []Category
	err        error
}

func (r *fakeFAQRepository) GetAll(_ context.Context, lang string) ([]Category, error) {
	r.lang = lang
	return r.categories, r.err
}

func TestGetFAQs(t *testing.T) {
	repository := &fakeFAQRepository{categories: []Category{
		{
			ID:   "uuid-cat-1",
			Name: "Informasi Umum",
			FAQs: []FAQ{
				{ID: "uuid-1", Question: "Apa itu Website Lelang Bank Jakarta?", Answer: "Website Lelang Bank Jakarta merupakan media informasi..."},
			},
		},
	}}
	service := NewService(repository)
	handler := NewFAQHandler(service)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/faqs?lang=id", nil)
	recorder := httptest.NewRecorder()

	router := gin.New()
	router.GET("/api/v1/faqs", handler.GetAll)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if repository.lang != "id" {
		t.Fatalf("expected lang id, got %q", repository.lang)
	}

	var response faqListResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != "success" || len(response.Data) != 1 {
		t.Fatalf("unexpected response: %+v", response)
	}
	if response.Data[0].CategoryID != "uuid-cat-1" || response.Data[0].FAQs[0].ID != "uuid-1" {
		t.Fatalf("unexpected faq response: %+v", response.Data[0])
	}
}

func TestGetFAQsRejectsUnsupportedLanguage(t *testing.T) {
	service := NewService(&fakeFAQRepository{})
	handler := NewFAQHandler(service)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/faqs?lang=jp", nil)
	recorder := httptest.NewRecorder()

	router := gin.New()
	router.GET("/api/v1/faqs", handler.GetAll)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}
