package catalogs

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct{ service *Service }

func NewCatalogHandler(service *Service) *CatalogHandler { return &CatalogHandler{service: service} }

func (h *CatalogHandler) GetAll(c *gin.Context) {
	catalog, err := h.service.GetLatest(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "katalog tidak ditemukan"})
		return
	}
	title := strings.TrimSuffix(catalog.FileName, filepath.Ext(catalog.FileName))
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{
		"id": catalog.ID, "title": title, "month": catalog.CreatedAt.Month().String(), "year": catalog.CreatedAt.Format("2006"),
		"latest_update": catalog.UpdatedAt, "download_url": catalog.FileURL,
	}})
}

func (h *CatalogHandler) GetActive(c *gin.Context) {
	catalog, err := h.service.GetActive(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "katalog tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"file_name": catalog.FileName, "file_url": catalog.FileURL, "size_bytes": catalog.Size, "published_at": catalog.PublishedAt}})
}
