package database

import (
	"strings"
	"testing"
)

func TestOracleDSNFromJDBCURL(t *testing.T) {
	dsn, err := BuildOracleDSN(
		"jdbc:oracle:thin:@//localhost:1521/FREEPDB1",
		"test-user",
		"test-password",
	)
	if err != nil {
		t.Fatalf("build Oracle DSN: %v", err)
	}

	for _, expected := range []string{"localhost", "1521", "FREEPDB1", "test-user"} {
		if !strings.Contains(dsn, expected) {
			t.Fatalf("expected DSN to contain %q", expected)
		}
	}
}

func TestOracleDSNRejectsInvalidURL(t *testing.T) {
	if _, err := BuildOracleDSN("localhost", "user", "password"); err == nil {
		t.Fatal("expected invalid URL error")
	}
}
