package database

import (
	"time"

	"gorm.io/gorm"
)

// migrationExample is only an example table for demonstrating the migration flow.
type migrationExample struct {
	ID        string    `gorm:"column:ID;primaryKey;size:36"`
	Name      string    `gorm:"column:NAME;not null;size:255"`
	CreatedAt time.Time `gorm:"column:CREATED_AT;not null"`
}

func (migrationExample) TableName() string {
	return "GORM_MIGRATION_EXAMPLE"
}

// AllMigrations is the registry. Add each new migration here once.
func AllMigrations() []Migration {
	return []Migration{
		{
			Version:     "001",
			Description: "create example table",
			Checksum:    "v001-create-gorm-migration-example-v1",
			Up: func(db *gorm.DB, schema string) error {
				table := qualifiedTable(schema, "GORM_MIGRATION_EXAMPLE")
				if schema == "" && db.Migrator().HasTable(&migrationExample{}) {
					return nil
				}
				if schema != "" {
					exists, err := oracleTableExists(db, schema, "GORM_MIGRATION_EXAMPLE")
					if err != nil || exists {
						return err
					}
				}
				return db.Table(table).Migrator().CreateTable(&migrationExample{})
			},
		},
	}
}
