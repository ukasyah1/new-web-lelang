package database

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed migration/*.sql
var sqlMigrationFiles embed.FS

var sqlMigrationNamePattern = regexp.MustCompile(`^V([0-9]+)__(.+)\.sql$`)

// AllMigrations loads embedded SQL files named V{version}__{description}.sql.
func AllMigrations() ([]Migration, error) {
	entries, err := sqlMigrationFiles.ReadDir("migration")
	if err != nil {
		return nil, fmt.Errorf("read SQL migrations: %w", err)
	}

	migrations := make([]Migration, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		matches := sqlMigrationNamePattern.FindStringSubmatch(entry.Name())
		if len(matches) != 3 {
			return nil, fmt.Errorf("nama migration %q harus mengikuti V001__description.sql", entry.Name())
		}

		path := filepath.ToSlash(filepath.Join("migration", entry.Name()))
		content, err := sqlMigrationFiles.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}
		checksum := sha256.Sum256(content)
		migrations = append(migrations, Migration{
			Version:     matches[1],
			Description: strings.ReplaceAll(matches[2], "_", " "),
			Checksum:    hex.EncodeToString(checksum[:]),
			SQL:         string(content),
		})
	}

	return migrations, nil
}
