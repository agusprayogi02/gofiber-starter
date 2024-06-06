package config

import (
	"crypto/rsa"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Config struct {
	PORT            string
	DB_USER         string
	DB_PASS         string
	DB_URL          string
	DB_NAME         string
	LOCATION_CERT   string
	NGROK_AUTHTOKEN string
}

var (
	ENV         *Config
	PRIVATE_KEY *rsa.PrivateKey
)

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

	privateKeyData, err := os.ReadFile(ENV.LOCATION_CERT)
	if err != nil {
		log.Fatalf("Error reading private key: %v", err)
	}

	PRIVATE_KEY, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
	}
}
