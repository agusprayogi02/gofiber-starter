package utils

import (
	"fmt"
	"starter-gofiber/pkg/apierror"
	"strings"
	"time"

	"gorm.io/gorm"
)

// FilterOperator represents filter comparison operators
type FilterOperator string

const (
	OpEqual              FilterOperator = "eq"       // Equal (=)
	OpNotEqual           FilterOperator = "ne"       // Not Equal (!=)
	OpGreaterThan        FilterOperator = "gt"       // Greater Than (>)
	OpGreaterThanOrEqual FilterOperator = "gte"      // Greater Than or Equal (>=)
	OpLessThan           FilterOperator = "lt"       // Less Than (<)
	OpLessThanOrEqual    FilterOperator = "lte"      // Less Than or Equal (<=)
	OpLike               FilterOperator = "like"     // Like (%value%)
	OpNotLike            FilterOperator = "notlike"  // Not Like
	OpIn                 FilterOperator = "in"       // In (value1, value2, ...)
	OpNotIn              FilterOperator = "notin"    // Not In
	OpBetween            FilterOperator = "between"  // Between value1 AND value2
	OpIsNull             FilterOperator = "isnull"   // IS NULL
	OpIsNotNull          FilterOperator = "notnull"  // IS NOT NULL
	OpStartsWith         FilterOperator = "starts"   // Starts with
	OpEndsWith           FilterOperator = "ends"     // Ends with
	OpContains           FilterOperator = "contains" // Contains (case-insensitive)
)

// Filter represents a single filter condition
type Filter struct {
	Field    string         `json:"field" query:"field"`
	Operator FilterOperator `json:"operator" query:"operator"`
	Value    interface{}    `json:"value" query:"value"`
	Values   []interface{}  `json:"values,omitempty" query:"values"` // For IN, NOT IN, BETWEEN
}

// FilterGroup represents a group of filters with AND/OR logic
type FilterGroup struct {
	Logic   string   `json:"logic" query:"logic"` // AND or OR
	Filters []Filter `json:"filters" query:"filters"`
}

// SearchFilter represents complete search and filter parameters
type SearchFilter struct {
	Search       string       `json:"search" query:"search"`               // General search term
	SearchFields []string     `json:"search_fields" query:"search_fields"` // Fields to search in
	Filters      []Filter     `json:"filters" query:"filters"`
	FilterGroup  *FilterGroup `json:"filter_group,omitempty" query:"filter_group"`
}

// ApplySearch applies general search across multiple fields
func ApplySearch(db *gorm.DB, search string, fields []string) *gorm.DB {
	if search == "" || len(fields) == 0 {
		return db
	}

	// Build search condition
	var conditions []string
	var values []interface{}

	searchPattern := "%" + search + "%"
	for _, field := range fields {
		conditions = append(conditions, fmt.Sprintf("%s LIKE ?", field))
		values = append(values, searchPattern)
	}

	// Join conditions with OR
	query := strings.Join(conditions, " OR ")
	return db.Where(query, values...)
}

// ApplyFilter applies a single filter to query
func ApplyFilter(db *gorm.DB, filter Filter) *gorm.DB {
	if filter.Field == "" {
		return db
	}

	switch filter.Operator {
	case OpEqual:
		return db.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)

	case OpNotEqual:
		return db.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)

	case OpGreaterThan:
		return db.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)

	case OpGreaterThanOrEqual:
		return db.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)

	case OpLessThan:
		return db.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)

	case OpLessThanOrEqual:
		return db.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)

	case OpLike:
		return db.Where(fmt.Sprintf("%s LIKE ?", filter.Field), "%"+fmt.Sprint(filter.Value)+"%")

	case OpNotLike:
		return db.Where(fmt.Sprintf("%s NOT LIKE ?", filter.Field), "%"+fmt.Sprint(filter.Value)+"%")

	case OpIn:
		if len(filter.Values) > 0 {
			return db.Where(fmt.Sprintf("%s IN ?", filter.Field), filter.Values)
		}
		return db

	case OpNotIn:
		if len(filter.Values) > 0 {
			return db.Where(fmt.Sprintf("%s NOT IN ?", filter.Field), filter.Values)
		}
		return db

	case OpBetween:
		if len(filter.Values) >= 2 {
			return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), filter.Values[0], filter.Values[1])
		}
		return db

	case OpIsNull:
		return db.Where(fmt.Sprintf("%s IS NULL", filter.Field))

	case OpIsNotNull:
		return db.Where(fmt.Sprintf("%s IS NOT NULL", filter.Field))

	case OpStartsWith:
		return db.Where(fmt.Sprintf("%s LIKE ?", filter.Field), fmt.Sprint(filter.Value)+"%")

	case OpEndsWith:
		return db.Where(fmt.Sprintf("%s LIKE ?", filter.Field), "%"+fmt.Sprint(filter.Value))

	case OpContains:
		return db.Where(fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", filter.Field), "%"+fmt.Sprint(filter.Value)+"%")

	default:
		return db
	}
}

