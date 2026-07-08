package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"new-website-lelang/internal/domain/reference"
)

type ReferenceHandler struct {
	service *reference.Service
}

func NewReferenceHandler(service *reference.Service) *ReferenceHandler {
	return &ReferenceHandler{service: service}
}

type categoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"nama_kategori"`
}

type assetTypeResponse struct {
	ID         string `json:"id"`
	CategoryID string `json:"kategori_id"`
	Name       string `json:"nama_tipe"`
}

type provinceResponse struct {
	ID   string `json:"id"`
	Name string `json:"nama_provinsi"`
}

type salesMethodResponse struct {
	ID   string `json:"id"`
	Name string `json:"nama_metode"`
}

type kpknlResponse struct {
	ID   string `json:"id"`
	Code string `json:"kode_kpknl"`
	Name string `json:"nama_kpknl"`
}

type referenceDataResponse struct {
	Categories   []categoryResponse    `json:"kategori"`
	AssetTypes   []assetTypeResponse   `json:"tipe_aset"`
	Provinces    []provinceResponse    `json:"provinsi"`
	SalesMethods []salesMethodResponse `json:"metode_penjualan"`
	KPKNLs       []kpknlResponse       `json:"kpknl"`
}

type referenceResponse struct {
	Status string                `json:"status"`
	Data   referenceDataResponse `json:"data"`
}

func (h *ReferenceHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, referenceResponse{
		Status: "success",
		Data:   mapReferenceData(data),
	})
}

func mapReferenceData(data reference.Data) referenceDataResponse {
	result := referenceDataResponse{
		Categories:   make([]categoryResponse, len(data.Categories)),
		AssetTypes:   make([]assetTypeResponse, len(data.AssetTypes)),
		Provinces:    make([]provinceResponse, len(data.Provinces)),
		SalesMethods: make([]salesMethodResponse, len(data.SalesMethods)),
		KPKNLs:       make([]kpknlResponse, len(data.KPKNLs)),
	}

	for i, item := range data.Categories {
		result.Categories[i] = categoryResponse{ID: item.ID, Name: item.Name}
	}
	for i, item := range data.AssetTypes {
		result.AssetTypes[i] = assetTypeResponse{ID: item.ID, CategoryID: item.CategoryID, Name: item.Name}
	}
	for i, item := range data.Provinces {
		result.Provinces[i] = provinceResponse{ID: item.ID, Name: item.Name}
	}
	for i, item := range data.SalesMethods {
		result.SalesMethods[i] = salesMethodResponse{ID: item.ID, Name: item.Name}
	}
	for i, item := range data.KPKNLs {
		result.KPKNLs[i] = kpknlResponse{ID: item.ID, Code: item.Code, Name: item.Name}
	}

	return result
}
