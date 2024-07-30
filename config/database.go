package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"starter-gofiber/entity"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

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
	} else {
		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=UTC", url[0], ENV.DB_USER, ENV.DB_PASS, ENV.DB_NAME, url[len(url)-1])
		dial = postgres.Open(dsn)
	}
	db, err := gorm.Open(dial, &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // Menggunakan logger bawaan dari Golang
			logger.Config{
				LogLevel: logger.Info, // Set level log menjadi Info untuk menampilkan semua log query
			},
		),
	})
	if err != nil {
		panic(err)
	}

	err = createEnum(db)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(entity.User{})
	if err != nil {
		panic(err)
	}

	DB = db
}
