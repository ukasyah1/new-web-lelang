package assets

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct{ service *Service }

func NewAssetHandler(services ...*Service) *AssetHandler {
	var service *Service
	if len(services) > 0 {
		service = services[0]
	}
	return &AssetHandler{service: service}
}

type assetSearchQuery struct {
	Search         string   `form:"search"`
	CategoryID     string   `form:"kategori_id"`
	AssetTypeID    string   `form:"tipe_aset_id"`
	ProvinceID     string   `form:"provinsi_id"`
	CityID         string   `form:"kota_id"`
	District       string   `form:"kecamatan"`
	TagID          string   `form:"tag_id"`
	SalesMethodIDs []string `form:"metode_penjualan_id[]"`
	MinimumPrice   *int64   `form:"harga_min" binding:"omitempty,min=0"`
	MaximumPrice   *int64   `form:"harga_max" binding:"omitempty,min=0"`
	Page           int      `form:"page" binding:"omitempty,min=1"`
	Limit          int      `form:"limit" binding:"omitempty,min=1,max=100"`
}

type assetListResponse struct {
	Status string            `json:"status"`
	Meta   assetMetaResponse `json:"meta"`
	Data   []assetResponse   `json:"data"`
}
type assetMetaResponse struct {
	TotalData   int64 `json:"total_data"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
}
type namedReferenceResponse struct {
	ID   string `json:"id"`
	Name string `json:"nama"`
}
type kpknlAssetResponse struct {
	ID   string `json:"id"`
	Code string `json:"kode"`
	Name string `json:"nama"`
}
type auctionEventResponse struct {
	EventID       string             `json:"event_id"`
	KPKNL         kpknlAssetResponse `json:"kpknl"`
	StartDate     *time.Time         `json:"start_date"`
	EndDate       *time.Time         `json:"end_date"`
	Timezone      string             `json:"zona_waktu"`
	AuctionStatus string             `json:"status_lelang"`
}
type assetResponse struct {
	ID             string                   `json:"id"`
	CollateralCode string                   `json:"kode_agunan"`
	Name           string                   `json:"nama_aset"`
	Category       namedReferenceResponse   `json:"kategori"`
	AssetType      namedReferenceResponse   `json:"tipe_aset"`
	SalesMethods   []namedReferenceResponse `json:"metode_penjualan"`
	Tags           []namedReferenceResponse `json:"tags"`
	Province       namedReferenceResponse   `json:"provinsi"`
	City           namedReferenceResponse   `json:"kota"`
	AuctionPrice   float64                  `json:"harga_lelang"`
	OriginalPrice  float64                  `json:"harga_coret"`
	AuctionEvent   auctionEventResponse     `json:"auction_event"`
	ImageURLs      []string                 `json:"image_urls"`
	LandArea       int64                    `json:"luas_tanah"`
	BuildingArea   int64                    `json:"luas_bangunan"`
	Certificate    string                   `json:"jenis_sertifikat"`
	Facilities     []string                 `json:"fasilitas"`
	ViewCount      int64                    `json:"view_count"`
}

type assetDetailResponse struct {
	ID             string                   `json:"id"`
	CollateralCode string                   `json:"kode_agunan"`
	Name           string                   `json:"nama_aset"`
	Category       namedReferenceResponse   `json:"kategori"`
	AssetType      namedReferenceResponse   `json:"tipe_aset"`
	SalesMethods   []namedReferenceResponse `json:"metode_penjualan"`
	Tags           []namedReferenceResponse `json:"tags"`
	Province       namedReferenceResponse   `json:"provinsi"`
	City           namedReferenceResponse   `json:"kota"`
	Address        string                   `json:"alamat"`
	AuctionPrice   float64                  `json:"harga_lelang"`
	OriginalPrice  float64                  `json:"harga_coret"`
	LandArea       int64                    `json:"luas_tanah"`
	BuildingArea   int64                    `json:"luas_bangunan"`
	Certificate    string                   `json:"jenis_sertifikat"`
	Coordinates    string                   `json:"koordinat"`
	ImageURLs      []string                 `json:"image_urls"`
	Description    string                   `json:"deskripsi"`
	Facilities     []string                 `json:"fasilitas"`
	AuctionEvent   auctionEventResponse     `json:"auction_event"`
	ViewCount      int64                    `json:"view_count"`
	UpdatedAt      *time.Time               `json:"updated_at"`
}

type assetDetailAPIResponse struct {
	Status string              `json:"status"`
	Data   assetDetailResponse `json:"data"`
}

func (h *AssetHandler) GetAll(c *gin.Context) {
	var query assetSearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondError(c, http.StatusBadRequest, "query parameter tidak valid")
		return
	}
	if query.MinimumPrice != nil && query.MaximumPrice != nil && *query.MinimumPrice > *query.MaximumPrice {
		respondError(c, http.StatusBadRequest, "harga_min tidak boleh lebih besar dari harga_max")
		return
	}
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}
	if h.service == nil {
		respondError(c, http.StatusInternalServerError, "service assets belum dikonfigurasi")
		return
	}
	result, err := h.service.Search(c.Request.Context(), SearchQuery{Search: query.Search, CategoryID: query.CategoryID, AssetTypeID: query.AssetTypeID, ProvinceID: query.ProvinceID, CityID: query.CityID, District: query.District, MinimumPrice: query.MinimumPrice, MaximumPrice: query.MaximumPrice, Page: query.Page, Limit: query.Limit})
	if err != nil {
		respondError(c, http.StatusInternalServerError, "gagal mengambil data assets")
		return
	}
	totalPages := (result.Total + int64(query.Limit) - 1) / int64(query.Limit)
	data := make([]assetResponse, len(result.Assets))
	for i, a := range result.Assets {
		data[i] = mapAssetResponse(a)
	}
	c.JSON(http.StatusOK, assetListResponse{Status: "success", Meta: assetMetaResponse{result.Total, query.Page, totalPages}, Data: data})
}

func (h *AssetHandler) GetByID(c *gin.Context) {
	if h.service == nil {
		respondError(c, http.StatusInternalServerError, "service assets belum dikonfigurasi")
		return
	}

	asset, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if errors.Is(err, ErrNotFound) {
		respondError(c, http.StatusNotFound, "asset tidak ditemukan")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "gagal mengambil detail asset")
		return
	}

	c.JSON(http.StatusOK, assetDetailAPIResponse{Status: "success", Data: mapAssetDetailResponse(asset)})
}

func mapAssetResponse(a Asset) assetResponse {
	return assetResponse{ID: a.ID, CollateralCode: a.Code, Name: a.Name, Category: namedReferenceResponse{a.Category.ID, a.Category.Name}, AssetType: namedReferenceResponse{a.AssetType.ID, a.AssetType.Name}, SalesMethods: []namedReferenceResponse{}, Tags: []namedReferenceResponse{}, Province: namedReferenceResponse{a.Province.ID, a.Province.Name}, City: namedReferenceResponse{a.City.ID, a.City.Name}, AuctionPrice: a.AuctionPrice, OriginalPrice: a.OriginalPrice, AuctionEvent: auctionEventResponse{EventID: a.ID, KPKNL: kpknlAssetResponse{a.KPKNL.ID, a.KPKNL.Code, a.KPKNL.Name}, StartDate: a.StartDate, EndDate: a.EndDate, Timezone: a.Timezone, AuctionStatus: a.Status}, ImageURLs: a.ImageURLs, LandArea: a.LandArea, BuildingArea: a.BuildingArea, Certificate: a.Certificate, Facilities: a.Facilities, ViewCount: a.ViewCount}
}

func mapAssetDetailResponse(a Asset) assetDetailResponse {
	return assetDetailResponse{ID: a.ID, CollateralCode: a.Code, Name: a.Name, Category: namedReferenceResponse{a.Category.ID, a.Category.Name}, AssetType: namedReferenceResponse{a.AssetType.ID, a.AssetType.Name}, SalesMethods: []namedReferenceResponse{}, Tags: []namedReferenceResponse{}, Province: namedReferenceResponse{a.Province.ID, a.Province.Name}, City: namedReferenceResponse{a.City.ID, a.City.Name}, Address: a.Address, AuctionPrice: a.AuctionPrice, OriginalPrice: a.OriginalPrice, LandArea: a.LandArea, BuildingArea: a.BuildingArea, Certificate: a.Certificate, Coordinates: a.Coordinates, ImageURLs: a.ImageURLs, Description: a.Description, Facilities: a.Facilities, AuctionEvent: auctionEventResponse{EventID: a.ID, KPKNL: kpknlAssetResponse{a.KPKNL.ID, a.KPKNL.Code, a.KPKNL.Name}, StartDate: a.StartDate, EndDate: a.EndDate, Timezone: a.Timezone, AuctionStatus: a.Status}, ViewCount: a.ViewCount, UpdatedAt: a.UpdatedAt}
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func respondError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, errorResponse{Status: "error", Message: message})
}
