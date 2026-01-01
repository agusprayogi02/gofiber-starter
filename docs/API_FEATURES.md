# API Features Documentation

Dokumentasi lengkap untuk fitur-fitur API advanced seperti compression, pagination, bulk operations, export data, search & filter, dan sorting.

## Table of Contents

- [Response Compression](#response-compression)
- [Cursor-based Pagination](#cursor-based-pagination)
- [Bulk Operations](#bulk-operations)
- [Export Data](#export-data)
- [Search & Filter](#search--filter)
- [Sorting Support](#sorting-support)
- [Combined Usage](#combined-usage)

---

## Response Compression

Middleware untuk mengompress response dengan Gzip untuk mengurangi bandwidth dan mempercepat transfer data.

### Setup

```go
import "starter-gofiber/internal/handler/middleware"

// Default compression (balanced)
app.Use(middleware.CompressionDefault())

// Or with custom level
app.Use(middleware.CompressionBestSpeed())    // Fastest, larger size
app.Use(middleware.CompressionBestSize())     // Slowest, smallest size
```

### Features

- Gzip compression untuk semua response
- 3 level kompresi: Default, Best Speed, Best Size
- Otomatis skip untuk response yang sudah compressed
- Skip untuk file kecil (<200 bytes)

### Performance Impact

| Level | Compression Ratio | Speed | Use Case |
|-------|------------------|-------|----------|
| Best Speed | ~50-60% | Fastest | High traffic APIs |
| Default | ~60-70% | Balanced | General purpose |
| Best Size | ~70-80% | Slowest | Large payloads |

---

## Cursor-based Pagination

Alternative untuk offset pagination yang lebih efisien untuk dataset besar.

### Basic Usage

```go
import "starter-gofiber/pkg/apierror"

func GetPosts(c *fiber.Ctx) error {
    // Parse cursor parameters dari query string
    pagination := pagination.ParseCursorParams(
        c.Query("cursor"),      // Next cursor dari response sebelumnya
        c.Query("limit"),       // Jumlah data per page (default: 10, max: 100)
        c.Query("sort_by"),     // Field untuk sorting (default: "id")
        c.Query("sort_order"),  // asc atau desc (default: "asc")
    )
    
    // Apply cursor pagination ke GORM query
    db := config.DB
    var posts []post.Post
    
    db, err := pagination.ApplyCursorPagination(db, pagination)
    if err != nil {
        return err
    }
    
    db.Find(&posts)
    
    // Build response dengan next cursor
    response := pagination.BuildCursorResponse(posts, pagination)
    
    return c.JSON(response)
}
```

### Query Examples

```bash
# First page (10 items)
GET /posts?limit=10

# Next page using cursor
GET /posts?limit=10&cursor=eyJsYXN0X2lkIjoxMCwibGFzdF92YWx1ZSI6IjIwMjQtMDEtMDEifQ==

# Custom sorting
GET /posts?sort_by=created_at&sort_order=desc&limit=20
```

### Response Format

```json
{
  "data": [...],
  "next_cursor": "eyJsYXN0X2lkIjoyMCwibGFzdF92YWx1ZSI6IjIwMjQtMDEtMDIifQ==",
  "has_more": true,
  "count": 10
}
```

### Advantages vs Offset Pagination

| Feature | Cursor | Offset |
|---------|--------|--------|
| Performance | O(1) | O(n) |
| Consistency | ✅ No missing/duplicate | ❌ Can miss/duplicate |
| Large datasets | ✅ Efficient | ❌ Slow |
| Random access | ❌ Sequential only | ✅ Jump to page |

---

## Bulk Operations

Operasi create, update, dan delete dalam jumlah banyak dengan error tracking.

### Bulk Create

```go
import "starter-gofiber/pkg/apierror"

func BulkCreateUsers(c *fiber.Ctx) error {
    var users []user.User
    if err := c.BodyParser(&users); err != nil {
        return err
    }
    
    // Simple bulk create (all or nothing)
    result, err := database.BulkCreate(config.DB, &users, 100) // batch size 100
    if err != nil {
        return err
    }
    
    // Or with individual validation
    result, err := database.BulkCreateWithValidation(
        config.DB,
        &users,
        func(user *user.User) error {
            // Custom validation per item
            if user.Email == "" {
                return errors.New("email required")
            }
            return nil
        },
        100,
    )
    
    return c.JSON(fiber.Map{
        "success": result.Success,
        "failed": result.Failed,
        "errors": result.Errors,
    })
}
```

### Bulk Update

```go
func BulkUpdateUsers(c *fiber.Ctx) error {
    var req struct {
        IDs     []uint                 `json:"ids"`
        Updates map[string]interface{} `json:"updates"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return err
    }
    
    // Update multiple records dengan field yang sama
    result, err := database.BulkUpdate(
        config.DB,
        &user.User{},
        req.IDs,
        req.Updates,
    )
    
    // Or with individual validation
    result, err := database.BulkUpdateWithValidation(
        config.DB,
        &user.User{},
        req.IDs,
        req.Updates,
        func(user *user.User) error {
            // Validation after update
            return nil
        },
    )
    
    return c.JSON(result)
}
```

### Bulk Delete

```go
func BulkDeleteUsers(c *fiber.Ctx) error {
    var req struct {
        IDs []uint `json:"ids"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return err
    }
    
    // Soft delete
    result, err := database.BulkDelete(config.DB, &user.User{}, req.IDs)
    
    // Hard delete (permanent)
    result, err := database.BulkDeletePermanent(config.DB, &user.User{}, req.IDs)
    
    return c.JSON(result)
}
```

### Bulk Restore

```go
func BulkRestoreUsers(c *fiber.Ctx) error {
    var req struct {
        IDs []uint `json:"ids"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return err
    }
    
    result, err := database.BulkRestore(config.DB, &user.User{}, req.IDs)
    
    return c.JSON(result)
}
```

### Bulk Upsert

```go
func BulkUpsertUsers(c *fiber.Ctx) error {
    var users []user.User
    if err := c.BodyParser(&users); err != nil {
        return err
    }
    
    // Insert or update berdasarkan unique key (email)
    result, err := database.BulkUpsert(
        config.DB,
        &users,
        []string{"email"}, // Conflict columns
        []string{"name", "age", "updated_at"}, // Columns to update
    )
    
    return c.JSON(result)
}
```

### Error Handling

```json
{
  "success": 95,
  "failed": 5,
  "errors": [
    {
      "index": 3,
      "id": 0,
      "error": "email already exists"
    },
    {
      "index": 7,
      "id": 0,
      "error": "invalid phone number"
    }
  ]
}
```

---

## Export Data

Export data ke format CSV, Excel, atau PDF.

### CSV Export

```go
import "starter-gofiber/pkg/apierror"

func ExportUsersCSV(c *fiber.Ctx) error {
    var users []user.User
    config.DB.Find(&users)
    
    headers := []string{"ID", "Name", "Email", "Created At"}
    
    filename, err := utils.ExportToCSV(users, headers, "users.csv")
    if err != nil {
        return err
    }
    
    return c.SendFile(filename)
}
```

### Excel Export

```go
func ExportUsersExcel(c *fiber.Ctx) error {
    var users []user.User
    config.DB.Find(&users)
    
    headers := []string{"ID", "Name", "Email", "Phone", "Created At"}
    
    // With custom config
    config := utils.ExportConfig{
        Filename:  "users.xlsx",
        SheetName: "Users Data",
        Headers:   headers,
        Format:    utils.FormatExcel,
    }
    
    filename, err := utils.ExportToExcel(users, config)
    if err != nil {
        return err
    }
    
    return c.SendFile(filename)
}
```

### PDF Export

```go
func ExportUsersPDF(c *fiber.Ctx) error {
    var users []user.User
    config.DB.Find(&users)
    
    headers := []string{"ID", "Name", "Email", "Phone", "Created At"}
    
    config := utils.ExportConfig{
        Filename: "users.pdf",
        Title:    "Users Report",
        Headers:  headers,
        Format:   utils.FormatPDF,
    }
    
    filename, err := utils.ExportToPDF(users, config)
    if err != nil {
        return err
    }
    
    return c.SendFile(filename)
}
```

### Generic Export

```go
func ExportUsers(c *fiber.Ctx) error {
    format := c.Query("format", "csv") // csv, excel, or pdf
    
    var users []user.User
    config.DB.Find(&users)
    
    headers := []string{"ID", "Name", "Email", "Phone", "Created At"}
    
    exportFormat := utils.FormatCSV
    switch format {
    case "excel":
        exportFormat = utils.FormatExcel
    case "pdf":
        exportFormat = utils.FormatPDF
    }
    
    config := utils.DefaultExportConfig(exportFormat)
    config.Filename = "users_" + time.Now().Format("20060102_150405")
    config.Headers = headers
    config.Title = "Users Report"
    
    filename, err := utils.ExportData(users, headers, config)
    if err != nil {
        return err
    }
    
    return c.SendFile(filename)
}
```

### Features

**CSV:**
- Simple comma-separated format
- UTF-8 encoding
- Fast generation

**Excel:**
- Bold colored headers (#4472C4 background, white text)
- Auto-fit columns
- Custom sheet names
- Professional styling

**PDF:**
- Colored table headers
- Automatic pagination (35 rows per page)
- Repeated headers on each page
- Page numbers in footer
- Professional layout

### Type Support

Export helper otomatis convert semua Go types:
- int, int8, int16, int32, int64
- uint, uint8, uint16, uint32, uint64
- float32, float64
- bool
- string
- time.Time (formatted as "2006-01-02 15:04:05")
- Pointers (handled safely)
- Nested structs (converted to string)

---

## Search & Filter

Advanced search dan filtering dengan multiple operators.

### Basic Search

```go
import "starter-gofiber/pkg/apierror"

func SearchUsers(c *fiber.Ctx) error {
    searchFilter := utils.SearchFilter{
        Search:      c.Query("search"),
        SearchFields: []string{"name", "email", "phone"},
    }
    
    db := utils.ApplySearchFilter(config.DB, searchFilter)
    
    var users []user.User
    db.Find(&users)
    
    return c.JSON(users)
}
```

Query: `GET /users?search=john`

### Single Filter

```go
func FilterUsers(c *fiber.Ctx) error {
    filter := utils.Filter{
        Field:    "age",
        Operator: utils.OpGreaterThanOrEqual,
        Value:    18,
    }
    
    db := utils.ApplyFilter(config.DB, filter)
    
    var users []user.User
    db.Find(&users)
    
    return c.JSON(users)
}
```

### Multiple Filters (AND)

```go
func FilterUsersMultiple(c *fiber.Ctx) error {
    filters := []utils.Filter{
        {
            Field:    "status",
            Operator: utils.OpEqual,
            Value:    "active",
        },
        {
            Field:    "age",
            Operator: utils.OpGreaterThanOrEqual,
            Value:    18,
        },
        {
            Field:    "email",
            Operator: utils.OpLike,
            Value:    "gmail",
        },
    }
    
    db := utils.ApplyFilters(config.DB, filters)
    
    var users []user.User
    db.Find(&users)
    
    return c.JSON(users)
}
```

### Filter Group (OR Logic)

```go
func FilterUsersOR(c *fiber.Ctx) error {
    filterGroup := utils.FilterGroup{
        Logic: "OR",
        Filters: []utils.Filter{
            {
                Field:    "role",
                Operator: utils.OpEqual,
                Value:    "admin",
            },
            {
                Field:    "role",
                Operator: utils.OpEqual,
                Value:    "moderator",
            },
        },
    }
    
    db := utils.ApplyFilterGroup(config.DB, filterGroup)
    
    var users []user.User
    db.Find(&users)
    
    return c.JSON(users)
}
```

### From Query String

```go
func FilterFromQuery(c *fiber.Ctx) error {
    // Parse: ?filter_age_gte=18&filter_status_eq=active&filter_email_like=gmail
    params := make(map[string]string)
    c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
        params[string(key)] = string(value)
    })
    
    filters := utils.BuildFilterFromQuery(params)
    
    db := utils.ApplyFilters(config.DB, filters)
    
    var users []user.User
    db.Find(&users)
    
    return c.JSON(users)
}
```

Query examples:
```bash
# Equal
GET /users?filter_status_eq=active

# Greater than or equal
GET /users?filter_age_gte=18

# Like (contains)
GET /users?filter_email_like=gmail

# Multiple filters
GET /users?filter_status_eq=active&filter_age_gte=18&filter_role_in=admin,moderator
```

### Available Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `eq` | Equal (=) | `status = 'active'` |
| `ne` | Not Equal (!=) | `status != 'banned'` |
| `gt` | Greater Than (>) | `age > 18` |
| `gte` | Greater Than or Equal (>=) | `age >= 18` |
| `lt` | Less Than (<) | `age < 65` |
| `lte` | Less Than or Equal (<=) | `age <= 65` |
| `like` | Like (%value%) | `name LIKE '%john%'` |
| `notlike` | Not Like | `email NOT LIKE '%temp%'` |
| `in` | In (value1, value2, ...) | `role IN ('admin', 'mod')` |
| `notin` | Not In | `status NOT IN ('banned', 'deleted')` |
| `between` | Between value1 AND value2 | `age BETWEEN 18 AND 65` |
| `isnull` | IS NULL | `deleted_at IS NULL` |
| `notnull` | IS NOT NULL | `email IS NOT NULL` |
| `starts` | Starts with | `name LIKE 'john%'` |
| `ends` | Ends with | `email LIKE '%@gmail.com'` |
| `contains` | Contains (case-insensitive) | `LOWER(name) LIKE '%john%'` |

### Helper Functions

```go
// Date range filter
filter := utils.DateRangeFilter(
    "created_at",
    time.Now().AddDate(0, -1, 0), // From 1 month ago
    time.Now(),                     // To now
)

// Multi-value filter (IN)
filter := utils.MultiValueFilter(
    "role",
    []interface{}{"admin", "moderator", "editor"},
)

// Text search filter
filter := utils.TextSearchFilter("name", "john")

// Numeric range filter
filter := utils.RangeFilter("age", 18, 65)
```

### Field Validation

```go
func FilterWithValidation(c *fiber.Ctx) error {
    allowedFields := []string{"name", "email", "age", "status", "role"}
    
    filter := utils.Filter{
        Field:    c.Query("field"),
        Operator: utils.FilterOperator(c.Query("op")),
        Value:    c.Query("value"),
    }
    
    // Validate field is allowed
    if err := utils.ValidateFilter(filter, allowedFields); err != nil {
        return err
    }
    
    db := utils.ApplyFilter(config.DB, filter)
    
    var users []user.User
    db.Find(&users)
    
    return c.JSON(users)
}
```

---

## Sorting Support

Multi-field sorting dengan validation.

### Basic Sorting

```go
import "starter-gofiber/pkg/apierror"

func GetUsers(c *fiber.Ctx) error {
    // Simple sorting
    db := utils.ApplySort(
        config.DB,
        "created_at",
        utils.SortDesc,
    )
    
    var users []user.User
    db.Find(&users)
    
    return c.JSON(users)
}
```

### Multi-field Sorting

```go
func GetUsersMultiSort(c *fiber.Ctx) error {
    sortFields := []utils.SortField{
        {Field: "status", Order: utils.SortAsc},
        {Field: "created_at", Order: utils.SortDesc},
    }
    
    allowedFields := []string{"name", "email", "status", "created_at"}
    
    db := utils.ApplyMultiSort(config.DB, sortFields, allowedFields)
    
    var users []user.User
    db.Find(&users)
    
    return c.JSON(users)
}
```

### From Query String

```go
func GetUsersWithSort(c *fiber.Ctx) error {
    // Parse dari query params
    params := make(map[string]string)
    c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
        params[string(key)] = string(value)
    })
    
    allowedFields := []string{"name", "email", "age", "created_at", "updated_at"}
    sortConfig := utils.BuildSortFromQuery(params, allowedFields)
    
    db := utils.ApplySortConfig(config.DB, sortConfig)
    
    var users []user.User
    db.Find(&users)
    
    return c.JSON(users)
}
```

### Query String Formats

```bash
# Single field ascending
GET /users?sort=name

# Single field descending (- prefix)
GET /users?sort=-created_at

# Multiple fields with colon notation
GET /users?sort=status:asc,created_at:desc

# Separate parameters
GET /users?sort_by=name&order=asc

# Multiple sorts
GET /users?sort=status:asc,age:desc,name:asc
```

### Parse Sort String

```go
// From string: "name:asc,created_at:desc"
sortFields := utils.ParseSortString(c.Query("sort"))

// From separate params
sortField := utils.ParseSortParams(
    c.Query("sort_by"),
    c.Query("order"),
)
```

### With Default

```go
func GetUsers(c *fiber.Ctx) error {
    // Default sort if not provided
    sortConfig := utils.DefaultSortConfig(
        "created_at",           // Default field
        utils.SortDesc,        // Default order
        []string{"name", "email", "created_at"}, // Allowed fields
    )
    
    // Override with query params
    params := make(map[string]string)
    c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
        params[string(key)] = string(value)
    })
    
    if sort, ok := params["sort"]; ok {
        sortConfig.Fields = utils.ParseSortString(sort)
    }
    
    db := utils.ApplySortConfig(config.DB, sortConfig)
    
    var users []user.User
    db.Find(&users)
    
    return c.JSON(users)
}
```

### Field Validation

```go
func GetUsersValidated(c *fiber.Ctx) error {
    allowedFields := []string{"name", "email", "created_at"}
    
    sortFields := utils.ParseSortString(c.Query("sort"))
    
    // Validate all sort fields
    if err := utils.ValidateSortFields(sortFields, allowedFields); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    db := utils.ApplyMultiSort(config.DB, sortFields, allowedFields)
    
    var users []user.User
    db.Find(&users)
    
    return c.JSON(users)
}
```

---

## Combined Usage

Contoh penggunaan semua fitur secara bersamaan.

### Advanced List Endpoint

```go
func GetUsersAdvanced(c *fiber.Ctx) error {
    // 1. Parse search & filter
    searchFilter := utils.SearchFilter{
        Search:      c.Query("search"),
        SearchFields: []string{"name", "email"},
    }
    
    // Parse filters dari query
    params := make(map[string]string)
    c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
        params[string(key)] = string(value)
    })
    
    searchFilter.Filters = utils.BuildFilterFromQuery(params)
    
    // 2. Parse sorting
    allowedSortFields := []string{"name", "email", "created_at", "updated_at"}
    sortConfig := utils.BuildSortFromQuery(params, allowedSortFields)
    
    // Set default sort jika tidak ada
    if len(sortConfig.Fields) == 0 {
        sortConfig.DefaultField = "created_at"
        sortConfig.DefaultOrder = utils.SortDesc
    }
    
    // 3. Parse cursor pagination
    pagination := pagination.ParseCursorParams(
        c.Query("cursor"),
        c.Query("limit"),
        c.Query("sort_by"),
        c.Query("sort_order"),
    )
    
    // 4. Build query
    db := config.DB.Model(&user.User{})
    
    // Apply search & filter
    db = utils.ApplySearchFilter(db, searchFilter)
    
    // Apply sorting
    db = utils.ApplySortConfig(db, sortConfig)
    
    // Apply pagination
    db, err := pagination.ApplyCursorPagination(db, pagination)
    if err != nil {
        return err
    }
    
    // 5. Execute query
    var users []user.User
    db.Find(&users)
    
    // 6. Build response
    response := pagination.BuildCursorResponse(users, pagination)
    
    return c.JSON(response)
}
```

### Query Example

```bash
GET /users?search=john&filter_status_eq=active&filter_age_gte=18&sort=-created_at&limit=20&cursor=eyJ...
```

This will:
1. Search "john" in name and email
2. Filter status = "active" AND age >= 18
3. Sort by created_at descending
4. Return 20 items with cursor pagination

### With Export

```go
func ExportUsersAdvanced(c *fiber.Ctx) error {
    format := c.Query("format", "excel")
    
    // Same query building as above
    searchFilter := utils.SearchFilter{...}
    sortConfig := utils.BuildSortFromQuery(params, allowedFields)
    
    db := config.DB.Model(&user.User{})
    db = utils.ApplySearchFilter(db, searchFilter)
    db = utils.ApplySortConfig(db, sortConfig)
    
    // Get all data (no pagination for export)
    var users []user.User
    db.Find(&users)
    
    // Export
    headers := []string{"ID", "Name", "Email", "Status", "Created At"}
    
    var exportFormat utils.ExportFormat
    switch format {
    case "csv":
        exportFormat = utils.FormatCSV
    case "pdf":
        exportFormat = utils.FormatPDF
    default:
        exportFormat = utils.FormatExcel
    }
    
    config := utils.DefaultExportConfig(exportFormat)
    config.Headers = headers
    config.Title = "Users Report - " + time.Now().Format("2006-01-02")
    
    filename, err := utils.ExportData(users, headers, config)
    if err != nil {
        return err
    }
    
    return c.SendFile(filename)
}
```

### Complete Example Handler

```go
package handler

