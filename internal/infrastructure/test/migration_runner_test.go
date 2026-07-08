package infrastructure_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"new-website-lelang/internal/infrastructure/database"
)

func openMigrationTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	return db
}

func TestRunMigrationsAppliesMigrationOnlyOnce(t *testing.T) {
	db := openMigrationTestDB(t)
	migrations := []database.Migration{{
		Version:     "001",
		Description: "test migration",
		Checksum:    "test-v1",
		SQL:         "CREATE TABLE TEST_MIGRATION_ONCE (ID INTEGER PRIMARY KEY);",
	}}

	if err := database.RunMigrations(db, "", migrations); err != nil {
		t.Fatalf("first migration run: %v", err)
	}
	if err := database.RunMigrations(db, "", migrations); err != nil {
		t.Fatalf("second migration run: %v", err)
	}
	if !db.Migrator().HasTable("TEST_MIGRATION_ONCE") {
		t.Fatal("expected migration table to exist")
	}
}

func TestRunMigrationsRejectsChangedChecksum(t *testing.T) {
	db := openMigrationTestDB(t)
	migration := database.Migration{
		Version:     "001",
		Description: "test migration",
		Checksum:    "original",
		SQL:         "CREATE TABLE TEST_MIGRATION_CHECKSUM (ID INTEGER PRIMARY KEY);",
	}
	if err := database.RunMigrations(db, "", []database.Migration{migration}); err != nil {
		t.Fatalf("first migration run: %v", err)
	}

	migration.Checksum = "changed"
	if err := database.RunMigrations(db, "", []database.Migration{migration}); err == nil {
		t.Fatal("expected checksum error")
	}
}

func TestLoadExampleSQLMigration(t *testing.T) {
	migrations, err := database.AllMigrations()
	if err != nil {
		t.Fatalf("load SQL migrations: %v", err)
	}
	if len(migrations) != 1 || migrations[0].Version != "001" {
		t.Fatalf("unexpected migrations: %+v", migrations)
	}
	if migrations[0].SQL == "" || migrations[0].Checksum == "" {
		t.Fatal("expected SQL and checksum to be loaded")
	}
}

func TestMigrationExecutesMultipleStatements(t *testing.T) {
	db := openMigrationTestDB(t)
	migration := database.Migration{
		Version:     "001",
		Description: "multiple statements",
		Checksum:    "multi-v1",
		SQL: "CREATE TABLE TEST_FIRST (ID INTEGER PRIMARY KEY);" +
			"CREATE TABLE TEST_SECOND (ID INTEGER PRIMARY KEY);",
	}
	if err := database.RunMigrations(db, "", []database.Migration{migration}); err != nil {
		t.Fatalf("run migration: %v", err)
	}
	if !db.Migrator().HasTable("TEST_FIRST") || !db.Migrator().HasTable("TEST_SECOND") {
		t.Fatal("expected both migration tables to exist")
	}
}
