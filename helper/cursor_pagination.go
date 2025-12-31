package helper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// CursorPagination represents cursor-based pagination parameters
type CursorPagination struct {
	Cursor    string `json:"cursor" query:"cursor"`       // Cursor for next page
	Limit     int    `json:"limit" query:"limit"`         // Number of items per page
	SortBy    string `json:"sort_by" query:"sort_by"`     // Field to sort by
	SortOrder string `json:"sort_order" query:"sort_order"` // asc or desc
}

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
		return nil, &BadRequestError{
			Message: "Invalid cursor format",
			Order:   "H-Cursor-Decode-1",
		}
	}
	
	var cursor CursorData
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return nil, &BadRequestError{
			Message: "Invalid cursor data",
			Order:   "H-Cursor-Decode-2",
		}
	}
	
	return &cursor, nil
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
	
	// Default sort by
	if pagination.SortBy == "" {
		pagination.SortBy = "id"
	}
	
	// Decode cursor
	cursor, err := DecodeCursor(pagination.Cursor)
	if err != nil {
		return db, err
	}
	
	// Apply cursor condition
	if cursor != nil {
		if pagination.SortOrder == "desc" {
			if cursor.LastValue != "" {
				// Sort by custom field
				db = db.Where(pagination.SortBy+" < ? OR ("+pagination.SortBy+" = ? AND id < ?)", 
					cursor.LastValue, cursor.LastValue, cursor.LastID)
			} else {
				// Sort by ID only
				db = db.Where("id < ?", cursor.LastID)
			}
		} else {
			if cursor.LastValue != "" {
				db = db.Where(pagination.SortBy+" > ? OR ("+pagination.SortBy+" = ? AND id > ?)", 
					cursor.LastValue, cursor.LastValue, cursor.LastID)
			} else {
				db = db.Where("id > ?", cursor.LastID)
			}
		}
	}
	
	// Apply sorting and limit
	orderClause := fmt.Sprintf("%s %s, id %s", pagination.SortBy, pagination.SortOrder, pagination.SortOrder)
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
func ParseCursorParams(cursor, limit, sortBy, sortOrder string) CursorPagination {
	pagination := DefaultCursorPagination()
	
	pagination.Cursor = cursor
	
	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			pagination.Limit = l
		}
	}
	
	if sortBy != "" {
		pagination.SortBy = sortBy
	}
	
	if sortOrder == "asc" || sortOrder == "desc" {
		pagination.SortOrder = sortOrder
	}
	
	return pagination
}