import (
    "starter-gofiber/internal/config"
    "starter-gofiber/internal/domain/user"
    "starter-gofiber/pkg/apierror"
    "github.com/gofiber/fiber/v2"
    "time"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
    return &UserHandler{}
}

// List with search, filter, sort, and pagination
func (h *UserHandler) List(c *fiber.Ctx) error {
    // Parse all parameters
    params := make(map[string]string)
    c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
        params[string(key)] = string(value)
    })
    
    // Search & Filter
    searchFilter := utils.SearchFilter{
        Search:      c.Query("search"),
        SearchFields: []string{"name", "email", "phone"},
        Filters:     utils.BuildFilterFromQuery(params),
    }
    
    // Sorting
    allowedSortFields := []string{"id", "name", "email", "created_at", "updated_at"}
    sortConfig := utils.BuildSortFromQuery(params, allowedSortFields)
    if len(sortConfig.Fields) == 0 {
        sortConfig.DefaultField = "created_at"
        sortConfig.DefaultOrder = utils.SortDesc
    }
    
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
    response := pagination.BuildCursorResponse(users, pagination)
    return c.JSON(response)
}

// Export to CSV/Excel/PDF
func (h *UserHandler) Export(c *fiber.Ctx) error {
    format := c.Query("format", "excel")
    
    // Same filter logic
    params := make(map[string]string)
    c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
        params[string(key)] = string(value)
    })
    
    searchFilter := utils.SearchFilter{
        Search:      c.Query("search"),
        SearchFields: []string{"name", "email"},
        Filters:     utils.BuildFilterFromQuery(params),
    }
    
    db := config.DB.Model(&user.User{})
    db = utils.ApplySearchFilter(db, searchFilter)
    
    var users []user.User
    db.Find(&users)
    
    // Export
    headers := []string{"ID", "Name", "Email", "Phone", "Created At"}
    
    var exportFormat utils.ExportFormat
    switch format {
    case "csv":
        exportFormat = utils.FormatCSV
    case "pdf":
        exportFormat = utils.FormatPDF
    default:
        exportFormat = utils.FormatExcel
    }
    
    exportConfig := utils.DefaultExportConfig(exportFormat)
    exportConfig.Headers = headers
    exportConfig.Title = "Users Report"
    
    filename, err := utils.ExportData(users, headers, exportConfig)
    if err != nil {
        return err
    }
    
    return c.SendFile(filename)
}

