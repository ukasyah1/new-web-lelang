package schedules

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

type requestQuery struct {
	KPKNLID   string `form:"kpknl_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (h *Handler) GetAll(c *gin.Context) {
	var q requestQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "query parameter tidak valid"})
		return
	}
	parse := func(v string) (*time.Time, error) {
		if v == "" {
			return nil, nil
		}
		d, e := time.Parse("2006-01-02", v)
		return &d, e
	}
	start, e := parse(q.StartDate)
	if e != nil {
		c.JSON(400, gin.H{"status": "error", "message": "start_date harus berformat YYYY-MM-DD"})
		return
	}
	end, e := parse(q.EndDate)
	if e != nil {
		c.JSON(400, gin.H{"status": "error", "message": "end_date harus berformat YYYY-MM-DD"})
		return
	}
	if start != nil && end != nil && start.After(*end) {
		c.JSON(400, gin.H{"status": "error", "message": "start_date tidak boleh melebihi end_date"})
		return
	}
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Limit == 0 {
		q.Limit = 10
	}
	result, err := h.service.Search(c.Request.Context(), Query{KPKNLID: q.KPKNLID, StartDate: start, EndDate: end, Page: q.Page, Limit: q.Limit})
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": "gagal mengambil jadwal lelang"})
		return
	}
	data := make([]gin.H, len(result.Schedules))
	for i, s := range result.Schedules {
		contacts := make([]gin.H, len(s.Contacts))
		for j, x := range s.Contacts {
			contacts[j] = gin.H{"nama": x.Name, "no_hp": x.Phone}
		}
		data[i] = gin.H{"asset_id": s.AssetID, "tanggal_lelang": s.AuctionDate, "zona_waktu": s.Timezone, "kpknl": gin.H{"id": s.KPKNL.ID, "nama_kpknl": s.KPKNL.Name}, "alamat_aset": s.Address, "contacts": contacts, "link_lelang": s.AuctionLink}
	}
	totalPages := (result.Total + int64(q.Limit) - 1) / int64(q.Limit)
	c.JSON(http.StatusOK, gin.H{"status": "success", "meta": gin.H{"total_data": result.Total, "current_page": q.Page, "total_pages": totalPages}, "data": data})
}
