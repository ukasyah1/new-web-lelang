package catalogs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"new-website-lelang/internal/platform/dbutil"
)

type catalogRow struct {
	ID, FileName, FileURL string
	CreatedAt, UpdatedAt  time.Time
}
type CatalogRepository struct {
	db     *gorm.DB
	schema string
}

func NewCatalogRepository(db *gorm.DB, schema string) *CatalogRepository {
	return &CatalogRepository{db: db, schema: strings.ToUpper(strings.TrimSpace(schema))}
}
func (r *CatalogRepository) GetLatest(ctx context.Context) (Catalog, error) {
	var row catalogRow
	result := r.db.WithContext(ctx).Raw("SELECT ID, FILE_NAME, FILE_URL, CREATED_AT, UPDATED_AT FROM "+dbutil.QualifiedTable(r.schema, "M_KATALOG")+" WHERE STATUS = ? ORDER BY UPDATED_AT DESC, CREATED_AT DESC FETCH NEXT 1 ROWS ONLY", "PUBLISHED").Scan(&row)
	if result.Error != nil {
		return Catalog{}, fmt.Errorf("query latest catalog: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return Catalog{}, gorm.ErrRecordNotFound
	}
	return Catalog{ID: row.ID, FileName: row.FileName, FileURL: row.FileURL, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}, nil
}

func (r *CatalogRepository) GetActive(ctx context.Context) (MonthlyCatalog, error) {
	var row struct {
		FileName, FileURL string
		Size              int64
		PublishedAt       time.Time
	}
	result := r.db.WithContext(ctx).Raw("SELECT TITLE AS FILE_NAME, MINIO_PATH AS FILE_URL, NVL(FILE_SIZE, 0) AS SIZE, CREATED_AT AS PUBLISHED_AT FROM " + dbutil.QualifiedTable(r.schema, "MONTHLY_CATALOGS") + " ORDER BY CREATED_AT DESC FETCH NEXT 1 ROWS ONLY").Scan(&row)
	if result.Error != nil {
		return MonthlyCatalog{}, fmt.Errorf("query active monthly catalog: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return MonthlyCatalog{}, gorm.ErrRecordNotFound
	}
	return MonthlyCatalog{FileName: row.FileName, FileURL: row.FileURL, Size: row.Size, PublishedAt: row.PublishedAt}, nil
}
