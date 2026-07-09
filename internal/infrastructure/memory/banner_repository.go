package memory

import (
	"context"

	"new-website-lelang/internal/domain/banner"
)

type BannerRepository struct{}

func NewBannerRepository() *BannerRepository {
	return &BannerRepository{}
}

func (r *BannerRepository) GetActive(_ context.Context) ([]banner.Banner, error) {
	firstTarget := "/daftar-agunan"
	return []banner.Banner{
		{
			ID:         "uuid-1",
			ImageURL:   "https://api.lelang.com/api/cms/images/banner1_uuid",
			TargetURL:  &firstTarget,
			OrderIndex: 1,
		},
		{
			ID:         "uuid-2",
			ImageURL:   "https://api.lelang.com/api/cms/images/banner2_uuid",
			TargetURL:  nil,
			OrderIndex: 2,
		},
	}, nil
}
