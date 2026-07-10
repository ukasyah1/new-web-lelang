package assets

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type assetRow struct {
	ID, Code, Name, Status                                string
	CategoryID, CategoryName, AssetTypeID, AssetTypeName  *string
	ProvinceID, ProvinceName, CityID, CityName            *string
	KPKNLID, KPKNLCode, KPKNLName                         *string
	Address, Certificate, Timezone, ImageURLs, Facilities *string
	Coordinates, Description                              *string
	AuctionPrice, OriginalPrice                           *float64
	LandArea, BuildingArea, ViewCount                     *int64
	StartDate, EndDate, UpdatedAt                         *time.Time
}

type AssetRepository struct{ db *gorm.DB }

func NewAssetRepository(db *gorm.DB) *AssetRepository { return &AssetRepository{db: db} }

func (r *AssetRepository) Search(ctx context.Context, query SearchQuery) (SearchResult, error) {
	from := ` FROM CMS.ASSETS a
		LEFT JOIN CMS.M_TIPE_ASET ta ON ta.ID = a.TIPE_ASET_ID
		LEFT JOIN CMS.M_KATEGORI k ON k.ID = ta.KATEGORI_ID
		LEFT JOIN CMS.M_PROVINSI p ON p.ID = a.PROVINSI_ID
		LEFT JOIN CMS.M_KOTA ko ON ko.ID = a.KOTA_ID
		LEFT JOIN CMS.M_KPKNL kp ON kp.ID = a.KPKNL_ID`
	where, args := assetWhere(query)
	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(*)"+from+where, args...).Scan(&total).Error; err != nil {
		return SearchResult{}, fmt.Errorf("count CMS.ASSETS: %w", err)
	}
	selectSQL := assetSelectSQL
	args = append(args, (query.Page-1)*query.Limit, query.Limit)
	var rows []assetRow
	if err := r.db.WithContext(ctx).Raw(selectSQL+from+where+" ORDER BY a.CREATED_AT DESC OFFSET ? ROWS FETCH NEXT ? ROWS ONLY", args...).Scan(&rows).Error; err != nil {
		return SearchResult{}, fmt.Errorf("query CMS.ASSETS: %w", err)
	}
	resultAssets := make([]Asset, len(rows))
	for i, row := range rows {
		resultAssets[i] = mapAssetRow(row)
	}
	return SearchResult{Assets: resultAssets, Total: total}, nil
}

const assetFromSQL = ` FROM CMS.ASSETS a
	LEFT JOIN CMS.M_TIPE_ASET ta ON ta.ID = a.TIPE_ASET_ID
	LEFT JOIN CMS.M_KATEGORI k ON k.ID = ta.KATEGORI_ID
	LEFT JOIN CMS.M_PROVINSI p ON p.ID = a.PROVINSI_ID
	LEFT JOIN CMS.M_KOTA ko ON ko.ID = a.KOTA_ID
	LEFT JOIN CMS.M_KPKNL kp ON kp.ID = a.KPKNL_ID`

const assetSelectSQL = `SELECT a.ID, a.KODE_ASET AS CODE, a.NAMA_ASET AS NAME, a.STATUS,
	k.ID AS CATEGORY_ID, k.NAMA_KATEGORI AS CATEGORY_NAME, ta.ID AS ASSET_TYPE_ID, ta.NAMA_TIPE AS ASSET_TYPE_NAME,
	p.ID AS PROVINCE_ID, p.NAMA_PROVINSI AS PROVINCE_NAME, ko.ID AS CITY_ID, ko.NAMA_KOTA AS CITY_NAME,
	kp.ID AS KPKNL_ID, kp.KODE_KPKNL AS KPKNL_CODE, kp.NAMA_KANTOR AS KPKNL_NAME,
	a.ALAMAT AS ADDRESS, a.JENIS_SERTIFIKAT AS CERTIFICATE, a.ZONA_WAKTU AS TIMEZONE, a.IMAGE_URLS, a.FASILITAS,
	a.LATITUDE || ',' || a.LONGITUDE AS COORDINATES, a.DESKRIPSI AS DESCRIPTION,
	a.LIMIT_LELANG AS AUCTION_PRICE, a.NILAI_PASAR AS ORIGINAL_PRICE, a.LUAS_TANAH AS LAND_AREA,
	a.LUAS_BANGUNAN AS BUILDING_AREA, a.VIEW_COUNT, a.START_DATE, a.END_DATE, a.UPDATED_AT`

