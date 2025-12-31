package helper

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// SortOrder represents sort direction
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// SortField represents a single sort field
type SortField struct {
	Field string    `json:"field" query:"field"`
	Order SortOrder `json:"order" query:"order"`
}

// SortConfig represents sorting configuration
type SortConfig struct {
	Fields        []SortField `json:"fields" query:"fields"`
	DefaultField  string      `json:"-"`
	DefaultOrder  SortOrder   `json:"-"`
	AllowedFields []string    `json:"-"` // Whitelist of allowed fields
}

// ApplySort applies single field sorting to query
func ApplySort(db *gorm.DB, field string, order SortOrder) *gorm.DB {
	if field == "" {
		return db
	}
	
	// Sanitize field name (prevent SQL injection)
	field = sanitizeFieldName(field)
	
	// Validate order
	if order != SortAsc && order != SortDesc {
		order = SortAsc
	}
	
	return db.Order(fmt.Sprintf("%s %s", field, strings.ToUpper(string(order))))
}

// ApplyMultiSort applies multiple field sorting to query
func ApplyMultiSort(db *gorm.DB, fields []SortField, allowedFields []string) *gorm.DB {
	if len(fields) == 0 {
		return db
	}
	
	for _, sortField := range fields {
		// Validate field is allowed
		if len(allowedFields) > 0 && !isFieldAllowed(sortField.Field, allowedFields) {
			continue
		}
		
		db = ApplySort(db, sortField.Field, sortField.Order)
	}
	
	return db
}

// ApplySortConfig applies sorting configuration to query
func ApplySortConfig(db *gorm.DB, config SortConfig) *gorm.DB {
	// If no fields provided, use default
	if len(config.Fields) == 0 && config.DefaultField != "" {
		return ApplySort(db, config.DefaultField, config.DefaultOrder)
	}
	
	return ApplyMultiSort(db, config.Fields, config.AllowedFields)
}

// ParseSortString parses sort string to SortField array
// Examples:
//   - "name" -> [{field: "name", order: "asc"}]
//   - "name:desc" -> [{field: "name", order: "desc"}]
//   - "name:asc,created_at:desc" -> [{field: "name", order: "asc"}, {field: "created_at", order: "desc"}]
//   - "-created_at" -> [{field: "created_at", order: "desc"}] (- prefix means descending)
func ParseSortString(sortStr string) []SortField {
	if sortStr == "" {
		return []SortField{}
	}
	
	fields := make([]SortField, 0)
	parts := strings.Split(sortStr, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		field := SortField{
			Order: SortAsc, // Default order
		}
		
		// Check for - prefix (descending)
		if strings.HasPrefix(part, "-") {
			field.Order = SortDesc
			part = part[1:]
		}
		
		// Check for :asc or :desc suffix
		if strings.Contains(part, ":") {
			subParts := strings.Split(part, ":")
			if len(subParts) == 2 {
				field.Field = subParts[0]
				if strings.ToLower(subParts[1]) == "desc" {
					field.Order = SortDesc
				}
			}
		} else {
			field.Field = part
		}
		
		if field.Field != "" {
			fields = append(fields, field)
		}
	}
	
	return fields
}

// ParseSortParams parses sort parameters from query string
func ParseSortParams(sortBy, order string) SortField {
	sortOrder := SortAsc
	if strings.ToLower(order) == "desc" {
		sortOrder = SortDesc
	}
	
	return SortField{
		Field: sortBy,
		Order: sortOrder,
	}
}

// ValidateSortField validates if sort field is allowed
func ValidateSortField(field string, allowedFields []string) error {
	if len(allowedFields) == 0 {
		return nil // No restrictions
	}
	
	if !isFieldAllowed(field, allowedFields) {
		return &BadRequestError{
			Message: fmt.Sprintf("Field '%s' is not allowed for sorting", field),
			Order:   "H-Sort-Validate-1",
		}
	}
	
	return nil
}

// ValidateSortFields validates multiple sort fields
func ValidateSortFields(fields []SortField, allowedFields []string) error {
	for _, field := range fields {
		if err := ValidateSortField(field.Field, allowedFields); err != nil {
			return err
		}
	}
	return nil
}

// sanitizeFieldName sanitizes field name to prevent SQL injection
func sanitizeFieldName(field string) string {
	// Remove any characters that could be used for SQL injection
	// Allow only alphanumeric, underscore, and dot (for table.field)
	var sanitized strings.Builder
	for _, r := range field {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '.' {
			sanitized.WriteRune(r)
		}
	}
	return sanitized.String()
}

// isFieldAllowed checks if field is in allowed list
func isFieldAllowed(field string, allowedFields []string) bool {
	field = sanitizeFieldName(field)
	for _, allowed := range allowedFields {
		if field == allowed {
			return true
		}
	}
	return false
}

// DefaultSortConfig creates default sort configuration
func DefaultSortConfig(defaultField string, defaultOrder SortOrder, allowedFields []string) SortConfig {
	return SortConfig{
		DefaultField:  defaultField,
		DefaultOrder:  defaultOrder,
		AllowedFields: allowedFields,
	}
}

// BuildSortFromQuery builds sort configuration from query parameters
// Supports:
//   - ?sort=name (ascending)
//   - ?sort=-created_at (descending with - prefix)
//   - ?sort=name:asc,created_at:desc (multiple fields)
//   - ?sort_by=name&order=desc (separate parameters)
func BuildSortFromQuery(params map[string]string, allowedFields []string) SortConfig {
	config := SortConfig{
		AllowedFields: allowedFields,
	}
	
	// Check for "sort" parameter (supports multiple fields)
	if sortStr, ok := params["sort"]; ok {
		config.Fields = ParseSortString(sortStr)
		return config
	}
	
	// Check for "sort_by" and "order" parameters (single field)
	if sortBy, ok := params["sort_by"]; ok {
		order := SortAsc
		if orderStr, ok := params["order"]; ok && strings.ToLower(orderStr) == "desc" {
			order = SortDesc
		}
		config.Fields = []SortField{
			{
				Field: sortBy,
				Order: order,
			},
		}
		return config
	}
	
	return config
}

// CombineSortAndPagination combines sorting with cursor pagination
func CombineSortAndPagination(db *gorm.DB, sortConfig SortConfig, pagination CursorPagination) *gorm.DB {
	// Apply sorting first
	db = ApplySortConfig(db, sortConfig)
	
	// Then apply cursor pagination
	db, _ = ApplyCursorPagination(db, pagination)
	
	return db
}
