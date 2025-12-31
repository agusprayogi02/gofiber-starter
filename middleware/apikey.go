package middleware

import (
	"strings"

	"starter-gofiber/helper"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var apiKeyDB *gorm.DB

// InitAPIKeyMiddleware sets the database instance
func InitAPIKeyMiddleware(db *gorm.DB) {
	apiKeyDB = db
}

// APIKeyAuth middleware validates API key from header
func APIKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get API key from header
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return &helper.UnauthorizedError{
				Message: "API key required",
				Order:   "AK1",
			}
		}

		// Validate API key
		valid, userID := helper.ValidateAPIKey(apiKeyDB, apiKey)
		if !valid {
			return &helper.UnauthorizedError{
				Message: "Invalid API key",
				Order:   "AK2",
			}
		}

		// Store user ID in context
		c.Locals("api_user_id", userID)
		c.Locals("auth_method", "api_key")

		return c.Next()
	}
}

// OptionalAPIKeyAuth tries API key first, falls back to JWT
func OptionalAPIKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try API key first
		apiKey := c.Get("X-API-Key")
		if apiKey != "" {
			valid, userID := helper.ValidateAPIKey(apiKeyDB, apiKey)
			if valid {
				c.Locals("api_user_id", userID)
				c.Locals("auth_method", "api_key")
				return c.Next()
			}
		}

		// Fall back to JWT auth
		return AuthMiddleware()(c)
	}
}

// APIKeyOrJWT allows either API key or JWT authentication
func APIKeyOrJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for API key
		apiKey := c.Get("X-API-Key")
		if apiKey != "" {
			return APIKeyAuth()(c)
		}

		// Check for JWT
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			return AuthMiddleware()(c)
		}

		return &helper.UnauthorizedError{
			Message: "Authentication required (API key or Bearer token)",
			Order:   "AK3",
		}
	}
}
