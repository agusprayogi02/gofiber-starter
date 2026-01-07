package main

import (
	"fmt"
	"io"
	"os"

	"starter-gofiber/internal/config"

	atlas "ariga.io/atlas-provider-gorm/gormschema"
)

// Atlas GORM Loader
// This program is called by Atlas to load GORM schema from GORM models
// It uses the same model definitions from internal/config for consistency
// This ensures no duplication between AutoMigrate and Atlas migrations
func main() {
	// Load config to initialize ENV variables
	config.LoadConfig()

	// Get models from single source of truth
	models := config.GetModelsForMigration()

	// Determine dialect based on DB_TYPE
	dialect := "postgres" // default
	if config.ENV.DB_TYPE == "mysql" {
		dialect = "mysql"
	} else if config.ENV.DB_TYPE == "sqlserver" {
		dialect = "sqlserver"
	}

	// Load GORM schema and output Atlas schema
	stmts, err := atlas.New(dialect).Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	// Write schema to stdout for Atlas
	io.WriteString(os.Stdout, stmts)
}
