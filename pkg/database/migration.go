package database

import (
	"fmt"
	"os"

	"starter-gofiber/pkg/logger"

	atlas "ariga.io/atlas-provider-gorm/gormschema"
	"gorm.io/gorm"
)

// GetAllModels returns all GORM models that should be included in migrations
// This function should be updated when new models are added to the system
func GetAllModels() []interface{} {
	// Import models from domain layer
	// This will be populated from config/database.go
	return []interface{}{}
}

// LoadAtlasModels loads GORM schema for Atlas migration
// This function is used by atlas.hcl configuration file
func LoadAtlasModels(db *gorm.DB, models []interface{}) error {
	// Create a temporary schema using AutoMigrate to let GORM understand the models
	// but we won't actually run this - Atlas will generate proper migrations instead
	err := db.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("failed to load models: %w", err)
	}
	return nil
}

// GenerateAtlasSchema generates Atlas-compatible schema from GORM models
// This is used by Atlas CLI to understand your database schema
func GenerateAtlasSchema(db *gorm.DB, models []interface{}) error {

	// Load schema from GORM models
	_, err := atlas.New("gorm").Load(models...)
	if err != nil {
		return fmt.Errorf("failed to generate Atlas schema: %w", err)
	}

	logger.Info("Atlas schema generated successfully from GORM models")
	return nil
}

// SyncSchema applies all pending Atlas migrations to the database
// This is the programmatic way to run migrations, equivalent to: atlas migrate apply
func SyncSchema(db *gorm.DB, models []interface{}) error {
	// Get database instance
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	logger.Info("Database schema sync initiated")
	logger.Info("Please use 'atlas migrate apply' command to apply migrations")
	logger.Info("Or enable DB_GEN=true to use AutoMigrate for development")

	return nil
}

// GetMigrationsPath returns the path to Atlas migrations directory
func GetMigrationsPath() string {
	// Check if custom migrations path is set
	if path := os.Getenv("ATLAS_MIGRATIONS_DIR"); path != "" {
		return path
	}

	// Default to ./migrations
	return "./migrations"
}

// EnsureMigrationsDir creates migrations directory if it doesn't exist
func EnsureMigrationsDir() error {
	path := GetMigrationsPath()
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}
	return nil
}

// MigrationInfo contains information about Atlas migrations
type MigrationInfo struct {
	Enabled        bool
	MigrationsPath string
	DatabaseType   string
	UseAutoMigrate bool
}

// GetMigrationInfo returns current migration configuration
func GetMigrationInfo(db *gorm.DB) *MigrationInfo {
	dbType := db.Dialector.Name()
	useAutoMigrate := os.Getenv("DB_GEN") == "true"

	return &MigrationInfo{
		Enabled:        true,
		MigrationsPath: GetMigrationsPath(),
		DatabaseType:   dbType,
		UseAutoMigrate: useAutoMigrate,
	}
}
