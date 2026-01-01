package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeadersMiddleware adds security headers (Helmet-like)
func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// X-Frame-Options: Prevent clickjacking
		c.Set("X-Frame-Options", "SAMEORIGIN")

		// X-Content-Type-Options: Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// X-XSS-Protection: Enable browser XSS protection
		c.Set("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy: Control referrer information
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy: Control browser features
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Content-Security-Policy: Prevent XSS, injection attacks
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self'; " +
			"frame-ancestors 'self'; " +
			"base-uri 'self'; " +
			"form-action 'self'"
		c.Set("Content-Security-Policy", csp)

		// Strict-Transport-Security: Force HTTPS (only in production with HTTPS)
		if c.Protocol() == "https" {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		// X-Permitted-Cross-Domain-Policies: Restrict cross-domain policies
		c.Set("X-Permitted-Cross-Domain-Policies", "none")

		// X-DNS-Prefetch-Control: Control DNS prefetching
		c.Set("X-DNS-Prefetch-Control", "off")

		// X-Download-Options: Prevent opening files in browser context
		c.Set("X-Download-Options", "noopen")

		return c.Next()
	}
}

// CustomCSPMiddleware allows custom CSP configuration
func CustomCSPMiddleware(policy string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", policy)
		return c.Next()
	}
}
