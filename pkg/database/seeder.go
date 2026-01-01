package database

import (
	"fmt"

	"starter-gofiber/entity"
	"starter-gofiber/pkg/crypto"
	"starter-gofiber/pkg/logger"

	"gorm.io/gorm"
)

// Seeder represents a database seeder
type Seeder struct {
	Name string
	Run  func(*gorm.DB) error
}

var seeders []Seeder

// RegisterSeeder registers a new seeder
func RegisterSeeder(name string, run func(*gorm.DB) error) {
	seeders = append(seeders, Seeder{
		Name: name,
		Run:  run,
	})
}

// RunAllSeeders runs all registered seeders
func RunAllSeeders(db *gorm.DB) error {
	logger.Info("Running database seeders...")

	for _, seeder := range seeders {
		logger.Info(fmt.Sprintf("Running seeder: %s", seeder.Name))
		if err := seeder.Run(db); err != nil {
			return fmt.Errorf("failed to run seeder %s: %w", seeder.Name, err)
		}
	}

	logger.Info("All seeders completed successfully")
	return nil
}

// RunSeeder runs a specific seeder by name
func RunSeeder(db *gorm.DB, name string) error {
	for _, seeder := range seeders {
		if seeder.Name == name {
			logger.Info(fmt.Sprintf("Running seeder: %s", seeder.Name))
			return seeder.Run(db)
		}
	}
	return fmt.Errorf("seeder not found: %s", name)
}

// SeedUsers seeds sample users
func SeedUsers(db *gorm.DB) error {
	users := []entity.User{
		{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: mustHashPassword("password123"),
		},
		{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: mustHashPassword("password123"),
		},
		{
			Name:     "Demo User",
			Email:    "demo@example.com",
			Password: mustHashPassword("password123"),
		},
	}

	for _, user := range users {
		// Check if user already exists
		var existing entity.User
		result := db.Where("email = ?", user.Email).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			// Create new user
			if err := db.Create(&user).Error; err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Created user: %s", user.Email))
		} else {
			logger.Info(fmt.Sprintf("User already exists: %s", user.Email))
		}
	}

	return nil
}

// SeedPosts seeds sample posts
func SeedPosts(db *gorm.DB) error {
	// Get first user
	var user entity.User
	if err := db.First(&user).Error; err != nil {
		return fmt.Errorf("no users found, run user seeder first")
	}

	posts := []entity.Post{
		{
			Tweet:  "Getting Started with Go - A comprehensive guide to getting started with Go programming language",
			UserID: user.ID,
		},
		{
			Tweet:  "Building REST APIs with Fiber - Learn how to build fast and scalable REST APIs using Go Fiber framework",
			UserID: user.ID,
		},
		{
			Tweet:  "Database Best Practices - Essential best practices for working with databases in Go applications",
			UserID: user.ID,
		},
		{
			Tweet:  "Testing in Go - A deep dive into unit testing, integration testing, and benchmarking in Go",
			UserID: user.ID,
		},
		{
			Tweet:  "Microservices Architecture - Design patterns and best practices for building microservices with Go",
			UserID: user.ID,
		},
	}

	for _, post := range posts {
		// Check if post already exists
		var existing entity.Post
		result := db.Where("tweet = ?", post.Tweet).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&post).Error; err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Created post: %s", truncateText(post.Tweet, 50)))
		} else {
			logger.Info(fmt.Sprintf("Post already exists: %s", truncateText(post.Tweet, 50)))
		}
	}

	return nil
}

// mustHashPassword hashes password and panics on error (for seeder only)
func mustHashPassword(password string) string {
	hashed, err := crypto.HashPassword(password)
	if err != nil {
		panic(err)
	}
	return hashed
}

// truncateText truncates text to max length
func truncateText(text string, max int) string {
	if len(text) <= max {
		return text
	}
	return text[:max] + "..."
}

// TruncateTable truncates a table (use with caution!)
func TruncateTable(db *gorm.DB, tableName string) error {
	return db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error
}

// ResetDatabase drops all data and re-runs seeders (development only!)
func ResetDatabase(db *gorm.DB) error {
	logger.Info("Resetting database...")

	// Truncate tables in correct order (respect foreign keys)
	tables := []string{"posts", "api_keys", "users"}

	for _, table := range tables {
		if err := TruncateTable(db, table); err != nil {
			logger.Warn(fmt.Sprintf("Failed to truncate table %s: %v", table, err))
		}
	}

	// Re-run all seeders
	return RunAllSeeders(db)
}

// init registers default seeders
func init() {
	RegisterSeeder("users", SeedUsers)
	RegisterSeeder("posts", SeedPosts)
}
