package database

import (
	"fmt"
	"strings"

	"starter-gofiber/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// IndexDefinition represents a database index definition
type IndexDefinition struct {
	Name       string   // Index name
	Table      string   // Table name
	Columns    []string // Column names
	Unique     bool     // Whether index is unique
	Concurrent bool     // Whether to create index concurrently (PostgreSQL only)
	Where      string   // Partial index condition (PostgreSQL only)
}

// CreateIndex creates a database index
func CreateIndex(db *gorm.DB, index IndexDefinition) error {
	if len(index.Columns) == 0 {
		return fmt.Errorf("index columns cannot be empty")
	}

	// Generate index name if not provided
	if index.Name == "" {
		index.Name = fmt.Sprintf("idx_%s_%s", index.Table, strings.Join(index.Columns, "_"))
	}

	// Build index SQL based on database type
	dbType := db.Dialector.Name()
	var sql string

	columns := strings.Join(index.Columns, ", ")

	switch dbType {
	case "postgres":
		unique := ""
		if index.Unique {
			unique = "UNIQUE "
		}

		concurrent := ""
		if index.Concurrent {
			concurrent = "CONCURRENTLY "
		}

		where := ""
		if index.Where != "" {
			where = fmt.Sprintf(" WHERE %s", index.Where)
		}

		sql = fmt.Sprintf(
			"CREATE %sINDEX %s%s ON %s (%s)%s",
			unique,
			concurrent,
			index.Name,
			index.Table,
			columns,
			where,
		)

	case "mysql":
		unique := ""
		if index.Unique {
			unique = "UNIQUE "
		}

		sql = fmt.Sprintf(
			"CREATE %sINDEX %s ON %s (%s)",
			unique,
			index.Name,
			index.Table,
			columns,
		)

	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	logger.Info("Creating index",
		zap.String("name", index.Name),
		zap.String("table", index.Table),
		zap.Strings("columns", index.Columns),
	)

	if err := db.Exec(sql).Error; err != nil {
		// Check if index already exists
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "Duplicate key") {
			logger.Warn("Index already exists, skipping",
				zap.String("name", index.Name),
				zap.String("table", index.Table),
			)
			return nil
		}
		return fmt.Errorf("failed to create index %s: %w", index.Name, err)
	}

	logger.Info("Index created successfully",
		zap.String("name", index.Name),
		zap.String("table", index.Table),
	)

	return nil
}

// DropIndex drops a database index
func DropIndex(db *gorm.DB, tableName, indexName string) error {
	dbType := db.Dialector.Name()
	var sql string

	switch dbType {
	case "postgres":
		sql = fmt.Sprintf("DROP INDEX IF EXISTS %s", indexName)
	case "mysql":
		sql = fmt.Sprintf("DROP INDEX %s ON %s", indexName, tableName)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	logger.Info("Dropping index",
		zap.String("name", indexName),
		zap.String("table", tableName),
	)

	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to drop index %s: %w", indexName, err)
	}

	logger.Info("Index dropped successfully",
		zap.String("name", indexName),
	)

	return nil
}

// CreateIndexes creates multiple indexes
func CreateIndexes(db *gorm.DB, indexes []IndexDefinition) error {
	for _, index := range indexes {
		if err := CreateIndex(db, index); err != nil {
			return fmt.Errorf("failed to create index %s: %w", index.Name, err)
		}
	}
	return nil
}

// IndexExists checks if an index exists
func IndexExists(db *gorm.DB, tableName, indexName string) (bool, error) {
	dbType := db.Dialector.Name()
	var sql string
	var result int

	switch dbType {
	case "postgres":
		sql = `
			SELECT COUNT(*) 
			FROM pg_indexes 
			WHERE tablename = $1 AND indexname = $2
		`
		if err := db.Raw(sql, tableName, indexName).Scan(&result).Error; err != nil {
			return false, err
		}

	case "mysql":
		sql = `
			SELECT COUNT(*) 
			FROM information_schema.statistics 
			WHERE table_schema = DATABASE() 
			AND table_name = ? 
			AND index_name = ?
		`
		if err := db.Raw(sql, tableName, indexName).Scan(&result).Error; err != nil {
			return false, err
		}

	default:
		return false, fmt.Errorf("unsupported database type: %s", dbType)
	}

	return result > 0, nil
}

// GetTableIndexes returns all indexes for a table
func GetTableIndexes(db *gorm.DB, tableName string) ([]string, error) {
	dbType := db.Dialector.Name()
	var indexes []string

	switch dbType {
	case "postgres":
		sql := `
			SELECT indexname 
			FROM pg_indexes 
			WHERE tablename = $1
		`
		if err := db.Raw(sql, tableName).Scan(&indexes).Error; err != nil {
			return nil, err
		}

	case "mysql":
		sql := `
			SELECT DISTINCT index_name 
			FROM information_schema.statistics 
			WHERE table_schema = DATABASE() 
			AND table_name = ?
		`
		if err := db.Raw(sql, tableName).Scan(&indexes).Error; err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	return indexes, nil
}

// RecommendedIndexes returns recommended indexes for common tables
func RecommendedIndexes() map[string][]IndexDefinition {
	return map[string][]IndexDefinition{
		"users": {
			{Name: "idx_users_email", Table: "users", Columns: []string{"email"}, Unique: true},
			{Name: "idx_users_created_at", Table: "users", Columns: []string{"created_at"}},
			{Name: "idx_users_updated_at", Table: "users", Columns: []string{"updated_at"}},
		},
		"posts": {
			{Name: "idx_posts_user_id", Table: "posts", Columns: []string{"user_id"}},
			{Name: "idx_posts_created_at", Table: "posts", Columns: []string{"created_at"}},
			{Name: "idx_posts_status_created", Table: "posts", Columns: []string{"status", "created_at"}},
		},
		"refresh_tokens": {
			{Name: "idx_refresh_tokens_user_id", Table: "refresh_tokens", Columns: []string{"user_id"}},
			{Name: "idx_refresh_tokens_token", Table: "refresh_tokens", Columns: []string{"token"}, Unique: true},
			{Name: "idx_refresh_tokens_expires_at", Table: "refresh_tokens", Columns: []string{"expires_at"}},
		},
		"email_verifications": {
			{Name: "idx_email_verifications_user_id", Table: "email_verifications", Columns: []string{"user_id"}},
			{Name: "idx_email_verifications_token", Table: "email_verifications", Columns: []string{"token"}, Unique: true},
			{Name: "idx_email_verifications_expires_at", Table: "email_verifications", Columns: []string{"expires_at"}},
		},
		"password_resets": {
			{Name: "idx_password_resets_user_id", Table: "password_resets", Columns: []string{"user_id"}},
			{Name: "idx_password_resets_token", Table: "password_resets", Columns: []string{"token"}, Unique: true},
			{Name: "idx_password_resets_expires_at", Table: "password_resets", Columns: []string{"expires_at"}},
		},
	}
}
