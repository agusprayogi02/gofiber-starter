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
	LOCATION_CERT   string
	NGROK_AUTHTOKEN string
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