// ApplyFilters applies multiple filters to query
func ApplyFilters(db *gorm.DB, filters []Filter) *gorm.DB {
	for _, filter := range filters {
		db = ApplyFilter(db, filter)
	}
	return db
}

// ApplyFilterGroup applies filter group with AND/OR logic
func ApplyFilterGroup(db *gorm.DB, group FilterGroup) *gorm.DB {
	if len(group.Filters) == 0 {
		return db
	}

	logic := strings.ToUpper(group.Logic)
	if logic != "AND" && logic != "OR" {
		logic = "AND" // Default to AND
	}

	if logic == "OR" {
		// Build OR conditions
		orDB := db.Session(&gorm.Session{NewDB: true})
		for i, filter := range group.Filters {
			if i == 0 {
				orDB = ApplyFilter(orDB, filter)
			} else {
				orDB = orDB.Or(func(tx *gorm.DB) *gorm.DB {
					return ApplyFilter(tx, filter)
				})
			}
		}
		return db.Where(orDB)
	}

	// AND logic (default)
	return ApplyFilters(db, group.Filters)
}

// ApplySearchFilter applies complete search and filter
func ApplySearchFilter(db *gorm.DB, searchFilter SearchFilter) *gorm.DB {
	// Apply general search
	if searchFilter.Search != "" && len(searchFilter.SearchFields) > 0 {
		db = ApplySearch(db, searchFilter.Search, searchFilter.SearchFields)
	}

	// Apply individual filters
	if len(searchFilter.Filters) > 0 {
		db = ApplyFilters(db, searchFilter.Filters)
	}

	// Apply filter group
	if searchFilter.FilterGroup != nil {
		db = ApplyFilterGroup(db, *searchFilter.FilterGroup)
	}

	return db
}

// DateRangeFilter creates a date range filter
func DateRangeFilter(field string, from, to time.Time) Filter {
	return Filter{
		Field:    field,
		Operator: OpBetween,
		Values:   []interface{}{from, to},
	}
}

// MultiValueFilter creates IN filter
func MultiValueFilter(field string, values []interface{}) Filter {
	return Filter{
		Field:    field,
		Operator: OpIn,
		Values:   values,
	}
}

// TextSearchFilter creates a text search filter
func TextSearchFilter(field, value string) Filter {
	return Filter{
		Field:    field,
		Operator: OpContains,
		Value:    value,
	}
}

// RangeFilter creates a numeric range filter
func RangeFilter(field string, min, max interface{}) Filter {
	return Filter{
		Field:    field,
		Operator: OpBetween,
		Values:   []interface{}{min, max},
	}
}

// ValidateFilter validates filter parameters
func ValidateFilter(filter Filter, allowedFields []string) error {
	// Check if field is allowed
	if len(allowedFields) > 0 {
		found := false
		for _, allowed := range allowedFields {
			if filter.Field == allowed {
				found = true
				break
			}
		}
		if !found {
			return &apierror.BadRequestError{
				Message: fmt.Sprintf("Field '%s' is not allowed for filtering", filter.Field),
				Order:   "H-Filter-Validate-1",
			}
		}
	}

	// Validate operator-specific requirements
	switch filter.Operator {
	case OpIn, OpNotIn, OpBetween:
		if len(filter.Values) == 0 {
			return &apierror.BadRequestError{
				Message: fmt.Sprintf("Operator '%s' requires 'values' array", filter.Operator),
				Order:   "H-Filter-Validate-2",
			}
		}
	case OpIsNull, OpIsNotNull:
		// These don't need values
	default:
		if filter.Value == nil {
			return &apierror.BadRequestError{
				Message: fmt.Sprintf("Operator '%s' requires 'value' parameter", filter.Operator),
				Order:   "H-Filter-Validate-3",
			}
		}
	}

	return nil
}

// BuildFilterFromQuery builds filter from query parameters
// Example: ?filter_name_eq=John&filter_age_gte=18&filter_email_like=gmail
func BuildFilterFromQuery(params map[string]string) []Filter {
	filters := make([]Filter, 0)

	for key, value := range params {
		if !strings.HasPrefix(key, "filter_") {
			continue
		}

		// Remove "filter_" prefix
		parts := strings.Split(key[7:], "_")
		if len(parts) < 2 {
			continue
		}

		// Last part is operator, rest is field name
		operator := FilterOperator(parts[len(parts)-1])
		field := strings.Join(parts[:len(parts)-1], "_")

		filters = append(filters, Filter{
			Field:    field,
			Operator: operator,
			Value:    value,
		})
	}

	return filters
}