// Bulk create
func (h *UserHandler) BulkCreate(c *fiber.Ctx) error {
    var users []user.User
    if err := c.BodyParser(&users); err != nil {
        return err
    }
    
    result, err := database.BulkCreateWithValidation(
        config.DB,
        &users,
        func(user *user.User) error {
            // Validation
            if user.Email == "" {
                return errors.New("email required")
            }
            return nil
        },
        100,
    )
    
    if err != nil {
        return err
    }
    
    return c.JSON(result)
}

// Bulk update
func (h *UserHandler) BulkUpdate(c *fiber.Ctx) error {
    var req struct {
        IDs     []uint                 `json:"ids"`
        Updates map[string]interface{} `json:"updates"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return err
    }
    
    result, err := database.BulkUpdate(
        config.DB,
        &user.User{},
        req.IDs,
        req.Updates,
    )
    
    if err != nil {
        return err
    }
    
    return c.JSON(result)
}

// Bulk delete
func (h *UserHandler) BulkDelete(c *fiber.Ctx) error {
    var req struct {
        IDs []uint `json:"ids"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return err
    }
    
    result, err := database.BulkDelete(config.DB, &user.User{}, req.IDs)
    if err != nil {
        return err
    }
    
    return c.JSON(result)
}
```

---

## Best Practices

### 1. Compression
- Enable untuk production APIs
- Use Best Speed untuk high traffic
- Skip untuk file uploads/downloads

### 2. Pagination
- Prefer cursor over offset untuk large datasets
- Set reasonable limit (max 100)
- Always provide next_cursor in response

### 3. Bulk Operations
- Use validation version untuk user input
- Set appropriate batch size (100-1000)
- Handle partial failures gracefully
- Log bulk errors for debugging

### 4. Export
- Add export format to query (?format=excel)
- Set timeout untuk large exports
- Consider background jobs untuk exports >10k rows
- Add file cleanup job

### 5. Search & Filter
- Whitelist allowed filter fields
- Validate operator usage
- Use indexes on filtered fields
- Combine with pagination

### 6. Sorting
- Whitelist allowed sort fields
- Validate sort order
- Add indexes on sorted fields
- Set default sort untuk consistency

### 7. Combined Features
- Apply in order: Filter → Sort → Paginate
- Cache frequently used queries
- Monitor query performance
- Add request logging

---

## Performance Tips

### Database Indexes

```sql
-- For sorting and filtering
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_email ON users(email);

