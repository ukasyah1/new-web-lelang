package masterdata

// Category represents an auction asset category.
type Category struct {
	ID   string
	Name string
}

// AssetType represents an asset type belonging to a category.
type AssetType struct {
	ID         string
	CategoryID string
	Name       string
}

// Province represents a province available for auction searches.
type Province struct {
	ID   string
	Name string
}

// City represents a city belonging to a province.
type City struct {
	ID         string
	ProvinceID string
	Name       string
}

// District represents a district belonging to a city.
type District struct {
	ID     string
	CityID string
	Name   string
}

// SalesMethod represents a supported asset sales method.
type SalesMethod struct {
	ID   string
	Name string
}

// KPKNL represents a KPKNL office.
type KPKNL struct {
	ID   string
	Code string
	Name string
}

// Data groups reference data needed by the auction client.
type Data struct {
	Categories   []Category
	AssetTypes   []AssetType
	Provinces    []Province
	Districts    []District
	SalesMethods []SalesMethod
	KPKNLs       []KPKNL
}
