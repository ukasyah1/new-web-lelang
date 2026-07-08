package database

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Migration is one immutable, versioned database change.
type Migration struct {
	Version     string
	Description string
	Checksum    string
	Up          func(db *gorm.DB, schema string) error
}

type schemaMigration struct {
	Version     string    `gorm:"column:VERSION;primaryKey;size:50"`
	Description string    `gorm:"column:DESCRIPTION;not null;size:255"`
	Checksum    string    `gorm:"column:CHECKSUM;not null;size:100"`
	AppliedAt   time.Time `gorm:"column:APPLIED_AT;not null"`
}

func (schemaMigration) TableName() string {
	return "GORM_SCHEMA_MIGRATIONS"
}

// RunMigrations applies pending migrations in version order and records each success.
func RunMigrations(db *gorm.DB, schema string, available []Migration) error {
	schema = strings.ToUpper(strings.TrimSpace(schema))
	if schema != "" && !validOracleIdentifier(schema) {
		return fmt.Errorf("migration schema %q tidak valid", schema)
	}

	historyTable := qualifiedTable(schema, "GORM_SCHEMA_MIGRATIONS")
	if err := prepareMigrationHistory(db, schema, historyTable); err != nil {
		return fmt.Errorf("prepare migration history: %w", err)
	}
	historyDB := db.Table(historyTable)

	migrations := append([]Migration(nil), available...)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	seen := make(map[string]struct{}, len(migrations))
	for _, item := range migrations {
		if err := validateMigration(item, seen); err != nil {
			return err
		}

		var applied schemaMigration
		result := historyDB.Where("VERSION = ?", item.Version).Limit(1).Find(&applied)
		if result.Error != nil {
			return fmt.Errorf("check migration %s: %w", item.Version, result.Error)
		}
		if result.RowsAffected > 0 {
			if applied.Checksum != item.Checksum {
				return fmt.Errorf("migration %s checksum berubah setelah diterapkan", item.Version)
			}
			continue
		}

		if err := item.Up(db, schema); err != nil {
			return fmt.Errorf("apply migration %s (%s): %w", item.Version, item.Description, err)
		}

		history := schemaMigration{
			Version:     item.Version,
			Description: item.Description,
			Checksum:    item.Checksum,
			AppliedAt:   time.Now().UTC(),
		}
		if err := historyDB.Create(&history).Error; err != nil {
			return fmt.Errorf("record migration %s: %w", item.Version, err)
		}
	}

	return nil
}

var oracleIdentifierPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_$#]*$`)

func validOracleIdentifier(value string) bool {
	return oracleIdentifierPattern.MatchString(value)
}

func qualifiedTable(schema, table string) string {
	if schema == "" {
		return table
	}
	return schema + "." + table
}

func prepareMigrationHistory(db *gorm.DB, schema, table string) error {
	if schema == "" {
		return db.AutoMigrate(&schemaMigration{})
	}

	exists, err := oracleTableExists(db, schema, "GORM_SCHEMA_MIGRATIONS")
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return db.Table(table).Migrator().CreateTable(&schemaMigration{})
}

func oracleTableExists(db *gorm.DB, schema, table string) (bool, error) {
	var count int64
	result := db.Raw(
		"SELECT COUNT(*) FROM ALL_TABLES WHERE OWNER = ? AND TABLE_NAME = ?",
		schema,
		table,
	).Scan(&count)
	return count > 0, result.Error
}

func validateMigration(item Migration, seen map[string]struct{}) error {
	if item.Version == "" || item.Description == "" || item.Checksum == "" || item.Up == nil {
		return fmt.Errorf("migration version, description, checksum, dan Up wajib diisi")
	}
	if _, exists := seen[item.Version]; exists {
		return fmt.Errorf("migration version %s duplikat", item.Version)
	}
	seen[item.Version] = struct{}{}
	return nil
}
