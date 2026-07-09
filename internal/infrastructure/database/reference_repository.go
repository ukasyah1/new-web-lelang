package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"new-website-lelang/internal/domain/reference"
)

type categoryModel struct {
	ID   string `gorm:"column:ID;primaryKey"`
	Name string `gorm:"column:NAMA_KATEGORI;not null"`
}

func (categoryModel) TableName() string {
	return "M_KATEGORI"
}

type assetTypeModel struct {
	ID         string `gorm:"column:ID;primaryKey"`
	CategoryID string `gorm:"column:KATEGORI_ID;not null;index"`
	Name       string `gorm:"column:NAMA_TIPE;not null"`
}

func (assetTypeModel) TableName() string {
	return "M_TIPE_ASET"
}

type provinceModel struct {
	ID   string `gorm:"column:ID;primaryKey"`
	Name string `gorm:"column:NAMA_PROVINSI;not null"`
}

func (provinceModel) TableName() string {
	return "M_PROVINSI"
}

type cityModel struct {
	ID         string `gorm:"column:ID;primaryKey"`
	ProvinceID string `gorm:"column:PROVINSI_ID;not null;index"`
	Name       string `gorm:"column:NAMA_KOTA;not null"`
	CodePrefix string `gorm:"column:KODE_PREFIX"`
}

func (cityModel) TableName() string {
	return "M_KOTA"
}

type salesMethodModel struct {
	ID   string `gorm:"column:ID;primaryKey"`
	Name string `gorm:"column:NAMA_METODE;not null"`
}

func (salesMethodModel) TableName() string {
	return "M_METODE_PENJUALAN"
}

type kpknlModel struct {
	ID   string `gorm:"column:ID;primaryKey"`
	Code string `gorm:"column:KODE_KPKNL;not null;uniqueIndex"`
	Name string `gorm:"column:NAMA_KANTOR;not null"`
}

func (kpknlModel) TableName() string {
	return "M_KPKNL"
}

type ReferenceRepository struct {
	db *gorm.DB
}

func NewReferenceRepository(db *gorm.DB) *ReferenceRepository {
	return &ReferenceRepository{db: db}
}

// Prepare creates the required tables and inserts starter data once.
func (r *ReferenceRepository) Prepare() error {
	if err := r.db.AutoMigrate(
		&categoryModel{},
		&assetTypeModel{},
		&provinceModel{},
		&cityModel{},
		&salesMethodModel{},
		&kpknlModel{},
	); err != nil {
		return fmt.Errorf("migrate tables: %w", err)
	}

	return r.seed()
}

func (r *ReferenceRepository) seed() error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		data := []any{
			&categoryModel{ID: "uuid-1", Name: "Properti"},
			&categoryModel{ID: "uuid-2", Name: "Kendaraan"},
			&assetTypeModel{ID: "uuid-1", CategoryID: "uuid-1", Name: "Rumah"},
			&assetTypeModel{ID: "uuid-2", CategoryID: "uuid-1", Name: "Ruko"},
			&assetTypeModel{ID: "uuid-3", CategoryID: "uuid-2", Name: "Mobil"},
			&provinceModel{ID: "uuid-1", Name: "DKI Jakarta"},
			&provinceModel{ID: "uuid-2", Name: "Jawa Barat"},
			&cityModel{ID: "uuid-1", ProvinceID: "uuid-1", Name: "Jakarta Pusat", CodePrefix: "JKP"},
			&cityModel{ID: "uuid-2", ProvinceID: "uuid-1", Name: "Jakarta Selatan", CodePrefix: "JKS"},
			&cityModel{ID: "uuid-3", ProvinceID: "uuid-2", Name: "Bandung", CodePrefix: "BDG"},
			&salesMethodModel{ID: "uuid-1", Name: "Lelang"},
			&salesMethodModel{ID: "uuid-2", Name: "Jual Damai"},
			&kpknlModel{ID: "uuid-1", Code: "KPKNL-JKT1", Name: "KPKNL Jakarta I"},
			&kpknlModel{ID: "uuid-2", Code: "KPKNL-JKT2", Name: "KPKNL Jakarta II"},
		}

		for _, item := range data {
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error; err != nil {
				return fmt.Errorf("seed reference data: %w", err)
			}
		}
		return nil
	})
}

func (r *ReferenceRepository) GetAll(ctx context.Context) (reference.Data, error) {
	var categories []categoryModel
	var assetTypes []assetTypeModel
	var provinces []provinceModel
	var salesMethods []salesMethodModel
	var kpknls []kpknlModel

	db := r.db.WithContext(ctx)
	if err := db.Order("id").Find(&categories).Error; err != nil {
		return reference.Data{}, err
	}
	if err := db.Order("id").Find(&assetTypes).Error; err != nil {
		return reference.Data{}, err
	}
	if err := db.Order("id").Find(&provinces).Error; err != nil {
		return reference.Data{}, err
	}
	if err := db.Order("id").Find(&salesMethods).Error; err != nil {
		return reference.Data{}, err
	}
	if err := db.Order("id").Find(&kpknls).Error; err != nil {
		return reference.Data{}, err
	}

	result := reference.Data{
		Categories:   make([]reference.Category, len(categories)),
		AssetTypes:   make([]reference.AssetType, len(assetTypes)),
		Provinces:    make([]reference.Province, len(provinces)),
		SalesMethods: make([]reference.SalesMethod, len(salesMethods)),
		KPKNLs:       make([]reference.KPKNL, len(kpknls)),
	}

	for i, item := range categories {
		result.Categories[i] = reference.Category{ID: item.ID, Name: item.Name}
	}
	for i, item := range assetTypes {
		result.AssetTypes[i] = reference.AssetType{ID: item.ID, CategoryID: item.CategoryID, Name: item.Name}
	}
	for i, item := range provinces {
		result.Provinces[i] = reference.Province{ID: item.ID, Name: item.Name}
	}
	for i, item := range salesMethods {
		result.SalesMethods[i] = reference.SalesMethod{ID: item.ID, Name: item.Name}
	}
	for i, item := range kpknls {
		result.KPKNLs[i] = reference.KPKNL{ID: item.ID, Code: item.Code, Name: item.Name}
	}

	return result, nil
}

func (r *ReferenceRepository) GetCitiesByProvinceID(ctx context.Context, provinceID string) ([]reference.City, error) {
	var cities []cityModel
	if err := r.db.WithContext(ctx).
		Where("PROVINSI_ID = ?", provinceID).
		Order("ID").
		Find(&cities).Error; err != nil {
		return nil, err
	}

	result := make([]reference.City, len(cities))
	for i, item := range cities {
		result[i] = reference.City{ID: item.ID, ProvinceID: item.ProvinceID, Name: item.Name}
	}
	return result, nil
}
