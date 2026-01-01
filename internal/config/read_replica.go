package config

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DBWrite is the primary database for write operations
	DBWrite *gorm.DB

	// DBRead is the read replica for read operations
	DBRead *gorm.DB
)

// LoadReadReplica initializes read replica connection
func LoadReadReplica() {
	// Check if read replica is configured
	readHost := ENV.DB_READ_HOST
	if readHost == "" {
		// No read replica configured, use primary DB for reads
		DBRead = DB
		return
	}

	readPort := ENV.DB_READ_PORT
	if readPort == "" {
		readPort = ENV.DB_URL // Fallback to primary port
	}
	readUser := ENV.DB_READ_USER
	if readUser == "" {
		readUser = ENV.DB_USER
	}
	readPass := ENV.DB_READ_PASS
	if readPass == "" {
		readPass = ENV.DB_PASS
	}
	readName := ENV.DB_READ_NAME
	if readName == "" {
		readName = ENV.DB_NAME
	}

	var dial gorm.Dialector
	switch ENV.DB_TYPE {
	case "mysql":
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
			readUser, readPass, readHost, readPort, readName)
		dial = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=UTC",
			readHost, readPort, readUser, readPass, readName)
		dial = postgres.Open(dsn)
	case "sql server":
		dsn := fmt.Sprintf("sqlserver://%v:%v@%v:%v?database=%v&connection+timeout=30",
			readUser, readPass, readHost, readPort, readName)
		dial = sqlserver.Open(dsn)
	default:
		panic("database type not supported for read replica")
	}

	// Configure GORM logger for read replica
	logConfig := &gorm.Config{
		Logger: NewGormLogger(
			1*time.Second,
			logger.Warn,
		),
	}

	db, err := gorm.Open(dial, logConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to read replica: %v", err))
	}

	// Configure connection pool for read replica
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// Read replicas typically handle more connections
	if ENV.ENV_TYPE == "prod" {
		sqlDB.SetMaxIdleConns(50)  // Higher idle for read operations
		sqlDB.SetMaxOpenConns(300) // More connections for reads
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	} else {
		sqlDB.SetMaxIdleConns(15)
		sqlDB.SetMaxOpenConns(75)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	}

	DBRead = db
	DBWrite = DB // Primary DB for writes
}

// UseReadReplica returns read replica if available, otherwise primary DB
func UseReadReplica() *gorm.DB {
	if DBRead != nil {
		return DBRead
	}
	return DB
}

// UseWriteDB returns primary DB for write operations
func UseWriteDB() *gorm.DB {
	if DBWrite != nil {
		return DBWrite
	}
	return DB
}

// DBResolver helps route queries to appropriate database
type DBResolver struct {
	WriteDB *gorm.DB
	ReadDB  *gorm.DB
}

// NewDBResolver creates a new database resolver
func NewDBResolver() *DBResolver {
	return &DBResolver{
		WriteDB: UseWriteDB(),
		ReadDB:  UseReadReplica(),
	}
}

// GetRead returns read database instance
func (r *DBResolver) GetRead() *gorm.DB {
	return r.ReadDB
}

// GetWrite returns write database instance
func (r *DBResolver) GetWrite() *gorm.DB {
	return r.WriteDB
}
