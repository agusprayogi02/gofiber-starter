package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	Version     uint
	Description string
	Up          func(*gorm.DB) error
	Down        func(*gorm.DB) error
}

var migrations = make(map[uint]*Migration)

// RegisterMigration registers a new migration
func RegisterMigration(version uint, description string, up, down func(*gorm.DB) error) {
	migrations[version] = &Migration{
		Version:     version,
		Description: description,
		Up:          up,
		Down:        down,
	}
}

// RunMigrations runs all pending migrations
func RunMigrations(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Detect database type from GORM
	dbType := db.Dialector.Name()

	// Create appropriate driver based on database type
	var driver *migrate.Migrate
	var driverErr error

	if dbType == "postgres" {
		pgDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
		if err != nil {
			return fmt.Errorf("failed to create postgres migration driver: %w", err)
		}
		migrationsPath := getMigrationsPath()
		driver, driverErr = migrate.NewWithDatabaseInstance(
			"file://"+migrationsPath,
			"postgres",
			pgDriver,
		)
	} else {
		// Default to MySQL
		mysqlDriver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
		if err != nil {
			return fmt.Errorf("failed to create mysql migration driver: %w", err)
		}
		migrationsPath := getMigrationsPath()
		driver, driverErr = migrate.NewWithDatabaseInstance(
			"file://"+migrationsPath,
			"mysql",
			mysqlDriver,
		)
	}

	if driverErr != nil {
		return fmt.Errorf("failed to create migration instance: %w", driverErr)
	}

	// Run migrations
	if err := driver.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	Info("Migrations completed successfully")
	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(db *gorm.DB, steps int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Detect database type from GORM
	dbType := db.Dialector.Name()

	var driver *migrate.Migrate
	migrationsPath := getMigrationsPath()

	if dbType == "postgres" {
		pgDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
		if err != nil {
			return err
		}
		driver, err = migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", migrationsPath),
			"postgres",
			pgDriver,
		)
		if err != nil {
			return err
		}
	} else {
		mysqlDriver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
		if err != nil {
			return err
		}
		driver, err = migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", migrationsPath),
			"mysql",
			mysqlDriver,
		)
		if err != nil {
			return err
		}
	}

	// Rollback specified steps
	return driver.Steps(-steps)
}

// CreateMigration creates a new migration file
func CreateMigration(name string) error {
	timestamp := time.Now().Unix()
	migrationsPath := getMigrationsPath()

	// Ensure migrations directory exists
	if err := os.MkdirAll(migrationsPath, 0o755); err != nil {
		return err
	}

	// Create up migration file
	upFile := filepath.Join(migrationsPath, fmt.Sprintf("%d_%s.up.sql", timestamp, name))
	downFile := filepath.Join(migrationsPath, fmt.Sprintf("%d_%s.down.sql", timestamp, name))

	// Write template content
	upContent := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n-- Write your UP migration here\n", name, time.Now().Format(time.RFC3339))
	downContent := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n-- Write your DOWN migration here\n", name, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(upFile, []byte(upContent), 0o644); err != nil {
		return err
	}

	if err := os.WriteFile(downFile, []byte(downContent), 0o644); err != nil {
		return err
	}

	Info(fmt.Sprintf("Created migration files:\n  - %s\n  - %s", upFile, downFile))
	return nil
}

// GetMigrationVersion returns current migration version
func GetMigrationVersion(db *gorm.DB) (uint, bool, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return 0, false, err
	}

	// Detect database type from GORM
	dbType := db.Dialector.Name()

	var driver *migrate.Migrate
	migrationsPath := getMigrationsPath()

	if dbType == "postgres" {
		pgDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
		if err != nil {
			return 0, false, err
		}
		driver, err = migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", migrationsPath),
			"postgres",
			pgDriver,
		)
		if err != nil {
			return 0, false, err
		}
	} else {
		mysqlDriver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
		if err != nil {
			return 0, false, err
		}
		driver, err = migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", migrationsPath),
			"mysql",
			mysqlDriver,
		)
		if err != nil {
			return 0, false, err
		}
	}

	version, dirty, err := driver.Version()
	return version, dirty, err
}

// ForceMigrationVersion forces migration to a specific version
func ForceMigrationVersion(db *gorm.DB, version int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Detect database type from GORM
	dbType := db.Dialector.Name()

	var driver *migrate.Migrate
	migrationsPath := getMigrationsPath()

	if dbType == "postgres" {
		pgDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
		if err != nil {
			return err
		}
		driver, err = migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", migrationsPath),
			"postgres",
			pgDriver,
		)
		if err != nil {
			return err
		}
	} else {
		mysqlDriver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
		if err != nil {
			return err
		}
		driver, err = migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", migrationsPath),
			"mysql",
			mysqlDriver,
		)
		if err != nil {
			return err
		}
	}

	return driver.Force(version)
}

// getMigrationsPath returns the path to migrations directory
func getMigrationsPath() string {
	// Try to find migrations directory
	paths := []string{
		"./migrations",
		"../migrations",
		"../../migrations",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath
		}
	}

	// Default to ./migrations
	absPath, _ := filepath.Abs("./migrations")
	return absPath
}
