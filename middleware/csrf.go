package middleware

import (
	"crypto/subtle"
	"time"

	"starter-gofiber/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
)

// CSRFConfig holds CSRF middleware configuration
func CSRFMiddleware() fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "csrf_token",
		CookieSameSite: "Strict",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		Expiration:     24 * time.Hour,
		KeyGenerator:   utils.UUIDv4,
		Extractor:      csrf.CsrfFromHeader("X-CSRF-Token"),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return &helper.ForbiddenError{
				Message: "CSRF token validation failed",
				Order:   "CSRF1",
			}
		},
	})
}

// GetCSRFToken returns CSRF token for frontend
func GetCSRFToken(c *fiber.Ctx) string {
	token := c.Locals("csrf")
	if token == nil {
		return ""
	}
	return token.(string)
}

// ValidateCSRFToken manually validates CSRF token
func ValidateCSRFToken(c *fiber.Ctx, token string) bool {
	sessionToken := c.Locals("csrf")
	if sessionToken == nil {
		return false
	}

	// Constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(
		[]byte(token),
		[]byte(sessionToken.(string)),
	) == 1
}
