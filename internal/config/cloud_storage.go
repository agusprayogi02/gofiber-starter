package config

import (
	"os"

	"starter-gofiber/internal/infrastructure/storage"
)

// GetS3Config returns S3 configuration from environment variables
func GetS3Config() storage.S3Config {
	return storage.S3Config{
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region:    os.Getenv("AWS_REGION"),
		Bucket:    os.Getenv("AWS_S3_BUCKET"),
		Endpoint:  os.Getenv("AWS_S3_ENDPOINT"), // Optional, for S3-compatible services
	}
}

// InitS3Client initializes S3 client
func InitS3Client() (*storage.S3Client, error) {
	cfg := GetS3Config()
	return storage.NewS3Client(cfg)
}
