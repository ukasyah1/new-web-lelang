package database

import (
	"context"
	"testing"
)

func TestBannerRepositoryGetActive(t *testing.T) {
	db, err := OpenSQLite(":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Exec("ATTACH DATABASE ':memory:' AS CMS").Error; err != nil {
		t.Fatalf("attach CMS schema: %v", err)
	}
	if err := db.Exec(`
		CREATE TABLE CMS.BANNER_SLIDER_PROMOSI (
			ID TEXT PRIMARY KEY,
			IMAGE_URL TEXT NOT NULL,
			TARGET_URL TEXT,
			ORDER_INDEX INTEGER NOT NULL,
			IS_ACTIVE INTEGER
		)`).Error; err != nil {
		t.Fatalf("create banner table: %v", err)
	}
	if err := db.Exec(`
		INSERT INTO CMS.BANNER_SLIDER_PROMOSI (ID, IMAGE_URL, TARGET_URL, ORDER_INDEX, IS_ACTIVE)
		VALUES
			('uuid-2', 'image-2', NULL, 2, 1),
			('uuid-1', 'image-1', '/target-1', 1, 1),
			('uuid-3', 'image-3', '/target-3', 3, 0)`).Error; err != nil {
		t.Fatalf("seed banner table: %v", err)
	}

	banners, err := NewBannerRepository(db).GetActive(context.Background())
	if err != nil {
		t.Fatalf("get active banners: %v", err)
	}

	if len(banners) != 2 {
		t.Fatalf("expected 2 active banners, got %d", len(banners))
	}
	if banners[0].ID != "uuid-1" || banners[0].OrderIndex != 1 {
		t.Fatalf("expected banners ordered by order_index, got %+v", banners)
	}
	if banners[1].TargetURL != nil {
		t.Fatalf("expected nullable target URL, got %q", *banners[1].TargetURL)
	}
}