-- For search
CREATE INDEX idx_users_name ON users(name);
CREATE FULLTEXT INDEX idx_users_search ON users(name, email);

-- Composite indexes untuk multi-field sorts
CREATE INDEX idx_users_status_created ON users(status, created_at);
```

### Caching

```go
// Cache search results
cacheKey := fmt.Sprintf("users:search:%s:filter:%v:page:%s", 
    searchFilter.Search,
    searchFilter.Filters,
    pagination.Cursor,
)

// Check cache
if cached, err := cache.Get(cacheKey); err == nil {
    return c.JSON(cached)
}

// Get from DB and cache
result := pagination.BuildCursorResponse(users, pagination)
cache.Set(cacheKey, result, 5*time.Minute)
```

### Query Optimization

```go
// Select only needed fields
db.Select("id", "name", "email", "created_at").Find(&users)

// Preload relations efficiently
db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
    return db.Select("id", "title", "user_id")
}).Find(&users)

// Use count cache
var total int64
db.Model(&user.User{}).Count(&total) // Cached count
```

---

## Error Handling

```go
func HandleErrors(c *fiber.Ctx) error {
    // Validate filter
    filter := utils.Filter{...}
    if err := utils.ValidateFilter(filter, allowedFields); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid filter",
            "detail": err.Error(),
        })
    }
    
    // Validate sort
    sortFields := utils.ParseSortString(c.Query("sort"))
    if err := utils.ValidateSortFields(sortFields, allowedFields); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid sort field",
            "detail": err.Error(),
        })
    }
    
    // Validate cursor
    pagination := pagination.ParseCursorParams(...)
    if pagination.Limit > 100 {
        return c.Status(400).JSON(fiber.Map{
            "error": "Limit exceeded",
            "detail": "Maximum limit is 100",
        })
    }
    
    return nil
}
```

---

## Testing

```go
func TestSearchFilter(t *testing.T) {
    // Setup test DB
    db := setupTestDB()
    
    // Create test data
    users := []user.User{
        {Name: "John Doe", Email: "john@example.com", Age: 25},
        {Name: "Jane Smith", Email: "jane@example.com", Age: 30},
    }
    db.Create(&users)
    
    // Test search
    searchFilter := utils.SearchFilter{
        Search:      "john",
        SearchFields: []string{"name", "email"},
    }
    
    query := utils.ApplySearchFilter(db, searchFilter)
    
    var results []user.User
    query.Find(&results)
    
    assert.Equal(t, 1, len(results))
    assert.Equal(t, "John Doe", results[0].Name)
}
```

---

## Summary

Fitur-fitur API ini menyediakan:

✅ **Response Compression** - Reduce bandwidth 50-80%  
✅ **Cursor Pagination** - O(1) performance untuk large datasets  
✅ **Bulk Operations** - Batch processing dengan error tracking  
✅ **Export Data** - CSV, Excel (styled), PDF (paginated)  
✅ **Search & Filter** - 15+ operators, AND/OR logic  
✅ **Sorting** - Multi-field, validated, SQL-injection safe  

Combine semua fitur untuk API endpoints yang powerful dan production-ready!
