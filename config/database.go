package config

import (
	"fmt"
	"strings"
	"time"

	"starter-gofiber/entity"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB, DB2 *gorm.DB

func createEnum(db *gorm.DB) error {
	return db.Exec(`
        DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
                CREATE TYPE user_role AS ENUM ('admin', 'user');
            END IF;
        END$$;
    `).Error
}

func LoadDB() {
	url := strings.Split(ENV.DB_URL, ":")
	var dial gorm.Dialector
	switch ENV.DB_TYPE {
	case "mysql":
		dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", ENV.DB_USER, ENV.DB_PASS, ENV.DB_URL, ENV.DB_NAME)
		dial = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=UTC", url[0], url[len(url)-1], ENV.DB_USER, ENV.DB_PASS, ENV.DB_NAME)
		dial = postgres.Open(dsn)
	case "sql server":
		dsn := fmt.Sprintf("sqlserver://%v:%v@%v?database=%v&connection+timeout=30&TrustServerCertificate=false&encrypt=true", ENV.DB_USER, ENV.DB_PASS, ENV.DB_URL, ENV.DB_NAME)
		dial = sqlserver.Open(dsn)
	default:
		panic("database is not supported")
	}

	// Configure GORM with structured logging
	var logConfig *gorm.Config
	if ENV.ENV_TYPE == "dev" {
		logConfig = &gorm.Config{
			Logger: NewGormLogger(
				200*time.Millisecond, // Slow query threshold
				logger.Info,          // Log all queries in dev
			),
		}
	} else {
		logConfig = &gorm.Config{
			Logger: NewGormLogger(
				1*time.Second, // Slow query threshold for production
				logger.Warn,   // Only log warnings and errors in production
			),
		}
	}

	db, err := gorm.Open(dial, logConfig)
	if err != nil {
		panic(err)
	}

	if ENV.DB_GEN {
		if ENV.DB_TYPE == "postgres" {
			err = createEnum(db)
			if err != nil {
				panic(err)
			}
		}

		err = db.AutoMigrate(
			entity.User{},
			entity.Post{},
			entity.RefreshToken{},
			entity.PasswordReset{},
			entity.EmailVerification{},
			entity.APIKey{},
		)
		if err != nil {
			panic(err)
		}
	}

	// Configure connection pool with optimized settings
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// Set connection pool settings based on environment
	if ENV.ENV_TYPE == "prod" {
		// Production settings - higher limits
		sqlDB.SetMaxIdleConns(25)           // Minimum idle connections
		sqlDB.SetMaxOpenConns(200)          // Maximum open connections
		sqlDB.SetConnMaxLifetime(time.Hour) // Connection lifetime
		sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	} else {
		// Development settings - conservative limits
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	}

	DB = db
}

func LoadDB2() {
	url := strings.Split(ENV.DB_2_URL, ":")
	var dial gorm.Dialector
	switch ENV.DB_2_TYPE {
	case "mysql":
		dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", ENV.DB_2_USER, ENV.DB_2_PASS, ENV.DB_2_URL, ENV.DB_2_NAME)
		dial = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=UTC", url[0], url[len(url)-1], ENV.DB_2_USER, ENV.DB_2_PASS, ENV.DB_2_NAME)
		dial = postgres.Open(dsn)
	case "sql server":
		dsn := fmt.Sprintf("sqlserver://%v:%v@%v?database=%v&connection+timeout=30", ENV.DB_2_USER, ENV.DB_2_PASS, ENV.DB_2_URL, ENV.DB_2_NAME)
		dial = sqlserver.Open(dsn)
	default:
		panic("database is not supported")
	}

	// Configure GORM with structured logging for DB2
	var logConfig *gorm.Config
	if ENV.ENV_TYPE == "dev" {
		logConfig = &gorm.Config{
			Logger: NewGormLogger(
				200*time.Millisecond,
				logger.Info,
			),
		}
	} else {
		logConfig = &gorm.Config{
			Logger: NewGormLogger(
				1*time.Second,
				logger.Warn,
			),
		}
	}

	db, err := gorm.Open(dial, logConfig)
	if err != nil {
		panic(err)
	}

	if ENV.DB_2_GEN {
		if ENV.DB_2_TYPE == "postgres" {
			err = createEnum(db)
			if err != nil {
				panic(err)
			}
		}

		err = db.AutoMigrate(
			entity.User{},
			entity.Post{},
			entity.RefreshToken{},
			entity.PasswordReset{},
			entity.EmailVerification{},
			entity.APIKey{},
		)
		if err != nil {
			panic(err)
		}
	}

	// Configure connection pool for DB2 with optimized settings
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// Set connection pool settings based on environment
	if ENV.ENV_TYPE == "prod" {
		// Production settings - higher limits
		sqlDB.SetMaxIdleConns(25)
		sqlDB.SetMaxOpenConns(200)
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	} else {
		// Development settings - conservative limits
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	}

	DB2 = db
}
