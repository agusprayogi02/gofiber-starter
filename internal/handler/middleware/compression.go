package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

// CompressionConfig represents compression middleware configuration
type CompressionConfig struct {
	Level int // Compression level (0-2): 0=Default, 1=Best Speed, 2=Best Compression
}

// DefaultCompressionConfig returns default compression configuration
func DefaultCompressionConfig() CompressionConfig {
	return CompressionConfig{
		Level: int(compress.LevelDefault), // 0 = default balanced compression
	}
}

// Compression returns a compression middleware with custom config
func Compression(config CompressionConfig) fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.Level(config.Level),
	})
}

// CompressionDefault returns a compression middleware with default config
func CompressionDefault() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelDefault,
	})
}

// CompressionBestSpeed returns a compression middleware optimized for speed
func CompressionBestSpeed() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
}

// CompressionBestSize returns a compression middleware optimized for size
func CompressionBestSize() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	})
}
