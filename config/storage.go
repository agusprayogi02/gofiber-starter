package config

import (
	"fmt"
	"os"

	"github.com/gofiber/storage/sqlite3/v2"
)

var STORAGE *sqlite3.Storage

func LoadStorage() {
	path_store := fmt.Sprintf("./assets/%s_storage.db", ENV.DB_NAME)
	if _, err := os.Stat(path_store); os.IsNotExist(err) {
		os.Create(path_store)
	}
	store := sqlite3.New(sqlite3.Config{
		Database: path_store,
	})
	STORAGE = store
}
