package database

import (
	"errors"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
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
	migrations := []Migration{{
		Version:     "001",
		Description: "test migration",
		Checksum:    "test-v1",
		SQL:         "CREATE TABLE TEST_MIGRATION_ONCE (ID INTEGER PRIMARY KEY);",
	}}

	if err := RunMigrations(db, "", migrations); err != nil {
		t.Fatalf("first migration run: %v", err)
	}
	if err := RunMigrations(db, "", migrations); err != nil {
		t.Fatalf("second migration run: %v", err)
	}
	if !db.Migrator().HasTable("TEST_MIGRATION_ONCE") {
		t.Fatal("expected migration table to exist")
	}
}

func TestRunMigrationsRejectsChangedChecksum(t *testing.T) {
	db := openMigrationTestDB(t)
	migration := Migration{
		Version:     "001",
		Description: "test migration",
		Checksum:    "original",
		SQL:         "CREATE TABLE TEST_MIGRATION_CHECKSUM (ID INTEGER PRIMARY KEY);",
	}
	if err := RunMigrations(db, "", []Migration{migration}); err != nil {
		t.Fatalf("first migration run: %v", err)
	}

	migration.Checksum = "changed"
	if err := RunMigrations(db, "", []Migration{migration}); err == nil {
		t.Fatal("expected checksum error")
	}
}

func TestRunMigrationsSkipsAppliedVersionAndContinuesNext(t *testing.T) {
	db := openMigrationTestDB(t)
	if err := db.AutoMigrate(&schemaMigration{}); err != nil {
		t.Fatalf("prepare migration history: %v", err)
	}
	if err := db.Exec("CREATE TABLE TEST_ALREADY_APPLIED (ID INTEGER PRIMARY KEY)").Error; err != nil {
		t.Fatalf("create applied table: %v", err)
	}
	applied := schemaMigration{
		Version:     "002",
		Description: "already applied",
		Checksum:    "checksum-002",
		AppliedAt:   time.Now().UTC(),
	}
	if err := db.Create(&applied).Error; err != nil {
		t.Fatalf("record applied migration: %v", err)
	}

	migrations := []Migration{
		{
			Version:     "002",
			Description: "already applied",
			Checksum:    "checksum-002",
			SQL:         "CREATE TABLE TEST_ALREADY_APPLIED (ID INTEGER PRIMARY KEY);",
		},
		{
			Version:     "003",
			Description: "next migration",
			Checksum:    "checksum-003",
			SQL:         "CREATE TABLE TEST_NEXT_MIGRATION (ID INTEGER PRIMARY KEY);",
		},
	}

	if err := RunMigrations(db, "", migrations); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	if !db.Migrator().HasTable("TEST_NEXT_MIGRATION") {
		t.Fatal("expected next migration to run")
	}
}

func TestLoadExampleSQLMigration(t *testing.T) {
	migrations, err := AllMigrations()
	if err != nil {
		t.Fatalf("load SQL migrations: %v", err)
	}
	if len(migrations) != 9 ||
		migrations[0].Version != "001" ||
		migrations[1].Version != "002" ||
		migrations[2].Version != "003" ||
		migrations[3].Version != "004" ||
		migrations[4].Version != "005" ||
		migrations[5].Version != "006" ||
		migrations[6].Version != "007" ||
		migrations[7].Version != "008" ||
		migrations[8].Version != "009" {
		t.Fatalf("unexpected migrations: %+v", migrations)
	}
	for _, migration := range migrations {
		if migration.SQL == "" || migration.Checksum == "" {
			t.Fatalf("expected SQL and checksum to be loaded for migration %s", migration.Version)
		}
	}
}

func TestMigrationExecutesMultipleStatements(t *testing.T) {
	db := openMigrationTestDB(t)
	migration := Migration{
		Version:     "001",
		Description: "multiple statements",
		Checksum:    "multi-v1",
		SQL: "CREATE TABLE TEST_FIRST (ID INTEGER PRIMARY KEY);" +
			"CREATE TABLE TEST_SECOND (ID INTEGER PRIMARY KEY);",
	}
	if err := RunMigrations(db, "", []Migration{migration}); err != nil {
		t.Fatalf("run migration: %v", err)
	}
	if !db.Migrator().HasTable("TEST_FIRST") || !db.Migrator().HasTable("TEST_SECOND") {
		t.Fatal("expected both migration tables to exist")
	}
}

func TestCreateTableAlreadyExistsOracleErrorIsIgnorable(t *testing.T) {
	err := errors.New("ORA-00955: name is already used by an existing object")
	statement := "CREATE TABLE CMS.M_FAQ_CATEGORY (ID VARCHAR2(36) PRIMARY KEY)"

	if !isIgnorableMigrationError(statement, err) {
		t.Fatal("expected ORA-00955 from CREATE TABLE to be ignorable")
	}
}

func TestNonCreateTableOracleErrorIsNotIgnorable(t *testing.T) {
	err := errors.New("ORA-00955: name is already used by an existing object")
	statement := "ALTER TABLE CMS.M_FAQ_CATEGORY ADD NAME VARCHAR2(100)"

	if isIgnorableMigrationError(statement, err) {
		t.Fatal("expected ORA-00955 from non-CREATE TABLE statement to fail")
	}
}
