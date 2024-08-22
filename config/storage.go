package config

import (
	"fmt"

	"github.com/gofiber/storage/sqlite3/v2"
)

var STORAGE *sqlite3.Storage

func LoadStorage() {
	store := sqlite3.New(sqlite3.Config{
		Database: fmt.Sprintf("./assets/%s_storage.db", ENV.DB_NAME),
	})
	STORAGE = store
}
