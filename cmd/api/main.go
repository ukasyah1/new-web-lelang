package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"new-website-lelang/internal/domain/award"
	"new-website-lelang/internal/domain/banner"
	"new-website-lelang/internal/domain/catalog"
	"new-website-lelang/internal/domain/faq"
	"new-website-lelang/internal/domain/reference"
	"new-website-lelang/internal/infrastructure/database"
	"new-website-lelang/internal/interfaces/httpapi"
)

func main() {
	loadEnvironment()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run describes the application startup from top to bottom.
func run() error {
	config := loadConfig()

	referenceHandler, err := buildReferenceHandler(config.sqlitePath)
	if err != nil {
		return err
	}
	db, err := openApplicationOracleDB(config)
	if err != nil {
		return err
	}
	awardHandler := buildAwardHandler(db)
	faqHandler := buildFAQHandler(db)
	bannerHandler := buildBannerHandler(db)
	assetHandler := catalog.NewAssetHandler()

	router := httpapi.NewRouter(referenceHandler, assetHandler, awardHandler, faqHandler, bannerHandler)
	server := newHTTPServer(config.port, router)

	return startAndWaitForShutdown(server, config.port)
}

type appConfig struct {
	port             string
	sqlitePath       string
	databaseURL      string
	databaseUsername string
	databasePassword string
	runMigrations    bool
	migrationSchema  string
}

func loadConfig() appConfig {
	return appConfig{
		port:             getEnv("PORT", "80"),
		sqlitePath:       getEnv("SQLITE_PATH", "lelang.db"),
		databaseURL:      getEnv("DATABASE_URL", os.Getenv("DATABASE_PATH")),
		databaseUsername: os.Getenv("DATABASE_USERNAME"),
		databasePassword: os.Getenv("DATABASE_PASSWORD"),
		runMigrations:    getEnv("RUN_MIGRATIONS", "false") == "true",
		migrationSchema:  getEnv("MIGRATION_SCHEMA", "CMS"),
	}
}

// loadEnvironment uses .env when available. The example file is a local fallback.
func loadEnvironment() {
	if err := godotenv.Load(); err != nil {
		_ = godotenv.Load(".env.example")
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// buildReferenceHandler connects database -> repository -> service -> HTTP handler.
func buildReferenceHandler(databasePath string) (*reference.ReferenceHandler, error) {
	db, err := database.OpenSQLite(databasePath)
	if err != nil {
		return nil, err
	}

	repository := database.NewReferenceRepository(db)
	if err := repository.Prepare(); err != nil {
		return nil, fmt.Errorf("prepare reference repository: %w", err)
	}

	service := reference.NewService(repository)
	return reference.NewReferenceHandler(service), nil
}

func openApplicationOracleDB(config appConfig) (*gorm.DB, error) {
	db, err := database.OpenOracle(
		config.databaseURL,
		config.databaseUsername,
		config.databasePassword,
	)
	if err != nil {
		return nil, err
	}
	if config.runMigrations {
		available, err := database.AllMigrations()
		if err != nil {
			return nil, fmt.Errorf("load database migrations: %w", err)
		}
		if err := database.RunMigrations(db, config.migrationSchema, available); err != nil {
			return nil, fmt.Errorf("run database migrations: %w", err)
		}
	}

	return db, nil
}

func buildAwardHandler(db *gorm.DB) *award.AwardHandler {
	repository := database.NewAwardRepository(db)
	service := award.NewService(repository)
	return award.NewAwardHandler(service)
}

func buildFAQHandler(db *gorm.DB) *faq.FAQHandler {
	repository := database.NewFAQRepository(db)
	service := faq.NewService(repository)
	return faq.NewFAQHandler(service)
}

func buildBannerHandler(db *gorm.DB) *banner.BannerHandler {
	repository := database.NewBannerRepository(db)
	service := banner.NewService(repository)
	return banner.NewBannerHandler(service)
}

func newHTTPServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func startAndWaitForShutdown(server *http.Server, port string) error {
	serverError := make(chan error, 1)
	go func() {
		log.Printf("API listening on http://localhost:%s", port)
		serverError <- server.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverError:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	case <-ctx.Done():
		log.Println("Shutting down API...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	return nil
}
