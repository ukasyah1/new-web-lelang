package settings

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"new-website-lelang/internal/platform/dbutil"
	"strings"
)

type SettingsRepository struct {
	db     *gorm.DB
	schema string
}

func NewSettingsRepository(db *gorm.DB, schema string) *SettingsRepository {
	return &SettingsRepository{db: db, schema: strings.ToUpper(strings.TrimSpace(schema))}
}
func (r *SettingsRepository) Get(ctx context.Context, lang string) (map[string]string, []Link, error) {
	settingsMap := map[string]string{}
	var rows []struct{ Key, Value string }
	column := "VALUE_ID"
	title := "TITLE_ID"
	if lang == "en" {
		column = "VALUE_EN"
		title = "TITLE_EN"
	}
	if err := r.db.WithContext(ctx).Raw("SELECT SETTING_KEY AS KEY, " + column + " AS VALUE FROM " + dbutil.QualifiedTable(r.schema, "GLOBAL_SETTINGS")).Scan(&rows).Error; err != nil {
		return nil, nil, fmt.Errorf("query global settings: %w", err)
	}
	for _, x := range rows {
		settingsMap[x.Key] = x.Value
	}
	var links []Link
	if err := r.db.WithContext(ctx).Raw("SELECT ID,NVL(PARENT_ID, '') AS PARENT_ID," + title + " AS TITLE,TARGET_URL AS URL,ORDER_INDEX AS \"ORDER\" FROM " + dbutil.QualifiedTable(r.schema, "NAVBAR_LINKS") + " WHERE IS_ACTIVE=1 ORDER BY PARENT_ID NULLS FIRST, ORDER_INDEX").Scan(&links).Error; err != nil {
		return nil, nil, fmt.Errorf("query navbar links: %w", err)
	}
	return settingsMap, links, nil
}
