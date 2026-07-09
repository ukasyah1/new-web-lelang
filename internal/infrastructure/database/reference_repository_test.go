package database

import (
	"context"
	"testing"
)

func TestReferenceRepositoryPrepareAndGetAll(t *testing.T) {
	db, err := OpenSQLite(":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	repository := NewReferenceRepository(db)
	if err := repository.Prepare(); err != nil {
		t.Fatalf("prepare repository: %v", err)
	}

	data, err := repository.GetAll(context.Background())
	if err != nil {
		t.Fatalf("get reference data: %v", err)
	}

	if len(data.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(data.Categories))
	}
	if len(data.AssetTypes) != 3 {
		t.Fatalf("expected 3 asset types, got %d", len(data.AssetTypes))
	}

	cities, err := repository.GetCitiesByProvinceID(context.Background(), "uuid-1")
	if err != nil {
		t.Fatalf("get cities by province: %v", err)
	}
	if len(cities) != 2 {
		t.Fatalf("expected 2 cities, got %d", len(cities))
	}
	if cities[0].ProvinceID != "uuid-1" || cities[0].Name != "Jakarta Pusat" {
		t.Fatalf("unexpected first city: %+v", cities[0])
	}
}
