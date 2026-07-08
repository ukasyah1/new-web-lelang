package database

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestRunMigrationsAppliesMigrationOnlyOnce(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	applied := 0
	migrations := []Migration{
		{
			Version:     "001",
			Description: "test migration",
			Checksum:    "test-v1",
			Up: func(_ *gorm.DB, _ string) error {
				applied++
				return nil
			},
		},
	}

	if err := RunMigrations(db, "", migrations); err != nil {
		t.Fatalf("first migration run: %v", err)
	}
	if err := RunMigrations(db, "", migrations); err != nil {
		t.Fatalf("second migration run: %v", err)
	}
	if applied != 1 {
		t.Fatalf("expected migration once, got %d", applied)
	}
}

func TestRunMigrationsRejectsChangedChecksum(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	migration := Migration{
		Version:     "001",
		Description: "test migration",
		Checksum:    "original",
		Up:          func(_ *gorm.DB, _ string) error { return nil },
	}
	if err := RunMigrations(db, "", []Migration{migration}); err != nil {
		t.Fatalf("first migration run: %v", err)
	}

	migration.Checksum = "changed"
	if err := RunMigrations(db, "", []Migration{migration}); err == nil {
		t.Fatal("expected checksum error")
	}
}

func TestExampleMigrationCreatesTable(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	if err := RunMigrations(db, "", AllMigrations()); err != nil {
		t.Fatalf("run example migration: %v", err)
	}
	if !db.Migrator().HasTable(&migrationExample{}) {
		t.Fatal("expected example table to exist")
	}
}
