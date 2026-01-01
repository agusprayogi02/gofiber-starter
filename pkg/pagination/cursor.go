package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// ValidationError represents a validation error for pagination
type ValidationError struct {
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// CursorPagination represents cursor-based pagination parameters
type CursorPagination struct {
	Cursor        string   `json:"cursor" query:"cursor"`         // Cursor for next page
	Limit         int      `json:"limit" query:"limit"`           // Number of items per page
	SortBy        string   `json:"sort_by" query:"sort_by"`       // Field to sort by
	SortOrder     string   `json:"sort_order" query:"sort_order"` // asc or desc
	AllowedFields []string `json:"-"`                             // Whitelist of allowed fields (optional)
}

// DefaultAllowedFields returns default safe fields for sorting
var DefaultAllowedFields = []string{"id", "created_at", "updated_at"}

// CursorResponse represents paginated response with cursor
type CursorResponse struct {
	Data       interface{} `json:"data"`
	NextCursor string      `json:"next_cursor,omitempty"`
	HasMore    bool        `json:"has_more"`
	Count      int         `json:"count"`
}

// CursorData represents cursor internal data
type CursorData struct {
	LastID    uint   `json:"last_id"`
	LastValue string `json:"last_value,omitempty"` // For non-ID sorting
}

// DefaultCursorPagination returns default cursor pagination params
func DefaultCursorPagination() CursorPagination {
	return CursorPagination{
		Limit:     10,
		SortBy:    "id",
		SortOrder: "desc",
	}
}

// EncodeCursor encodes cursor data to base64 string
func EncodeCursor(lastID uint, lastValue string) string {
	cursor := CursorData{
		LastID:    lastID,
		LastValue: lastValue,
	}

	jsonData, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(jsonData)
}

// DecodeCursor decodes cursor from base64 string
func DecodeCursor(cursorStr string) (*CursorData, error) {
	if cursorStr == "" {
		return nil, nil
	}

	decoded, err := base64.URLEncoding.DecodeString(cursorStr)
	if err != nil {
		return nil, &ValidationError{
			Message: "Invalid cursor format",
			Field:   "cursor",
		}
	}

	var cursor CursorData
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return nil, &ValidationError{
			Message: "Invalid cursor data",
			Field:   "cursor",
		}
	}

	return &cursor, nil
}

// sanitizeFieldName sanitizes field name to prevent SQL injection
// Only allows alphanumeric characters, underscore, and dot (for table.field notation)
func sanitizeFieldName(field string) string {
	// Remove any characters that could be used for SQL injection
	var sanitized strings.Builder
	for _, r := range field {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '.' {
			sanitized.WriteRune(r)
		}
	}
	return sanitized.String()
}

