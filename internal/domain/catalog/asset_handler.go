package catalog

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const hardcodedAssetTotal = 45

type AssetHandler struct{}

func NewAssetHandler() *AssetHandler {
	return &AssetHandler{}
}

// assetSearchQuery documents every filter accepted by the public endpoint.
// The values are parsed now and can later be forwarded to a database service.
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
	TotalData   int `json:"total_data"`
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
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
	StartDate     string             `json:"start_date"`
	EndDate       string             `json:"end_date"`
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
	AuctionPrice   int64                    `json:"harga_lelang"`
	OriginalPrice  int64                    `json:"harga_coret"`
	AuctionEvent   auctionEventResponse     `json:"auction_event"`
	ImageURLs      []string                 `json:"image_urls"`
	LandArea       int                      `json:"luas_tanah"`
	BuildingArea   int                      `json:"luas_bangunan"`
	Certificate    string                   `json:"jenis_sertifikat"`
	Facilities     []string                 `json:"fasilitas"`
	ViewCount      int                      `json:"view_count"`
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

	totalPages := (hardcodedAssetTotal + query.Limit - 1) / query.Limit
	c.JSON(http.StatusOK, assetListResponse{
		Status: "success",
		Meta: assetMetaResponse{
			TotalData:   hardcodedAssetTotal,
			CurrentPage: query.Page,
			TotalPages:  totalPages,
		},
		Data: hardcodedAssets(),
	})
}

func hardcodedAssets() []assetResponse {
	return []assetResponse{
		{
			ID:             "uuid-101",
			CollateralCode: "AG-JKT-001",
			Name:           "Rumah Strategis Tebet",
			Category:       namedReferenceResponse{ID: "uuid-1", Name: "Properti"},
			AssetType:      namedReferenceResponse{ID: "uuid-1", Name: "Rumah"},
			SalesMethods: []namedReferenceResponse{
				{ID: "uuid-1", Name: "E-Auction"},
				{ID: "uuid-2", Name: "Jual Beli"},
			},
			Tags: []namedReferenceResponse{
				{ID: "uuid-1", Name: "Promo"},
			},
			Province:      namedReferenceResponse{ID: "uuid-1", Name: "DKI Jakarta"},
			City:          namedReferenceResponse{ID: "uuid-1", Name: "Jakarta Selatan"},
			AuctionPrice:  450000000,
			OriginalPrice: 550000000,
			AuctionEvent: auctionEventResponse{
				EventID: "uuid-1",
				KPKNL: kpknlAssetResponse{
					ID:   "uuid-1",
					Code: "KPKNL-JKT1",
					Name: "KPKNL Jakarta I",
				},
				StartDate:     "2026-12-12T09:00:00Z",
				EndDate:       "2026-12-12T12:00:00Z",
				Timezone:      "WIB",
				AuctionStatus: "MENUNGGU",
			},
			ImageURLs: []string{
				"https://api.lelang.com/api/cms/images/550e8400-e29b-41d4-a716-446655440000",
				"https://api.lelang.com/api/cms/images/550e8400-e29b-41d4-a716-446655440001",
			},
			LandArea:     120,
			BuildingArea: 90,
			Certificate:  "Hak Milik (SHM)",
			Facilities:   []string{"2 Kamar Tidur", "Carport", "Dekat Stasiun"},
			ViewCount:    703,
		},
	}
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func respondError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, errorResponse{Status: "error", Message: message})
}
