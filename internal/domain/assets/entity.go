package assets

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("asset not found")

type NamedReference struct {
	ID   string
	Name string
}

type KPKNL struct {
	ID   string
	Code string
	Name string
}

type Asset struct {
	ID, Code, Name, Certificate, Status, Timezone string
	Address, Coordinates, Description             string
	Category, AssetType, Province, City           NamedReference
	KPKNL                                         KPKNL
	AuctionPrice, OriginalPrice                   float64
	LandArea, BuildingArea, ViewCount             int64
	StartDate, EndDate                            *time.Time
	UpdatedAt                                     *time.Time
	ImageURLs, Facilities                         []string
}

type SearchQuery struct {
	Search, CategoryID, AssetTypeID, ProvinceID, CityID, District string
	MinimumPrice, MaximumPrice                                    *int64
	Page, Limit                                                   int
}

type SearchResult struct {
	Assets []Asset
	Total  int64
}
