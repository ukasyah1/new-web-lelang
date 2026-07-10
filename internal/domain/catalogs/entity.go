package catalogs

import "time"

type Catalog struct {
	ID, FileName, FileURL string
	CreatedAt, UpdatedAt  time.Time
}

type MonthlyCatalog struct {
	FileName    string
	FileURL     string
	Size        int64
	PublishedAt time.Time
}
