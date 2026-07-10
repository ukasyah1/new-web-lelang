package settings

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Link struct {
	ID, ParentID, Title, URL string
	Order                    int
}
type Repository interface {
	Get(context.Context, string) (map[string]string, []Link, error)
}
type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) Get(ctx context.Context, lang string) (map[string]string, []Link, error) {
	return s.repo.Get(ctx, lang)
}

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

type linkResponse struct {
	ID       string         `json:"id"`
	Title    string         `json:"title"`
	URL      string         `json:"url"`
	Children []linkResponse `json:"children"`
}

func (h *Handler) GetAll(c *gin.Context) {
	lang := strings.ToLower(strings.TrimSpace(c.DefaultQuery("lang", "id")))
	if lang != "id" && lang != "en" {
		c.JSON(400, gin.H{"status": "error", "message": "lang harus id atau en"})
		return
	}
	values, links, err := h.service.Get(c.Request.Context(), lang)
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": "gagal mengambil konfigurasi sistem"})
		return
	}
	children := map[string][]linkResponse{}
	roots := []linkResponse{}
	for _, l := range links {
		item := linkResponse{ID: l.ID, Title: l.Title, URL: l.URL, Children: []linkResponse{}}
		if l.ParentID == "" {
			roots = append(roots, item)
		} else {
			children[l.ParentID] = append(children[l.ParentID], item)
		}
	}
	for i := range roots {
		roots[i].Children = children[roots[i].ID]
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"site_settings": values, "navbar": roots}})
}
