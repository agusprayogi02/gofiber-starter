package middleware

import (
	"strings"
	"sync"

	"starter-gofiber/pkg/apierror"
	"starter-gofiber/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// IPFilterConfig holds IP whitelist/blacklist configuration
type IPFilterConfig struct {
	Whitelist []string
	Blacklist []string
	mu        sync.RWMutex
}

var globalIPFilter = &IPFilterConfig{
	Whitelist: []string{},
	Blacklist: []string{},
}

// IPWhitelistMiddleware allows only whitelisted IPs
func IPWhitelistMiddleware(whitelist []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := c.IP()

		// Check if IP is in whitelist
		if !isIPInList(clientIP, whitelist) {
			logger.Warn("Access denied for IP: " + clientIP)
			return &apierror.ForbiddenError{
				Message: "Access denied",
				Order:   "IP1",
			}
		}

		return c.Next()
	}
}

// IPBlacklistMiddleware blocks blacklisted IPs
func IPBlacklistMiddleware(blacklist []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := c.IP()

		// Check if IP is in blacklist
		if isIPInList(clientIP, blacklist) {
			logger.Warn("Blocked IP attempting access: " + clientIP)
			return &apierror.ForbiddenError{
				Message: "Access denied",
				Order:   "IP2",
			}
		}

		return c.Next()
	}
}

// DynamicIPFilterMiddleware uses global filter config
func DynamicIPFilterMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := c.IP()

		globalIPFilter.mu.RLock()
		defer globalIPFilter.mu.RUnlock()

		// Check whitelist first (if configured)
		if len(globalIPFilter.Whitelist) > 0 {
			if !isIPInList(clientIP, globalIPFilter.Whitelist) {
				logger.Warn("IP not in whitelist: " + clientIP)
				return &apierror.ForbiddenError{
					Message: "Access denied",
					Order:   "IP3",
				}
			}
		}

		// Check blacklist
		if isIPInList(clientIP, globalIPFilter.Blacklist) {
			logger.Warn("Blacklisted IP: " + clientIP)
			return &apierror.ForbiddenError{
				Message: "Access denied",
				Order:   "IP4",
			}
		}

		return c.Next()
	}
}

// AddToWhitelist adds IP to whitelist
func AddToWhitelist(ip string) {
	globalIPFilter.mu.Lock()
	defer globalIPFilter.mu.Unlock()

	if !isIPInList(ip, globalIPFilter.Whitelist) {
		globalIPFilter.Whitelist = append(globalIPFilter.Whitelist, ip)
		logger.Info("Added IP to whitelist: " + ip)
	}
}

// RemoveFromWhitelist removes IP from whitelist
func RemoveFromWhitelist(ip string) {
	globalIPFilter.mu.Lock()
	defer globalIPFilter.mu.Unlock()

	globalIPFilter.Whitelist = removeIP(globalIPFilter.Whitelist, ip)
	logger.Info("Removed IP from whitelist: " + ip)
}

// AddToBlacklist adds IP to blacklist
func AddToBlacklist(ip string) {
	globalIPFilter.mu.Lock()
	defer globalIPFilter.mu.Unlock()

	if !isIPInList(ip, globalIPFilter.Blacklist) {
		globalIPFilter.Blacklist = append(globalIPFilter.Blacklist, ip)
		logger.Info("Added IP to blacklist: " + ip)
	}
}

// RemoveFromBlacklist removes IP from blacklist
func RemoveFromBlacklist(ip string) {
	globalIPFilter.mu.Lock()
	defer globalIPFilter.mu.Unlock()

	globalIPFilter.Blacklist = removeIP(globalIPFilter.Blacklist, ip)
	logger.Info("Removed IP from blacklist: " + ip)
}

// GetWhitelist returns current whitelist
func GetWhitelist() []string {
	globalIPFilter.mu.RLock()
	defer globalIPFilter.mu.RUnlock()

	whitelist := make([]string, len(globalIPFilter.Whitelist))
	copy(whitelist, globalIPFilter.Whitelist)
	return whitelist
}

// GetBlacklist returns current blacklist
func GetBlacklist() []string {
	globalIPFilter.mu.RLock()
	defer globalIPFilter.mu.RUnlock()

	blacklist := make([]string, len(globalIPFilter.Blacklist))
	copy(blacklist, globalIPFilter.Blacklist)
	return blacklist
}

// isIPInList checks if IP is in the list (supports wildcards and CIDR)
func isIPInList(ip string, list []string) bool {
	for _, allowedIP := range list {
		// Exact match
		if ip == allowedIP {
			return true
		}

		// Wildcard match (e.g., 192.168.1.*)
		if strings.Contains(allowedIP, "*") {
			pattern := strings.ReplaceAll(allowedIP, "*", ".*")
			pattern = "^" + pattern + "$"
			// Simple wildcard matching (for production, use proper CIDR)
			if matchWildcard(ip, allowedIP) {
				return true
			}
		}

		// CIDR support can be added for production:
		// - Parse CIDR notation (e.g., "192.168.1.0/24")
		// - Check if IP falls within CIDR range
		// - Use net.ParseCIDR() and net.Contains() for validation
	}
	return false
}

// matchWildcard performs simple wildcard matching
func matchWildcard(ip, pattern string) bool {
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return ip == pattern
	}

	// Check prefix
	if !strings.HasPrefix(ip, parts[0]) {
		return false
	}

	// Check suffix
	if len(parts) > 1 && parts[len(parts)-1] != "" {
		if !strings.HasSuffix(ip, parts[len(parts)-1]) {
			return false
		}
	}

	return true
}

// removeIP removes IP from slice
func removeIP(list []string, ip string) []string {
	result := []string{}
	for _, item := range list {
		if item != ip {
			result = append(result, item)
		}
	}
	return result
}
