package banner

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BannerHandler struct {
	service *Service
}

func NewBannerHandler(service *Service) *BannerHandler {
	return &BannerHandler{service: service}
}

type bannerListResponse struct {
	Status string           `json:"status"`
	Data   []bannerResponse `json:"data"`
}

type bannerResponse struct {
	ID         string  `json:"id"`
	ImageURL   string  `json:"image_url"`
	TargetURL  *string `json:"target_url"`
	OrderIndex int     `json:"order_index"`
}

func (h *BannerHandler) GetAll(c *gin.Context) {
	banners, err := h.service.GetActive(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "gagal mengambil data banners")
		return
	}

	c.JSON(http.StatusOK, bannerListResponse{
		Status: "success",
		Data:   mapBanners(banners),
	})
}

func mapBanners(banners []Banner) []bannerResponse {
	response := make([]bannerResponse, len(banners))
	for i, item := range banners {
		response[i] = bannerResponse{
			ID:         item.ID,
			ImageURL:   item.ImageURL,
			TargetURL:  item.TargetURL,
			OrderIndex: item.OrderIndex,
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
