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
pagination := pagination.ParseCursorParams(
    c.Query("cursor"),
    c.Query("limit"),
    c.Query("sort_by"),
    c.Query("sort_order"),
)

// Apply ke query
db, _ := pagination.ApplyCursorPagination(db, pagination)
db.Find(&results)

// Build response
response := pagination.BuildCursorResponse(results, pagination)
```

**Query:** `GET /users?limit=20&cursor=eyJ...`

## Bulk Operations

```go
// Bulk Create
result, _ := database.BulkCreate(db, &items, 100)
result, _ := database.BulkCreateWithValidation(db, &items, validateFunc, 100)

// Bulk Update
result, _ := database.BulkUpdate(db, &Model{}, ids, updates)
result, _ := database.BulkUpdateWithValidation(db, &Model{}, ids, updates, validateFunc)

// Bulk Delete
result, _ := database.BulkDelete(db, &Model{}, ids)              // Soft delete
result, _ := database.BulkDeletePermanent(db, &Model{}, ids)     // Hard delete
result, _ := database.BulkRestore(db, &Model{}, ids)             // Restore

// Bulk Upsert
result, _ := database.BulkUpsert(db, &items, []string{"email"}, []string{"name", "age"})
```

## Export Data

```go
// CSV
filename, _ := utils.ExportToCSV(data, headers, "output.csv")

// Excel
config := utils.ExportConfig{
    Filename:  "output.xlsx",
    SheetName: "Sheet1",
    Headers:   headers,
    Format:    utils.FormatExcel,
}
filename, _ := utils.ExportToExcel(data, config)

// PDF
config.Format = utils.FormatPDF
config.Title = "Report Title"
filename, _ := utils.ExportToPDF(data, config)

// Generic
filename, _ := utils.ExportData(data, headers, config)
```

## Search & Filter

```go
// Simple search
searchFilter := utils.SearchFilter{
    Search:      "keyword",
    SearchFields: []string{"name", "email"},
}
db = utils.ApplySearchFilter(db, searchFilter)

// Single filter
filter := utils.Filter{
    Field:    "age",
    Operator: utils.OpGreaterThanOrEqual,
    Value:    18,
}
db = utils.ApplyFilter(db, filter)

// Multiple filters (AND)
filters := []utils.Filter{
    {Field: "status", Operator: utils.OpEqual, Value: "active"},
    {Field: "age", Operator: utils.OpGreaterThanOrEqual, Value: 18},
}
db = utils.ApplyFilters(db, filters)

// Filter group (OR)
filterGroup := utils.FilterGroup{
    Logic: "OR",
    Filters: []utils.Filter{
        {Field: "role", Operator: utils.OpEqual, Value: "admin"},
        {Field: "role", Operator: utils.OpEqual, Value: "moderator"},
    },
}
db = utils.ApplyFilterGroup(db, filterGroup)

// From query string
filters := utils.BuildFilterFromQuery(params)
db = utils.ApplyFilters(db, filters)
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
db = utils.ApplySort(db, "created_at", utils.SortDesc)

// Multiple fields
sortFields := []utils.SortField{
    {Field: "status", Order: utils.SortAsc},
    {Field: "created_at", Order: utils.SortDesc},
}
db = utils.ApplyMultiSort(db, sortFields, allowedFields)

// From query string
sortConfig := utils.BuildSortFromQuery(params, allowedFields)
db = utils.ApplySortConfig(db, sortConfig)

// Parse sort string
sortFields := utils.ParseSortString("status:asc,created_at:desc")
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
    searchFilter := utils.SearchFilter{
        Search:      c.Query("search"),
        SearchFields: []string{"name", "email"},
        Filters:     utils.BuildFilterFromQuery(params),
    }
    
    // Sorting
    allowedSortFields := []string{"name", "email", "created_at"}
    sortConfig := utils.BuildSortFromQuery(params, allowedSortFields)
    
    // Pagination
    pagination := pagination.ParseCursorParams(
        c.Query("cursor"),
        c.Query("limit"),
        c.Query("sort_by"),
        c.Query("sort_order"),
    )
    
    // Build query
    db := config.DB.Model(&user.User{})
    db = utils.ApplySearchFilter(db, searchFilter)
    db = utils.ApplySortConfig(db, sortConfig)
    db, _ = pagination.ApplyCursorPagination(db, pagination)
    
    // Execute
    var users []user.User
    db.Find(&users)
    
    // Response
    return c.JSON(pagination.BuildCursorResponse(users, pagination))
}
```

**Full Query Example:**
```
GET /users?search=john&filter_status_eq=active&filter_age_gte=18&sort=-created_at&limit=20&cursor=eyJ...
```

## Helper Functions

```go
// Date range
filter := utils.DateRangeFilter("created_at", fromDate, toDate)

// Multi-value (IN)
filter := utils.MultiValueFilter("role", []interface{}{"admin", "mod"})

// Text search
filter := utils.TextSearchFilter("name", "john")

// Numeric range
filter := utils.RangeFilter("age", 18, 65)

// Validate filter
err := utils.ValidateFilter(filter, allowedFields)

// Validate sort
err := utils.ValidateSortFields(sortFields, allowedFields)
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
