package httpapi

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"new-website-lelang/internal/domain/award"
)

type AwardHandler struct {
	service *award.Service
}

func NewAwardHandler(service *award.Service) *AwardHandler {
	return &AwardHandler{service: service}
}

type awardListResponse struct {
	Status string          `json:"status"`
	Data   []awardResponse `json:"data"`
}

type awardResponse struct {
	ID        string     `json:"id"`
	ImageSrc  *string    `json:"img_src"`
	InputDate *time.Time `json:"input_date"`
	Sequence  *int64     `json:"seq"`
	Status    *string    `json:"status"`
	CreatedBy *string    `json:"created_by"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedBy *string    `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	IsDeleted int        `json:"is_deleted"`
	DeletedBy *string    `json:"deleted_by"`
	DeletedAt *time.Time `json:"deleted_at"`
	FileName  *string    `json:"file_name"`
}

func (h *AwardHandler) GetAll(c *gin.Context) {
	records, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "gagal mengambil data awards")
		return
	}

	data := make([]awardResponse, len(records))
	for i, record := range records {
		data[i] = mapAwardResponse(record)
	}

	c.JSON(http.StatusOK, awardListResponse{Status: "success", Data: data})
}

func mapAwardResponse(record award.Award) awardResponse {
	return awardResponse{
		ID:        record.ID,
		ImageSrc:  record.ImageSrc,
		InputDate: record.InputDate,
		Sequence:  record.Sequence,
		Status:    record.Status,
		CreatedBy: record.CreatedBy,
		CreatedAt: record.CreatedAt,
		UpdatedBy: record.UpdatedBy,
		UpdatedAt: record.UpdatedAt,
		IsDeleted: record.IsDeleted,
		DeletedBy: record.DeletedBy,
		DeletedAt: record.DeletedAt,
		FileName:  record.FileName,
	}
}
