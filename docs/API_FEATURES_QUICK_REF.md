# API Features - Quick Reference

Cheat sheet untuk semua API features yang tersedia.

## Response Compression

```go
// Setup
app.Use(middleware.CompressionDefault())
app.Use(middleware.CompressionBestSpeed())  // High traffic
app.Use(middleware.CompressionBestSize())   // Large payloads
```

## Cursor Pagination

```go
// Parse dari query
pagination := helper.ParseCursorParams(
    c.Query("cursor"),
    c.Query("limit"),
    c.Query("sort_by"),
    c.Query("sort_order"),
)

// Apply ke query
db, _ := helper.ApplyCursorPagination(db, pagination)
db.Find(&results)

// Build response
response := helper.BuildCursorResponse(results, pagination)
```

**Query:** `GET /users?limit=20&cursor=eyJ...`

## Bulk Operations

```go
// Bulk Create
result, _ := helper.BulkCreate(db, &items, 100)
result, _ := helper.BulkCreateWithValidation(db, &items, validateFunc, 100)

// Bulk Update
result, _ := helper.BulkUpdate(db, &Model{}, ids, updates)
result, _ := helper.BulkUpdateWithValidation(db, &Model{}, ids, updates, validateFunc)

// Bulk Delete
result, _ := helper.BulkDelete(db, &Model{}, ids)              // Soft delete
result, _ := helper.BulkDeletePermanent(db, &Model{}, ids)     // Hard delete
result, _ := helper.BulkRestore(db, &Model{}, ids)             // Restore

// Bulk Upsert
result, _ := helper.BulkUpsert(db, &items, []string{"email"}, []string{"name", "age"})
```

## Export Data

```go
// CSV
filename, _ := helper.ExportToCSV(data, headers, "output.csv")

// Excel
config := helper.ExportConfig{
    Filename:  "output.xlsx",
    SheetName: "Sheet1",
    Headers:   headers,
    Format:    helper.FormatExcel,
}
filename, _ := helper.ExportToExcel(data, config)

// PDF
config.Format = helper.FormatPDF
config.Title = "Report Title"
filename, _ := helper.ExportToPDF(data, config)

// Generic
filename, _ := helper.ExportData(data, headers, config)
```

## Search & Filter

```go
// Simple search
searchFilter := helper.SearchFilter{
    Search:      "keyword",
    SearchFields: []string{"name", "email"},
}
db = helper.ApplySearchFilter(db, searchFilter)

// Single filter
filter := helper.Filter{
    Field:    "age",
    Operator: helper.OpGreaterThanOrEqual,
    Value:    18,
}
db = helper.ApplyFilter(db, filter)

// Multiple filters (AND)
filters := []helper.Filter{
    {Field: "status", Operator: helper.OpEqual, Value: "active"},
    {Field: "age", Operator: helper.OpGreaterThanOrEqual, Value: 18},
}
db = helper.ApplyFilters(db, filters)

// Filter group (OR)
filterGroup := helper.FilterGroup{
    Logic: "OR",
    Filters: []helper.Filter{
        {Field: "role", Operator: helper.OpEqual, Value: "admin"},
        {Field: "role", Operator: helper.OpEqual, Value: "moderator"},
    },
}
db = helper.ApplyFilterGroup(db, filterGroup)

// From query string
filters := helper.BuildFilterFromQuery(params)
db = helper.ApplyFilters(db, filters)
```

**Operators:**
- `eq`, `ne`, `gt`, `gte`, `lt`, `lte`
- `like`, `notlike`, `in`, `notin`
- `between`, `isnull`, `notnull`
- `starts`, `ends`, `contains`

**Query:** `GET /users?filter_status_eq=active&filter_age_gte=18&filter_role_in=admin,moderator`

## Sorting

```go
// Single field
db = helper.ApplySort(db, "created_at", helper.SortDesc)

// Multiple fields
sortFields := []helper.SortField{
    {Field: "status", Order: helper.SortAsc},
    {Field: "created_at", Order: helper.SortDesc},
}
db = helper.ApplyMultiSort(db, sortFields, allowedFields)

// From query string
sortConfig := helper.BuildSortFromQuery(params, allowedFields)
db = helper.ApplySortConfig(db, sortConfig)

// Parse sort string
sortFields := helper.ParseSortString("status:asc,created_at:desc")
```

