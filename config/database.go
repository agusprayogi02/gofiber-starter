package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func LoadDB() {
	url := strings.Split(ENV.DB_URL, ":")
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=UTC", url[0], ENV.DB_USER, ENV.DB_PASS, ENV.DB_NAME, url[len(url)-1])
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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

	DB = db
}
