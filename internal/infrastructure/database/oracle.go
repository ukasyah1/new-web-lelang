package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	oracle "github.com/godoes/gorm-oracle"
	"gorm.io/gorm"
)

func OpenOracle(jdbcURL, username, password string) (*gorm.DB, error) {
	if strings.TrimSpace(jdbcURL) == "" || username == "" || password == "" {
		return nil, fmt.Errorf("DATABASE_URL, DATABASE_USERNAME, dan DATABASE_PASSWORD wajib diisi")
	}

	dsn, err := oracleDSN(jdbcURL, username, password)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(oracle.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open Oracle database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get Oracle connection pool: %w", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping Oracle database: %w", err)
	}

	return db, nil
}

func oracleDSN(jdbcURL, username, password string) (string, error) {
	address := strings.TrimSpace(jdbcURL)
	address = strings.TrimPrefix(address, "jdbc:oracle:thin:@")
	address = strings.TrimPrefix(address, "//")

	parsed, err := url.Parse("oracle://" + address)
	if err != nil {
		return "", fmt.Errorf("parse Oracle DATABASE_URL: %w", err)
	}

	host := parsed.Hostname()
	portText := parsed.Port()
	service := strings.Trim(parsed.Path, "/")
	if host == "" || portText == "" || service == "" {
		return "", fmt.Errorf("DATABASE_URL harus berbentuk jdbc:oracle:thin:@//host:port/service")
	}
	if _, _, err := net.SplitHostPort(parsed.Host); err != nil {
		return "", fmt.Errorf("invalid Oracle host and port: %w", err)
	}

	port, err := strconv.Atoi(portText)
	if err != nil {
		return "", fmt.Errorf("invalid Oracle port: %w", err)
	}

	return oracle.BuildUrl(host, port, service, username, password, nil), nil
}
