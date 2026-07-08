package award

import "time"

// Award represents one row from CMS.MST_AWARDS.
type Award struct {
	ID        string
	ImageSrc  *string
	InputDate *time.Time
	Sequence  *int64
	Status    *string
	CreatedBy *string
	CreatedAt *time.Time
	UpdatedBy *string
	UpdatedAt *time.Time
	IsDeleted int
	DeletedBy *string
	DeletedAt *time.Time
	FileName  *string
}
