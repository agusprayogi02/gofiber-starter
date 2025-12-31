package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ENV_TYPE        string
	PORT            string
	DB_USER         string
	DB_PASS         string
	DB_URL          string
	DB_NAME         string
	DB_TYPE         string
	DB_GEN          bool
	DB_2_ENABLE     bool
	DB_2_URL        string
	DB_2_NAME       string
	DB_2_TYPE       string
	DB_2_USER       string
	DB_2_PASS       string
	DB_2_GEN        bool
	TIMEZONE        string
	LOCATION_CERT   string
	NGROK_AUTHTOKEN string
	SENTRY_DSN      string
	ENCRYPTION_KEY  string

	// Redis Configuration
	REDIS_HOST     string
	REDIS_PORT     string
	REDIS_PASSWORD string
	REDIS_DB       int
	REDIS_ENABLE   bool
}

var ENV *Config

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		panic(err)
	}
}