// validateFieldName validates that field name matches safe pattern
func validateFieldName(field string) bool {
	// Must match: alphanumeric, underscore, dot (for table.field)
	// Must not start with number or dot
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_.]*$`, field)
	return matched
}

// isFieldAllowed checks if field is in allowed whitelist
func isFieldAllowed(field string, allowedFields []string) bool {
	if len(allowedFields) == 0 {
		// If no whitelist provided, use default safe fields
		allowedFields = DefaultAllowedFields
	}

	sanitized := sanitizeFieldName(field)
	for _, allowed := range allowedFields {
		if sanitized == allowed {
			return true
		}
	}
	return false
}

// validateAndSanitizeSortBy validates and sanitizes SortBy field
// Returns sanitized field name and error if validation fails
func validateAndSanitizeSortBy(sortBy string, allowedFields []string) (string, error) {
	// Sanitize first - removes dangerous characters
	sanitized := sanitizeFieldName(sortBy)

	// If sanitization removed everything, it was malicious
	if sanitized == "" && sortBy != "" {
		return "id", &ValidationError{
			Message: fmt.Sprintf("Invalid sort field format: %s", sortBy),
			Field:   "sort_by",
		}
	}

	// Validate pattern - must match safe identifier pattern
	if !validateFieldName(sanitized) {
		return "id", &ValidationError{
			Message: fmt.Sprintf("Invalid sort field format: %s", sortBy),
			Field:   "sort_by",
		}
	}

	// Check against whitelist
	if !isFieldAllowed(sanitized, allowedFields) {
		allowedList := DefaultAllowedFields
		if len(allowedFields) > 0 {
			allowedList = allowedFields
		}
		return "id", &ValidationError{
			Message: fmt.Sprintf("Field '%s' is not allowed for sorting. Allowed fields: %v", sortBy, allowedList),
			Field:   "sort_by",
		}
	}

	return sanitized, nil
}

// ApplyCursorPagination applies cursor pagination to GORM query
func ApplyCursorPagination(db *gorm.DB, pagination CursorPagination) (*gorm.DB, error) {
	// Validate limit
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}
	if pagination.Limit > 100 {
		pagination.Limit = 100 // Max limit
	}

	// Validate sort order
	if pagination.SortOrder != "asc" && pagination.SortOrder != "desc" {
		pagination.SortOrder = "desc"
	}

	// Validate and sanitize SortBy field
	if pagination.SortBy == "" {
		pagination.SortBy = "id"
	} else {
		var err error
		pagination.SortBy, err = validateAndSanitizeSortBy(pagination.SortBy, pagination.AllowedFields)
		if err != nil {
			// Return error - caller should handle it
			// This prevents SQL injection by rejecting invalid fields
			return db, err
		}
	}

	// Decode cursor
	cursor, err := DecodeCursor(pagination.Cursor)
	if err != nil {
		return db, err
	}

	// Apply cursor condition - use parameterized queries for field name safety
	if cursor != nil {
		// Use GORM's Where with proper parameterization
		// Note: We've already validated SortBy, so it's safe to use in format string
		// But we still use parameterized values for the actual data
		if pagination.SortOrder == "desc" {
			if cursor.LastValue != "" {
				// Sort by custom field - field name is validated, values are parameterized
				db = db.Where(fmt.Sprintf("%s < ? OR (%s = ? AND id < ?)",
					pagination.SortBy, pagination.SortBy),
					cursor.LastValue, cursor.LastValue, cursor.LastID)
			} else {
				// Sort by ID only
				db = db.Where("id < ?", cursor.LastID)
			}
		} else {
			if cursor.LastValue != "" {
				db = db.Where(fmt.Sprintf("%s > ? OR (%s = ? AND id > ?)",
					pagination.SortBy, pagination.SortBy),
					cursor.LastValue, cursor.LastValue, cursor.LastID)
			} else {
				db = db.Where("id > ?", cursor.LastID)
			}
		}
	}

	// Apply sorting and limit - field name is validated, order is validated
	orderClause := fmt.Sprintf("%s %s, id %s", pagination.SortBy, strings.ToUpper(pagination.SortOrder), strings.ToUpper(pagination.SortOrder))
	db = db.Order(orderClause).Limit(pagination.Limit + 1) // +1 to check if has more

	return db, nil
}

// BuildCursorResponse builds cursor response from query results
func BuildCursorResponse(data interface{}, pagination CursorPagination, lastID uint, lastValue string) CursorResponse {
	// Check if we need reflection to count items
	// For simplicity, assume data is a slice
	count := getSliceLength(data)

	hasMore := count > pagination.Limit
	if hasMore {
		count = pagination.Limit
		// Remove last item (used only for hasMore check)
		data = sliceData(data, pagination.Limit)
	}

	nextCursor := ""
	if hasMore && lastID > 0 {
		nextCursor = EncodeCursor(lastID, lastValue)
	}

	return CursorResponse{
		Data:       data,
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Count:      count,
	}
}

// Helper to get slice length using type assertion
func getSliceLength(data interface{}) int {
	switch v := data.(type) {
	case []interface{}:
		return len(v)
	default:
		// Try reflection if needed
		return 0
	}
}

// Helper to slice data
func sliceData(data interface{}, limit int) interface{} {
	switch v := data.(type) {
	case []interface{}:
		if len(v) > limit {
			return v[:limit]
		}
		return v
	default:
		return data
	}
}

// ParseCursorParams parses cursor pagination params from query string
// allowedFields: optional whitelist of allowed sort fields
func ParseCursorParams(cursor, limit, sortBy, sortOrder string, allowedFields ...string) CursorPagination {
	pagination := DefaultCursorPagination()

	pagination.Cursor = cursor

	if len(allowedFields) > 0 {
		pagination.AllowedFields = allowedFields
	}

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			pagination.Limit = l
		}
	}

	if sortBy != "" {
		// Validation will happen in ApplyCursorPagination
		pagination.SortBy = sortBy
	}

	if sortOrder == "asc" || sortOrder == "desc" {
		pagination.SortOrder = sortOrder
	}

	return pagination
}
