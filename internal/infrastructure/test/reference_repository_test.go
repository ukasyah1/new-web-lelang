package infrastructure_test

import (
	"context"
	"testing"

	"new-website-lelang/internal/infrastructure/database"
)

func TestReferenceRepositoryPrepareAndGetAll(t *testing.T) {
	db, err := database.OpenSQLite(":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	repository := database.NewReferenceRepository(db)
	if err := repository.Prepare(); err != nil {
		t.Fatalf("prepare repository: %v", err)
	}

	data, err := repository.GetAll(context.Background())
	if err != nil {
		t.Fatalf("get reference data: %v", err)
	}
	if len(data.Categories) != 2 || len(data.AssetTypes) != 3 {
		t.Fatalf("unexpected reference data: %+v", data)
	}

	cities, err := repository.GetCitiesByProvinceID(context.Background(), "uuid-1")
	if err != nil {
		t.Fatalf("get cities by province: %v", err)
	}
	if len(cities) != 2 || cities[0].ProvinceID != "uuid-1" {
		t.Fatalf("unexpected cities: %+v", cities)
	}
}
