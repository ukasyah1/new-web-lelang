package masterdata

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ReferenceHandler struct {
	service *Service
}

func NewReferenceHandler(service *Service) *ReferenceHandler {
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

type cityResponse struct {
	ID         string `json:"id"`
	ProvinceID string `json:"provinsi_id"`
	Name       string `json:"nama_kota"`
}

type districtResponse struct {
	ID     string `json:"id"`
	CityID string `json:"kota_id"`
	Name   string `json:"nama_kecamatan"`
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
	Districts    []districtResponse    `json:"kecamatan"`
	SalesMethods []salesMethodResponse `json:"metode_penjualan"`
	KPKNLs       []kpknlResponse       `json:"kpknl"`
}

type referenceResponse struct {
	Status string                `json:"status"`
	Data   referenceDataResponse `json:"data"`
}

type citiesResponse struct {
	Status string         `json:"status"`
	Data   []cityResponse `json:"data"`
}

func (h *ReferenceHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll(c.Request.Context(), strings.TrimSpace(c.Query("value")))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, referenceResponse{
		Status: "success",
		Data:   mapReferenceData(data),
	})
}

func (h *ReferenceHandler) GetCitiesByProvince(c *gin.Context) {
	provinceID := strings.TrimSpace(c.Query("provinsi_id"))
	if provinceID == "" {
		respondError(c, http.StatusBadRequest, "provinsi_id wajib diisi")
		return
	}

	cities, err := h.service.GetCitiesByProvinceID(c.Request.Context(), provinceID, strings.TrimSpace(c.Query("value")))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, citiesResponse{
		Status: "success",
		Data:   mapCities(cities),
	})
}

func (h *ReferenceHandler) GetDistrictsByCity(c *gin.Context) {
	cityID := strings.TrimSpace(c.Query("kota_id"))
	if cityID == "" {
		respondError(c, http.StatusBadRequest, "kota_id wajib diisi")
		return
	}
	districts, err := h.service.GetDistrictsByCityID(c.Request.Context(), cityID, strings.TrimSpace(c.Query("value")))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal server error")
		return
	}
	result := make([]districtResponse, len(districts))
	for i, item := range districts {
		result[i] = districtResponse{ID: item.ID, CityID: item.CityID, Name: item.Name}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}

func (h *ReferenceHandler) GetProvinces(c *gin.Context) {
	provinces, err := h.service.GetProvinces(c.Request.Context(), strings.TrimSpace(c.Query("value")))
	if err != nil { respondError(c, http.StatusInternalServerError, "internal server error"); return }
	result := make([]provinceResponse, len(provinces))
	for i, item := range provinces { result[i] = provinceResponse{ID: item.ID, Name: item.Name} }
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}
func (h *ReferenceHandler) GetCategories(c *gin.Context) { data, err := h.service.GetCategories(c.Request.Context(), strings.TrimSpace(c.Query("value"))); if err != nil { respondError(c, 500, "internal server error"); return }; c.JSON(200, gin.H{"status":"success","data":data}) }
func (h *ReferenceHandler) GetAssetTypes(c *gin.Context) { data, err := h.service.GetAssetTypes(c.Request.Context(), strings.TrimSpace(c.Query("value"))); if err != nil { respondError(c, 500, "internal server error"); return }; c.JSON(200, gin.H{"status":"success","data":data}) }

func mapReferenceData(data Data) referenceDataResponse {
	result := referenceDataResponse{
		Categories:   make([]categoryResponse, len(data.Categories)),
		AssetTypes:   make([]assetTypeResponse, len(data.AssetTypes)),
		Provinces:    make([]provinceResponse, len(data.Provinces)),
		Districts:    make([]districtResponse, len(data.Districts)),
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
	for i, item := range data.Districts {
		result.Districts[i] = districtResponse{ID: item.ID, CityID: item.CityID, Name: item.Name}
	}
	for i, item := range data.SalesMethods {
		result.SalesMethods[i] = salesMethodResponse{ID: item.ID, Name: item.Name}
	}
	for i, item := range data.KPKNLs {
		result.KPKNLs[i] = kpknlResponse{ID: item.ID, Code: item.Code, Name: item.Name}
	}

	return result
}

func mapCities(cities []City) []cityResponse {
	result := make([]cityResponse, len(cities))
	for i, item := range cities {
		result[i] = cityResponse{ID: item.ID, ProvinceID: item.ProvinceID, Name: item.Name}
	}
	return result
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func respondError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, errorResponse{Status: "error", Message: message})
}
