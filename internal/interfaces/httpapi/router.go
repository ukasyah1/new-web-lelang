package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getAllHandler interface {
	GetAll(*gin.Context)
}

type referenceHandler interface {
	GetAll(*gin.Context)
	GetCitiesByProvince(*gin.Context)
}

func NewRouter(referenceHandler referenceHandler, assetHandler getAllHandler, optionalHandlers ...getAllHandler) *gin.Engine {
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
	api.GET("/assets", assetHandler.GetAll)
	if len(optionalHandlers) > 0 && optionalHandlers[0] != nil {
		api.GET("/awards", optionalHandlers[0].GetAll)
	}
	if len(optionalHandlers) > 1 && optionalHandlers[1] != nil {
		api.GET("/faqs", optionalHandlers[1].GetAll)
	}
	if len(optionalHandlers) > 2 && optionalHandlers[2] != nil {
		api.GET("/banners", optionalHandlers[2].GetAll)
	}

	router.NoMethod(func(c *gin.Context) {
		c.Header("Allow", http.MethodGet)
		respondError(c, http.StatusMethodNotAllowed, "method not allowed")
	})

	return router
}
