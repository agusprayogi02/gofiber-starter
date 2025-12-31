package middleware

import (
	"sync"
	"time"

	"starter-gofiber/helper"

	"github.com/gofiber/fiber/v2"
)

// UserRateLimiter implements per-user rate limiting
type UserRateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*rateLimitEntry
	max      int
	window   time.Duration
}

type rateLimitEntry struct {
	count      int
	resetTime  time.Time
	lastAccess time.Time
}

// NewUserRateLimiter creates a new user-based rate limiter
func NewUserRateLimiter(max int, window time.Duration) *UserRateLimiter {
	limiter := &UserRateLimiter{
		limiters: make(map[string]*rateLimitEntry),
		max:      max,
		window:   window,
	}

	// Clean up expired entries every minute
	go limiter.cleanup()

	return limiter
}

// Middleware returns fiber middleware for per-user rate limiting
func (rl *UserRateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user identifier (user_id from JWT or IP)
		identifier := rl.getIdentifier(c)

		// Check rate limit
		allowed, resetTime := rl.allow(identifier)
		if !allowed {
			return &helper.TooManyRequestsError{
				Message: "Rate limit exceeded. Please try again later.",
				Order:   "RL1",
			}
		}

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", string(rune(rl.max)))
		c.Set("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

		return c.Next()
	}
}

// getIdentifier returns user_id from JWT or falls back to IP
func (rl *UserRateLimiter) getIdentifier(c *fiber.Ctx) string {
	// Try to get user from JWT token
	userClaims, err := helper.GetUserFromToken(c)
	if err == nil && userClaims != nil {
		return "user:" + userClaims.Email
	}

	// Fallback to IP address
	return "ip:" + c.IP()
}

// allow checks if request is allowed under rate limit
func (rl *UserRateLimiter) allow(identifier string) (bool, time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.limiters[identifier]

	// Create new entry if doesn't exist or window expired
	if !exists || now.After(entry.resetTime) {
		rl.limiters[identifier] = &rateLimitEntry{
			count:      1,
			resetTime:  now.Add(rl.window),
			lastAccess: now,
		}
		return true, now.Add(rl.window)
	}

	// Update last access time
	entry.lastAccess = now

	// Check if limit exceeded
	if entry.count >= rl.max {
		return false, entry.resetTime
	}

	// Increment counter
	entry.count++
	return true, entry.resetTime
}

// cleanup removes stale entries
func (rl *UserRateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for id, entry := range rl.limiters {
			// Remove if reset time passed and no access in last 5 minutes
			if now.After(entry.resetTime) && now.Sub(entry.lastAccess) > 5*time.Minute {
				delete(rl.limiters, id)
			}
		}
		rl.mu.Unlock()
	}
}
