package faq

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FAQHandler struct {
	service *Service
}

func NewFAQHandler(service *Service) *FAQHandler {
	return &FAQHandler{service: service}
}

type faqListResponse struct {
	Status string                `json:"status"`
	Data   []faqCategoryResponse `json:"data"`
}

type faqCategoryResponse struct {
	CategoryID   string        `json:"category_id"`
	CategoryName string        `json:"category_name"`
	FAQs         []faqResponse `json:"faqs"`
}

type faqResponse struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (h *FAQHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll(c.Request.Context(), c.Query("lang"))
	if err != nil {
		if errors.Is(err, ErrUnsupportedLanguage) {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "gagal mengambil data faqs")
		return
	}

	c.JSON(http.StatusOK, faqListResponse{
		Status: "success",
		Data:   mapFAQCategories(categories),
	})
}

func mapFAQCategories(categories []Category) []faqCategoryResponse {
	response := make([]faqCategoryResponse, len(categories))
	for i, category := range categories {
		faqs := make([]faqResponse, len(category.FAQs))
		for j, item := range category.FAQs {
			faqs[j] = faqResponse{
				ID:       item.ID,
				Question: item.Question,
				Answer:   item.Answer,
			}
		}
		response[i] = faqCategoryResponse{
			CategoryID:   category.ID,
			CategoryName: category.Name,
			FAQs:         faqs,
		}
	}
	return response
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func respondError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, errorResponse{Status: "error", Message: message})
}
