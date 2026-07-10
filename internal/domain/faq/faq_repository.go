package faq

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type faqRow struct {
	CategoryID   string `gorm:"column:CATEGORY_ID"`
	CategoryName string `gorm:"column:CATEGORY_NAME"`
	FAQID        string `gorm:"column:FAQ_ID"`
	Question     string `gorm:"column:QUESTION"`
	Answer       string `gorm:"column:ANSWER"`
}

type FAQRepository struct {
	db *gorm.DB
}

func NewFAQRepository(db *gorm.DB) *FAQRepository {
	return &FAQRepository{db: db}
}

func (r *FAQRepository) GetAll(ctx context.Context, lang string) ([]Category, error) {
	categoryNameColumn, questionColumn, answerColumn := faqLanguageColumns(lang)
	rows := []faqRow{}
	result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`
		SELECT c.ID AS CATEGORY_ID,
		       c.%s AS CATEGORY_NAME,
		       f.ID AS FAQ_ID,
		       f.%s AS QUESTION,
		       f.%s AS ANSWER
		FROM CMS.M_FAQ_CATEGORY c
		JOIN CMS.FAQS f ON f.CATEGORY_ID = c.ID
		WHERE f.IS_ACTIVE = ?
		ORDER BY c.ORDER_INDEX ASC, f.ORDER_INDEX ASC`, categoryNameColumn, questionColumn, answerColumn), 1).Scan(&rows)
	if result.Error != nil {
		return nil, fmt.Errorf("query CMS.FAQS: %w", result.Error)
	}

	return mapFAQRows(rows), nil
}

func faqLanguageColumns(lang string) (string, string, string) {
	if lang == "en" {
		return "NAMA_KATEGORI_EN", "QUESTION_EN", "ANSWER_EN"
	}
	return "NAMA_KATEGORI_ID", "QUESTION_ID", "ANSWER_ID"
}

func mapFAQRows(rows []faqRow) []Category {
	categories := make([]Category, 0)
	categoryIndexes := make(map[string]int)

	for _, row := range rows {
		index, exists := categoryIndexes[row.CategoryID]
		if !exists {
			index = len(categories)
			categoryIndexes[row.CategoryID] = index
			categories = append(categories, Category{
				ID:   row.CategoryID,
				Name: row.CategoryName,
				FAQs: []FAQ{},
			})
		}

		categories[index].FAQs = append(categories[index].FAQs, FAQ{
			ID:       row.FAQID,
			Question: row.Question,
			Answer:   row.Answer,
		})
	}

	return categories
}
