package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getAllHandler interface {
	GetAll(*gin.Context)
}

type assetHandler interface {
	getAllHandler
	GetByID(*gin.Context)
}

type referenceHandler interface {
	GetAll(*gin.Context)
	GetCitiesByProvince(*gin.Context)
	GetDistrictsByCity(*gin.Context)
	GetProvinces(*gin.Context)
	GetCategories(*gin.Context)
	GetAssetTypes(*gin.Context)
}

func NewRouter(referenceHandler referenceHandler, assetHandler assetHandler, optionalHandlers ...getAllHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.HandleMethodNotAllowed = true

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	api := router.Group("/api/v1")
	api.GET("/reference-data", referenceHandler.GetAll)
	api.GET("/master-data", referenceHandler.GetAll)
	api.GET("/master-data/kota", referenceHandler.GetCitiesByProvince)
	api.GET("/master-data/kecamatan", referenceHandler.GetDistrictsByCity)
	api.GET("/master-data/provinsi", referenceHandler.GetProvinces)
	api.GET("/master-data/kategori", referenceHandler.GetCategories)
	api.GET("/master-data/tipe-aset", referenceHandler.GetAssetTypes)
	api.GET("/assets", assetHandler.GetAll)
	api.GET("/assets/:id", assetHandler.GetByID)
	if len(optionalHandlers) > 0 && optionalHandlers[0] != nil {
		api.GET("/awards", optionalHandlers[0].GetAll)
	}
	if len(optionalHandlers) > 1 && optionalHandlers[1] != nil {
		api.GET("/faqs", optionalHandlers[1].GetAll)
	}
	if len(optionalHandlers) > 2 && optionalHandlers[2] != nil {
		api.GET("/banners", optionalHandlers[2].GetAll)
	}
	if len(optionalHandlers) > 3 && optionalHandlers[3] != nil {
		api.GET("/catalogs/latest", optionalHandlers[3].GetAll)
		if handler, ok := optionalHandlers[3].(interface{ GetActive(*gin.Context) }); ok {
			api.GET("/catalogs/active", handler.GetActive)
		}
	}
	if len(optionalHandlers) > 4 && optionalHandlers[4] != nil {
		api.GET("/schedules", optionalHandlers[4].GetAll)
	}
	if len(optionalHandlers) > 5 && optionalHandlers[5] != nil {
		api.GET("/settings", optionalHandlers[5].GetAll)
	}

	router.NoMethod(func(c *gin.Context) {
		c.Header("Allow", http.MethodGet)
		respondError(c, http.StatusMethodNotAllowed, "method not allowed")
	})

	return router
}