**Query formats:**
- `?sort=name` - Ascending
- `?sort=-created_at` - Descending (- prefix)
- `?sort=status:asc,created_at:desc` - Multiple fields
- `?sort_by=name&order=desc` - Separate params

## Combined Usage

```go
func GetUsers(c *fiber.Ctx) error {
    // Parse params
    params := make(map[string]string)
    c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
        params[string(key)] = string(value)
    })
    
    // Search & Filter
    searchFilter := helper.SearchFilter{
        Search:      c.Query("search"),
        SearchFields: []string{"name", "email"},
        Filters:     helper.BuildFilterFromQuery(params),
    }
    
    // Sorting
    allowedSortFields := []string{"name", "email", "created_at"}
    sortConfig := helper.BuildSortFromQuery(params, allowedSortFields)
    
    // Pagination
    pagination := helper.ParseCursorParams(
        c.Query("cursor"),
        c.Query("limit"),
        c.Query("sort_by"),
        c.Query("sort_order"),
    )
    
    // Build query
    db := config.DB.Model(&entity.User{})
    db = helper.ApplySearchFilter(db, searchFilter)
    db = helper.ApplySortConfig(db, sortConfig)
    db, _ = helper.ApplyCursorPagination(db, pagination)
    
    // Execute
    var users []entity.User
    db.Find(&users)
    
    // Response
    return c.JSON(helper.BuildCursorResponse(users, pagination))
}
```

**Full Query Example:**
```
GET /users?search=john&filter_status_eq=active&filter_age_gte=18&sort=-created_at&limit=20&cursor=eyJ...
```

## Helper Functions

```go
// Date range
filter := helper.DateRangeFilter("created_at", fromDate, toDate)

// Multi-value (IN)
filter := helper.MultiValueFilter("role", []interface{}{"admin", "mod"})

// Text search
filter := helper.TextSearchFilter("name", "john")

// Numeric range
filter := helper.RangeFilter("age", 18, 65)

// Validate filter
err := helper.ValidateFilter(filter, allowedFields)

// Validate sort
err := helper.ValidateSortFields(sortFields, allowedFields)
```

## Response Formats

### Cursor Pagination Response
```json
{
  "data": [...],
  "next_cursor": "eyJsYXN0X2lkIjoyMCwibGFzdF92YWx1ZSI6IjIwMjQtMDEtMDIifQ==",
  "has_more": true,
  "count": 10
}
```

### Bulk Operation Response
```json
{
  "success": 95,
  "failed": 5,
  "errors": [
    {
      "index": 3,
      "id": 0,
      "error": "email already exists"
    }
  ]
}
```

## Best Practices

1. **Compression:** Enable untuk production
2. **Pagination:** Set max limit (100), always provide cursor
3. **Bulk Ops:** Use validation version, batch size 100-1000
4. **Export:** Add timeout, consider background jobs untuk >10k rows
5. **Search:** Whitelist fields, use indexes
6. **Sorting:** Whitelist fields, validate order
7. **Combined:** Filter → Sort → Paginate order

## Performance

```sql
-- Add indexes
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_status_created ON users(status, created_at);
CREATE FULLTEXT INDEX idx_users_search ON users(name, email);
```

```go
// Cache results
cacheKey := fmt.Sprintf("users:search:%s:page:%s", search, cursor)
if cached, err := cache.Get(cacheKey); err == nil {
    return c.JSON(cached)
}

// Select only needed fields
db.Select("id", "name", "email").Find(&users)

// Preload efficiently
db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
    return db.Select("id", "title")
}).Find(&users)
```

## Dependencies

```
github.com/gofiber/fiber/v2
gorm.io/gorm
github.com/xuri/excelize/v2 v2.10.0
github.com/jung-kurt/gofpdf v1.16.2
```

---

**Untuk dokumentasi lengkap, lihat [docs/API_FEATURES.md](../docs/API_FEATURES.md)**
