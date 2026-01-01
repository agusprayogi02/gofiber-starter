package validator

import (
	"html"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	// Strict sanitizer - strips all HTML
	strictPolicy = bluemonday.StrictPolicy()

	// UGC (User Generated Content) sanitizer - allows safe HTML
	ugcPolicy = bluemonday.UGCPolicy()

	// SQL injection patterns
	sqlInjectionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute|script|javascript|alert)`),
		regexp.MustCompile(`(?i)(\||;|--|/\*|\*/|xp_|sp_)`),
	}

	// XSS patterns
	xssPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
		regexp.MustCompile(`(?i)javascript:`),
		regexp.MustCompile(`(?i)on\w+\s*=`), // onclick, onload, etc
		regexp.MustCompile(`(?i)<iframe[^>]*>`),
	}
)

// SanitizeInput removes all HTML tags and dangerous characters
func SanitizeInput(input string) string {
	// Remove all HTML tags
	sanitized := strictPolicy.Sanitize(input)

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// SanitizeHTML allows safe HTML (for user-generated content like comments)
func SanitizeHTML(input string) string {
	return ugcPolicy.Sanitize(input)
}

// EscapeHTML escapes HTML special characters
func EscapeHTML(input string) string {
	return html.EscapeString(input)
}

// UnescapeHTML unescapes HTML special characters
func UnescapeHTML(input string) string {
	return html.UnescapeString(input)
}

// ValidateNoSQLInjection checks for SQL injection patterns
func ValidateNoSQLInjection(input string) bool {
	for _, pattern := range sqlInjectionPatterns {
		if pattern.MatchString(input) {
			return false
		}
	}
	return true
}

// ValidateNoXSS checks for XSS patterns
func ValidateNoXSS(input string) bool {
	for _, pattern := range xssPatterns {
		if pattern.MatchString(input) {
			return false
		}
	}
	return true
}

// SanitizeEmail validates and sanitizes email
func SanitizeEmail(email string) string {
	// Convert to lowercase
	email = strings.ToLower(email)

	// Trim whitespace
	email = strings.TrimSpace(email)

	// Basic sanitization
	email = SanitizeInput(email)

	return email
}

// SanitizeString performs comprehensive sanitization
func SanitizeString(input string) string {
	// Remove HTML
	sanitized := SanitizeInput(input)

	// Remove null bytes
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")

	// Normalize whitespace
	sanitized = strings.Join(strings.Fields(sanitized), " ")

	return sanitized
}

// IsSafe checks if input is safe (no XSS, no SQL injection)
func IsSafe(input string) bool {
	return ValidateNoXSS(input) && ValidateNoSQLInjection(input)
}

// SanitizeMap sanitizes all string values in a map
func SanitizeMap(data map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})

	for key, value := range data {
		switch v := value.(type) {
		case string:
			sanitized[key] = SanitizeString(v)
		case map[string]interface{}:
			sanitized[key] = SanitizeMap(v)
		default:
			sanitized[key] = value
		}
	}

	return sanitized
}

// SanitizeSlice sanitizes all string values in a slice
func SanitizeSlice(data []string) []string {
	sanitized := make([]string, len(data))

	for i, value := range data {
		sanitized[i] = SanitizeString(value)
	}

	return sanitized
}
