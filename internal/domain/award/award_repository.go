package award

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type awardModel struct {
	ID        string     `gorm:"column:ID"`
	ImageSrc  *string    `gorm:"column:IMG_SRC"`
	InputDate *time.Time `gorm:"column:INPUT_DATE"`
	Sequence  *int64     `gorm:"column:SEQ"`
	Status    *string    `gorm:"column:STATUS"`
	CreatedBy *string    `gorm:"column:CREATED_BY"`
	CreatedAt *time.Time `gorm:"column:CREATED_AT"`
	UpdatedBy *string    `gorm:"column:UPDATED_BY"`
	UpdatedAt *time.Time `gorm:"column:UPDATED_AT"`
	IsDeleted int        `gorm:"column:IS_DELETED"`
	DeletedBy *string    `gorm:"column:DELETED_BY"`
	DeletedAt *time.Time `gorm:"column:DELETED_AT"`
	FileName  *string    `gorm:"column:FILE_NAME"`
}

type AwardRepository struct {
	db *gorm.DB
}

func NewAwardRepository(db *gorm.DB) *AwardRepository {
	return &AwardRepository{db: db}
}

func (r *AwardRepository) GetAll(ctx context.Context) ([]Award, error) {
	var models []awardModel
	result := r.db.WithContext(ctx).Raw(`
		SELECT ID, IMG_SRC, INPUT_DATE, SEQ, STATUS,
		       CREATED_BY, CREATED_AT, UPDATED_BY, UPDATED_AT,
		       IS_DELETED, DELETED_BY, DELETED_AT, FILE_NAME
		FROM CMS.MST_AWARDS
		WHERE IS_DELETED = ?
		ORDER BY SEQ ASC`, 0).Scan(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("query CMS.MST_AWARDS: %w", result.Error)
	}

	records := make([]Award, len(models))
	for i, model := range models {
		records[i] = Award{
			ID:        model.ID,
			ImageSrc:  model.ImageSrc,
			InputDate: model.InputDate,
			Sequence:  model.Sequence,
			Status:    model.Status,
			CreatedBy: model.CreatedBy,
			CreatedAt: model.CreatedAt,
			UpdatedBy: model.UpdatedBy,
			UpdatedAt: model.UpdatedAt,
			IsDeleted: model.IsDeleted,
			DeletedBy: model.DeletedBy,
			DeletedAt: model.DeletedAt,
			FileName:  model.FileName,
		}
	}

	return records, nil
}