func (r *AssetRepository) GetByID(ctx context.Context, id string) (Asset, error) {
	var row assetRow
	result := r.db.WithContext(ctx).Raw(
		assetSelectSQL+assetFromSQL+" WHERE a.ID = ? AND a.STATUS IN ('PUBLISHED', 'LELANG')",
		id,
	).Scan(&row)
	if result.Error != nil {
		return Asset{}, fmt.Errorf("query CMS.ASSETS by ID: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return Asset{}, ErrNotFound
	}
	return mapAssetRow(row), nil
}

func assetWhere(q SearchQuery) (string, []any) {
	parts := []string{"a.STATUS IN ('PUBLISHED', 'LELANG')"}
	args := []any{}
	if q.Search != "" {
		parts = append(parts, "(UPPER(a.NAMA_ASET) LIKE ? OR UPPER(a.KODE_ASET) LIKE ? OR UPPER(a.ALAMAT) LIKE ?)")
		term := "%" + strings.ToUpper(q.Search) + "%"
		args = append(args, term, term, term)
	}
	filters := []struct{ value, clause string }{{q.CategoryID, "k.ID = ?"}, {q.AssetTypeID, "a.TIPE_ASET_ID = ?"}, {q.ProvinceID, "a.PROVINSI_ID = ?"}, {q.CityID, "a.KOTA_ID = ?"}, {q.District, "UPPER(a.KECAMATAN) = ?"}}
	for _, f := range filters {
		if f.value != "" {
			parts = append(parts, f.clause)
			value := f.value
			if strings.Contains(f.clause, "UPPER") {
				value = strings.ToUpper(value)
			}
			args = append(args, value)
		}
	}
	if q.MinimumPrice != nil {
		parts = append(parts, "a.LIMIT_LELANG >= ?")
		args = append(args, *q.MinimumPrice)
	}
	if q.MaximumPrice != nil {
		parts = append(parts, "a.LIMIT_LELANG <= ?")
		args = append(args, *q.MaximumPrice)
	}
	return " WHERE " + strings.Join(parts, " AND "), args
}

func mapAssetRow(r assetRow) Asset {
	a := Asset{ID: r.ID, Code: r.Code, Name: r.Name, Status: r.Status, ImageURLs: jsonArray(r.ImageURLs), Facilities: jsonArray(r.Facilities)}
	if r.CategoryID != nil {
		a.Category.ID = *r.CategoryID
	}
	if r.CategoryName != nil {
		a.Category.Name = *r.CategoryName
	}
	if r.AssetTypeID != nil {
		a.AssetType.ID = *r.AssetTypeID
	}
	if r.AssetTypeName != nil {
		a.AssetType.Name = *r.AssetTypeName
	}
	if r.ProvinceID != nil {
		a.Province.ID = *r.ProvinceID
	}
	if r.ProvinceName != nil {
		a.Province.Name = *r.ProvinceName
	}
	if r.CityID != nil {
		a.City.ID = *r.CityID
	}
	if r.CityName != nil {
		a.City.Name = *r.CityName
	}
	if r.KPKNLID != nil {
		a.KPKNL.ID = *r.KPKNLID
	}
	if r.KPKNLCode != nil {
		a.KPKNL.Code = *r.KPKNLCode
	}
	if r.KPKNLName != nil {
		a.KPKNL.Name = *r.KPKNLName
	}
	if r.Certificate != nil {
		a.Certificate = *r.Certificate
	}
	if r.Address != nil {
		a.Address = *r.Address
	}
	if r.Coordinates != nil {
		a.Coordinates = *r.Coordinates
	}
	if r.Description != nil {
		a.Description = *r.Description
	}
	if r.Timezone != nil {
		a.Timezone = *r.Timezone
	}
	if r.AuctionPrice != nil {
		a.AuctionPrice = *r.AuctionPrice
	}
	if r.OriginalPrice != nil {
		a.OriginalPrice = *r.OriginalPrice
	}
	if r.LandArea != nil {
		a.LandArea = *r.LandArea
	}
	if r.BuildingArea != nil {
		a.BuildingArea = *r.BuildingArea
	}
	if r.ViewCount != nil {
		a.ViewCount = *r.ViewCount
	}
	a.StartDate = r.StartDate
	a.EndDate = r.EndDate
	a.UpdatedAt = r.UpdatedAt
	return a
}

func jsonArray(value *string) []string {
	result := []string{}
	if value != nil {
		_ = json.Unmarshal([]byte(*value), &result)
	}
	return result
}
