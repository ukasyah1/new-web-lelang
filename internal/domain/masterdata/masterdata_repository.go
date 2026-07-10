package masterdata

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"new-website-lelang/internal/platform/dbutil"
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

type districtModel struct {
	ID     string `gorm:"column:ID;primaryKey"`
	CityID string `gorm:"column:KOTA_ID;not null;index"`
	Name   string `gorm:"column:NAMA_KECAMATAN;not null"`
}

func (districtModel) TableName() string {
	return "M_KECAMATAN"
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
	db     *gorm.DB
	schema string
}

func NewReferenceRepository(db *gorm.DB) *ReferenceRepository {
	return &ReferenceRepository{db: db}
}

// NewMasterDataRepository uses the schema selected through MIGRATION_SCHEMA.
func NewMasterDataRepository(db *gorm.DB, schema string) *ReferenceRepository {
	return &ReferenceRepository{db: db, schema: strings.ToUpper(strings.TrimSpace(schema))}
}

func (r *ReferenceRepository) table(name string) string {
	return dbutil.QualifiedTable(r.schema, name)
}

// Prepare creates the required tables and inserts starter data once.
func (r *ReferenceRepository) Prepare() error {
	if err := r.db.AutoMigrate(
		&categoryModel{},
		&assetTypeModel{},
		&provinceModel{},
		&cityModel{},
		&districtModel{},
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
			&districtModel{ID: "uuid-1", CityID: "uuid-1", Name: "Menteng"},
			&districtModel{ID: "uuid-2", CityID: "uuid-1", Name: "Tanah Abang"},
			&districtModel{ID: "uuid-3", CityID: "uuid-3", Name: "Coblong"},
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

func (r *ReferenceRepository) GetAll(ctx context.Context, value string) (Data, error) {
	var categories []categoryModel
	var assetTypes []assetTypeModel
	var provinces []provinceModel
	var districts []districtModel
	var salesMethods []salesMethodModel
	var kpknls []kpknlModel

	db := r.db.WithContext(ctx)
	value = strings.TrimSpace(value)
	nameFilter := func(query *gorm.DB, column string) *gorm.DB {
		if value == "" {
			return query
		}
		return query.Where("UPPER("+column+") LIKE ?", "%"+strings.ToUpper(value)+"%")
	}
	if err := nameFilter(db.Table(r.table("M_KATEGORI")), "NAMA_KATEGORI").Order("id").Find(&categories).Error; err != nil {
		return Data{}, err
	}
	if err := nameFilter(db.Table(r.table("M_TIPE_ASET")), "NAMA_TIPE").Order("id").Find(&assetTypes).Error; err != nil {
		return Data{}, err
	}
	if err := nameFilter(db.Table(r.table("M_PROVINSI")), "NAMA_PROVINSI").Order("id").Find(&provinces).Error; err != nil {
		return Data{}, err
	}
	if err := nameFilter(db.Table(r.table("M_KECAMATAN")), "NAMA_KECAMATAN").Order("id").Find(&districts).Error; err != nil {
		return Data{}, err
	}
	if err := db.Table(r.table("M_METODE_PENJUALAN")).Order("id").Find(&salesMethods).Error; err != nil {
		return Data{}, err
	}
	if err := db.Table(r.table("M_KPKNL")).Order("id").Find(&kpknls).Error; err != nil {
		return Data{}, err
	}

	result := Data{
		Categories:   make([]Category, len(categories)),
		AssetTypes:   make([]AssetType, len(assetTypes)),
		Provinces:    make([]Province, len(provinces)),
		Districts:    make([]District, len(districts)),
		SalesMethods: make([]SalesMethod, len(salesMethods)),
		KPKNLs:       make([]KPKNL, len(kpknls)),
	}

	for i, item := range categories {
		result.Categories[i] = Category{ID: item.ID, Name: item.Name}
	}
	for i, item := range assetTypes {
		result.AssetTypes[i] = AssetType{ID: item.ID, CategoryID: item.CategoryID, Name: item.Name}
	}
	for i, item := range provinces {
		result.Provinces[i] = Province{ID: item.ID, Name: item.Name}
	}
	for i, item := range districts {
		result.Districts[i] = District{ID: item.ID, CityID: item.CityID, Name: item.Name}
	}
	for i, item := range salesMethods {
		result.SalesMethods[i] = SalesMethod{ID: item.ID, Name: item.Name}
	}
	for i, item := range kpknls {
		result.KPKNLs[i] = KPKNL{ID: item.ID, Code: item.Code, Name: item.Name}
	}

	return result, nil
}

func (r *ReferenceRepository) GetCitiesByProvinceID(ctx context.Context, provinceID, value string) ([]City, error) {
	var cities []cityModel
	query := r.db.WithContext(ctx).Table(r.table("M_KOTA")).Where("PROVINSI_ID = ?", provinceID)
	if value = strings.TrimSpace(value); value != "" {
		query = query.Where("UPPER(NAMA_KOTA) LIKE ?", "%"+strings.ToUpper(value)+"%")
	}
	if err := query.
		Order("ID").
		Find(&cities).Error; err != nil {
		return nil, err
	}

	result := make([]City, len(cities))
	for i, item := range cities {
		result[i] = City{ID: item.ID, ProvinceID: item.ProvinceID, Name: item.Name}
	}
	return result, nil
}

func (r *ReferenceRepository) GetDistrictsByCityID(ctx context.Context, cityID, value string) ([]District, error) {
	var districts []districtModel
	query := r.db.WithContext(ctx).Table(r.table("M_KECAMATAN")).Where("KOTA_ID = ?", cityID)
	if value = strings.TrimSpace(value); value != "" {
		query = query.Where("UPPER(NAMA_KECAMATAN) LIKE ?", "%"+strings.ToUpper(value)+"%")
	}
	if err := query.Order("ID").Find(&districts).Error; err != nil {
		return nil, err
	}
	result := make([]District, len(districts))
	for i, item := range districts {
		result[i] = District{ID: item.ID, CityID: item.CityID, Name: item.Name}
	}
	return result, nil
}

func (r *ReferenceRepository) GetProvinces(ctx context.Context, value string) ([]Province, error) {
	var provinces []provinceModel
	query := r.db.WithContext(ctx).Table(r.table("M_PROVINSI"))
	if value = strings.TrimSpace(value); value != "" { query = query.Where("UPPER(NAMA_PROVINSI) LIKE ?", "%"+strings.ToUpper(value)+"%") }
	if err := query.Order("ID").Find(&provinces).Error; err != nil { return nil, err }
	result := make([]Province, len(provinces))
	for i, item := range provinces { result[i] = Province{ID: item.ID, Name: item.Name} }
	return result, nil
}
func (r *ReferenceRepository) GetCategories(ctx context.Context, value string) ([]Category, error) { var rows []categoryModel; q:=r.db.WithContext(ctx).Table(r.table("M_KATEGORI")); if value=strings.TrimSpace(value);value!="" { q=q.Where("UPPER(NAMA_KATEGORI) LIKE ?", "%"+strings.ToUpper(value)+"%") }; if err:=q.Order("ID").Find(&rows).Error;err!=nil{return nil,err}; result:=make([]Category,len(rows));for i,x:=range rows{result[i]=Category{ID:x.ID,Name:x.Name}};return result,nil }
func (r *ReferenceRepository) GetAssetTypes(ctx context.Context, value string) ([]AssetType, error) { var rows []assetTypeModel; q:=r.db.WithContext(ctx).Table(r.table("M_TIPE_ASET")); if value=strings.TrimSpace(value);value!="" { q=q.Where("UPPER(NAMA_TIPE) LIKE ?", "%"+strings.ToUpper(value)+"%") }; if err:=q.Order("ID").Find(&rows).Error;err!=nil{return nil,err}; result:=make([]AssetType,len(rows));for i,x:=range rows{result[i]=AssetType{ID:x.ID,CategoryID:x.CategoryID,Name:x.Name}};return result,nil }
