package config

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	if ENV.DB_TYPE == "mysql" {
		dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", ENV.DB_USER, ENV.DB_PASS, ENV.DB_URL, ENV.DB_NAME)
		dial = mysql.Open(dsn)
	} else if ENV.DB_TYPE == "postgres" {
		dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=UTC", url[0], url[len(url)-1], ENV.DB_USER, ENV.DB_PASS, ENV.DB_NAME)
		dial = postgres.Open(dsn)
	} else {
		dsn := fmt.Sprintf("sqlserver://%v:%v@%v?database=%v&connection+timeout=30", ENV.DB_USER, ENV.DB_PASS, ENV.DB_URL, ENV.DB_NAME)
		dial = sqlserver.Open(dsn)
	}
	var logConfig *gorm.Config
	if ENV.ENV_TYPE == "dev" {
		logConfig = &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					LogLevel: logger.Info, // Set level log menjadi Info untuk menampilkan semua log query
				},
			),
		}
	} else {
		logConfig = &gorm.Config{}
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

		err = db.AutoMigrate(entity.User{}, entity.Post{})
		if err != nil {
			panic(err)
		}
	}

	DB = db
}

func LoadDB2() {
	url := strings.Split(ENV.DB_2_URL, ":")
	var dial gorm.Dialector
	if ENV.DB_2_TYPE == "mysql" {
		dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", ENV.DB_2_USER, ENV.DB_2_PASS, ENV.DB_2_URL, ENV.DB_2_NAME)
		dial = mysql.Open(dsn)
	} else {
		dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=UTC", url[0], url[len(url)-1], ENV.DB_2_USER, ENV.DB_2_PASS, ENV.DB_2_NAME)
		dial = postgres.Open(dsn)
	}
	var logConfig *gorm.Config
	if ENV.ENV_TYPE == "dev" {
		logConfig = &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					LogLevel: logger.Info, // Set level log menjadi Info untuk menampilkan semua log query
				},
			),
		}
	} else {
		logConfig = &gorm.Config{}
	}
	db, err := gorm.Open(dial, logConfig)
	if err != nil {
		panic(err)
	}

	if ENV.DB_GEN {
		err = createEnum(db)
		if err != nil {
			panic(err)
		}

		err = db.AutoMigrate(entity.User{}, entity.Post{})
		if err != nil {
			panic(err)
		}
	}
	DB2 = db
}
