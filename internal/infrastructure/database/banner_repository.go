package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"new-website-lelang/internal/domain/banner"
)

type bannerRow struct {
	ID         string  `gorm:"column:ID"`
	ImageURL   string  `gorm:"column:IMAGE_URL"`
	TargetURL  *string `gorm:"column:TARGET_URL"`
	OrderIndex int     `gorm:"column:ORDER_INDEX"`
}

type BannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) *BannerRepository {
	return &BannerRepository{db: db}
}

func (r *BannerRepository) GetActive(ctx context.Context) ([]banner.Banner, error) {
	rows := []bannerRow{}
	result := r.db.WithContext(ctx).Raw(`
		SELECT ID,
		       IMAGE_URL,
		       TARGET_URL,
		       ORDER_INDEX
		FROM CMS.BANNER_SLIDER_PROMOSI
		WHERE IS_ACTIVE = ?
		ORDER BY ORDER_INDEX ASC`, 1).Scan(&rows)
	if result.Error != nil {
		return nil, fmt.Errorf("query CMS.BANNER_SLIDER_PROMOSI: %w", result.Error)
	}

	banners := make([]banner.Banner, len(rows))
	for i, row := range rows {
		banners[i] = banner.Banner{
			ID:         row.ID,
			ImageURL:   row.ImageURL,
			TargetURL:  row.TargetURL,
			OrderIndex: row.OrderIndex,
		}
	}
	return banners, nil
}
