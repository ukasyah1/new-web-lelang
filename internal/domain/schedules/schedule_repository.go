package schedules

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"new-website-lelang/internal/platform/dbutil"
	"strings"
	"time"
)

type scheduleRow struct {
	AssetID, Timezone, KPKNLID, KPKNLName, Address, AuctionLink string
	AuctionDate                                                 *time.Time
}
type contactRow struct{ AssetID, Name, Phone string }
type ScheduleRepository struct {
	db     *gorm.DB
	schema string
}

func NewScheduleRepository(db *gorm.DB, schema string) *ScheduleRepository {
	return &ScheduleRepository{db: db, schema: strings.ToUpper(strings.TrimSpace(schema))}
}
func (r *ScheduleRepository) Search(ctx context.Context, q Query) (Result, error) {
	from := " FROM " + dbutil.QualifiedTable(r.schema, "ASSETS") + " a JOIN " + dbutil.QualifiedTable(r.schema, "M_KPKNL") + " k ON k.ID=a.KPKNL_ID"
	parts := []string{"a.START_DATE IS NOT NULL", "a.END_DATE IS NOT NULL"}
	args := []any{}
	if q.KPKNLID != "" {
		parts = append(parts, "a.KPKNL_ID = ?")
		args = append(args, q.KPKNLID)
	}
	if q.StartDate != nil {
		parts = append(parts, "a.START_DATE >= ?")
		args = append(args, *q.StartDate)
	}
	if q.EndDate != nil {
		parts = append(parts, "a.START_DATE < ?")
		args = append(args, q.EndDate.AddDate(0, 0, 1))
	}
	where := " WHERE " + strings.Join(parts, " AND ")
	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(*)"+from+where, args...).Scan(&total).Error; err != nil {
		return Result{}, fmt.Errorf("count schedules: %w", err)
	}
	rows := []scheduleRow{}
	pageArgs := append(args, (q.Page-1)*q.Limit, q.Limit)
	sql := "SELECT a.ID AS ASSET_ID,a.START_DATE AS AUCTION_DATE,a.ZONA_WAKTU AS TIMEZONE,k.ID AS KPKNL_ID,k.NAMA_KANTOR AS KPKNL_NAME,a.ALAMAT AS ADDRESS,a.LINK_LELANG AS AUCTION_LINK" + from + where + " ORDER BY a.START_DATE ASC OFFSET ? ROWS FETCH NEXT ? ROWS ONLY"
	if err := r.db.WithContext(ctx).Raw(sql, pageArgs...).Scan(&rows).Error; err != nil {
		return Result{}, fmt.Errorf("query schedules: %w", err)
	}
	result := make([]Schedule, len(rows))
	ids := make([]string, len(rows))
	for i, row := range rows {
		ids[i] = row.AssetID
		result[i] = Schedule{AssetID: row.AssetID, AuctionDate: row.AuctionDate, Timezone: row.Timezone, KPKNL: KPKNL{ID: row.KPKNLID, Name: row.KPKNLName}, Address: row.Address, AuctionLink: row.AuctionLink, Contacts: []Contact{}}
	}
	if len(ids) > 0 {
		contacts := []contactRow{}
		contactSQL := "SELECT ac.ASSET_ID,p.NAMA AS NAME,p.NO_HP AS PHONE FROM " + dbutil.QualifiedTable(r.schema, "ASSET_CONTACTS") + " ac JOIN " + dbutil.QualifiedTable(r.schema, "M_PIC") + " p ON p.ID=ac.PIC_ID WHERE ac.ASSET_ID IN ? ORDER BY ac.ASSET_ID,ac.URUTAN"
		if err := r.db.WithContext(ctx).Raw(contactSQL, ids).Scan(&contacts).Error; err != nil {
			return Result{}, fmt.Errorf("query schedule contacts: %w", err)
		}
		index := map[string]int{}
		for i, item := range result {
			index[item.AssetID] = i
		}
		for _, item := range contacts {
			result[index[item.AssetID]].Contacts = append(result[index[item.AssetID]].Contacts, Contact{Name: item.Name, Phone: item.Phone})
		}
	}
	return Result{Total: total, Schedules: result}, nil
}
