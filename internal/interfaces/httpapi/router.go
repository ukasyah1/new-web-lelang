package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(referenceHandler *ReferenceHandler, assetHandler *AssetHandler, awardHandlers ...*AwardHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.HandleMethodNotAllowed = true

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	api := router.Group("/api/v1")
	api.GET("/reference-data", referenceHandler.GetAll)
	api.GET("/assets", assetHandler.GetAll)
	if len(awardHandlers) > 0 && awardHandlers[0] != nil {
		api.GET("/awards", awardHandlers[0].GetAll)
	}

	router.NoMethod(func(c *gin.Context) {
		c.Header("Allow", http.MethodGet)
		respondError(c, http.StatusMethodNotAllowed, "method not allowed")
	})

	return router
}
