# Architecture Documentation

Dokumentasi arsitektur dan design pattern yang digunakan dalam Starter Template Go Fiber.

## 📐 Architecture Overview

Project ini menggunakan **Clean Architecture** pattern dengan **layered architecture** untuk memisahkan concerns dan meningkatkan maintainability.

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP Layer                            │
│                     (Fiber Framework)                        │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    Handler Layer                             │
│              (Request/Response Handling)                     │
│         - Validate Input                                     │
│         - Call Service                                       │
│         - Return Response                                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    Service Layer                             │
│                  (Business Logic)                            │
│         - Implement Use Cases                                │
│         - Coordinate Repositories                            │
│         - Handle Transactions                                │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                  Repository Layer                            │
│                 (Data Access Logic)                          │
│         - CRUD Operations                                    │
│         - Database Queries                                   │
│         - ORM Abstraction                                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    Database Layer                            │
│                   (PostgreSQL/MySQL)                         │
└─────────────────────────────────────────────────────────────┘
```

## 🏗️ Project Structure

```
starter-gofiber/
├── assets/               # Static assets
│   ├── certs/           # SSL certificates
│   ├── docker/          # Docker configurations
│   └── rbac/            # RBAC policies (Casbin)
├── cmd/                 # Application entry points
│   ├── api/            # Main API server
│   │   └── main.go     # API entry point
│   └── worker/         # Background worker server
│       └── main.go     # Worker entry point
├── internal/            # Private application code
│   ├── config/         # Configuration management
│   │   ├── app.go      # App config
│   │   ├── database.go  # Database config
│   │   └── permission.go # Casbin enforcer
│   ├── domain/          # Domain layer (entities & interfaces)
│   │   ├── user/        # User domain
│   │   │   ├── entity.go      # User entity
│   │   │   ├── dto.go         # User DTOs
│   │   │   ├── repository.go # User repository interface
│   │   │   └── service.go    # User service interface
│   │   └── post/        # Post domain
│   │       ├── entity.go      # Post entity
│   │       ├── dto.go         # Post DTOs
│   │       ├── repository.go # Post repository interface
│   │       └── service.go    # Post service interface
│   ├── handler/         # HTTP handlers
│   │   ├── http/        # HTTP request handlers
│   │   │   ├── auth.go  # Auth endpoints
│   │   │   └── post.go  # Post endpoints
│   │   └── middleware/  # HTTP middlewares
│   │       ├── auth.go  # JWT authentication
│   │       └── authz.go # Authorization (Casbin)
│   ├── repository/      # Data access implementations
│   │   └── postgres/    # PostgreSQL repository
│   │       ├── user.go  # User repository implementation
│   │       └── post.go  # Post repository implementation
│   ├── service/         # Business logic implementations
│   │   ├── auth/        # Auth service
│   │   │   └── service.go
│   │   └── post/        # Post service
│   │       └── service.go
│   └── worker/          # Background worker (Asynq)
│       ├── jobs.go      # Job definitions
│       └── handlers.go  # Job handlers
├── pkg/                 # Public library code
│   ├── dto/             # Shared Data Transfer Objects
│   │   ├── paginate.go  # Pagination DTO
│   │   ├── response.go  # Response wrapper
│   │   └── token.go     # Token DTOs
│   ├── apierror/        # API error types
│   │   └── error.go
│   ├── crypto/          # Cryptographic utilities
│   │   ├── hash.go      # Hashing utilities
│   │   ├── password.go  # Password hashing
│   │   └── jwt.go       # JWT utilities
│   ├── database/        # Database utilities
│   │   ├── bulk.go      # Bulk operations
│   │   └── pagination.go # Pagination helpers
│   ├── pagination/      # Pagination utilities
│   │   └── cursor.go    # Cursor-based pagination
│   ├── utils/           # General utilities
│   │   ├── export.go    # Data export
│   │   └── filter.go    # Filtering utilities
│   └── response/        # Response utilities
│       └── response.go  # Response builder
├── router/              # Route definitions
│   ├── router.go        # Main router
│   ├── auth.go          # Auth routes
│   └── post.go          # Post routes
├── tests/               # Test files
│   ├── setup_test.go    # Test setup
│   ├── auth_test.go     # Auth tests
│   └── post_test.go     # Post tests
├── variables/           # Constants
│   └── constant.go      # App constants
├── docs/                # Documentation
├── go.mod               # Go modules
├── .env                 # Environment variables
└── docker-compose.yml   # Docker compose config
```

## 🎯 Design Patterns

### 1. Repository Pattern

**Purpose**: Abstraksi data access layer untuk memudahkan testing dan switching database.

**Implementation**:

```go
// repository/repository.go
type Repository[T any] interface {
    Create(data *T) error
    Update(id uint, data *T) error
    Delete(id uint) error
    FindById(id uint) (*T, error)
}

// repository/user.go
type UserRepository struct {
    DB *gorm.DB
}

func (r *UserRepository) Create(user *user.User) error {
    return r.DB.Create(user).Error
}
```

**Benefits**:
- ✅ Separation of concerns
- ✅ Easy to mock for testing
- ✅ Database agnostic
- ✅ Centralized query logic

### 2. Dependency Injection

**Purpose**: Loose coupling between components.

**Implementation**:

```go
// service/auth.go
type AuthService struct {
    userRepo     user.Repository
    tokenRepo    user.Repository
    passwordRepo user.Repository
}

func NewAuthService(
    userRepo user.Repository,
    tokenRepo user.Repository,
    passwordRepo user.Repository,
) *AuthService {
    return &AuthService{
        userRepo:     userRepo,
        tokenRepo:    tokenRepo,
        passwordRepo: passwordRepo,
    }
}
```

**Benefits**:
- ✅ Testability (easy to inject mocks)
- ✅ Flexibility (swap implementations)
- ✅ Clear dependencies
- ✅ Better code organization

### 3. Service Layer Pattern

**Purpose**: Encapsulate business logic terpisah dari HTTP layer.

**Implementation**:

```go
// service/auth.go
func (s *AuthService) Register(req user.RegisterRequest) (*user.LoginResponse, error) {
    // 1. Validate business rules
    if exist := s.userRepo.ExistEmail(req.Email); exist {
        return nil, &apierror.BadRequestError{
            Message: "Email already registered",
            Order:   "S1",
        }
    }
    
    // 2. Hash password
    hashedPassword, err := crypto.HashPassword(req.Password)
    
    // 3. Create user
    user := &user.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
    }
    
    // 4. Save to database
    err = s.userRepo.Create(user)
    
    // 5. Generate tokens
    accessToken, _ := crypto.GenerateJWT(user)
    refreshToken, _ := crypto.GenerateRefreshToken(user)
    
    // 6. Return response
    return &user.LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}
```

**Benefits**:
- ✅ Business logic reusability
- ✅ Transaction management
- ✅ Independent testing
- ✅ Single responsibility

### 4. DTO (Data Transfer Object) Pattern

**Purpose**: Memisahkan struktur request/response dari entity database.

**Implementation**:

```go
// dto/user.go
type RegisterRequest struct {
    Name     string `json:"name" validate:"required,min=3"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
    AccessToken  string      `json:"access_token"`
    RefreshToken string      `json:"refresh_token"`
    User         UserProfile `json:"user"`
}

// entity/user.go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"type:varchar(200)"`
    Email     string    `gorm:"uniqueIndex"`
    Password  string    `gorm:"type:varchar(150)"` // Tidak di-expose
    Role      UserRole  `gorm:"type:enum('admin','user')"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Benefits**:
- ✅ API contract clarity
- ✅ Hide sensitive data (password)
- ✅ Validation at boundary
- ✅ Version compatibility

### 5. Middleware Pattern

**Purpose**: Cross-cutting concerns (auth, logging, error handling).

**Implementation**:

```go
// middleware/auth.go
func AuthMiddleware(c *fiber.Ctx) error {
    // 1. Extract token from header
    token := c.Get("Authorization")
    
    // 2. Validate token
    user, err := crypto.GetUserFromToken(token)
    if err != nil {
        return apierror.UnauthorizedError{Message: "Invalid token"}
    }
    
    // 3. Store user in context
    c.Locals("user", user)
    
    // 4. Continue to next handler
    return c.Next()
}

// middleware/authz.go (Casbin)
func LoadAuthzMiddleware(enforcer *casbin.Enforcer) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(*user.User)
        
        // Check permission
        ok, _ := enforcer.Enforce(user.Role, c.Path(), c.Method())
        if !ok {
            return apierror.ForbiddenError{Message: "Access denied"}
        }
        
        return c.Next()
    }
}
```

**Benefits**:
- ✅ Reusable across routes
- ✅ Separation of concerns
- ✅ Clean handler code
- ✅ Easy to add/remove

### 6. Error Handling Pattern

**Purpose**: Consistent error response dengan proper HTTP status codes.

**Implementation**:

```go
// helper/error.go
type BadRequestError struct {
    Message string
    Order   string // "S1", "H2", etc for tracking
}

func (e BadRequestError) Error() string {
    return e.Message
}

// helper/error.go
func ErrorHelper(c *fiber.Ctx, err error) error {
    switch e := err.(type) {
    case *BadRequestError:
        return c.Status(400).JSON(Response{
            Success: false,
            Message: e.Message,
            Order:   e.Order,
        })
    case *UnauthorizedError:
        return c.Status(401).JSON(Response{
            Success: false,
            Message: e.Message,
        })
    case *ForbiddenError:
        return c.Status(403).JSON(Response{
            Success: false,
            Message: e.Message,
        })
    default:
        return c.Status(500).JSON(Response{
            Success: false,
            Message: "Internal server error",
        })
    }
}

// Usage in handler
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    resp, err := h.authService.Login(req)
    if err != nil {
        return apierror.ErrorHelper(c, err) // Auto handle error type
    }
    return c.JSON(resp)
}
```

**Benefits**:
- ✅ Consistent error format
- ✅ Proper HTTP status codes
- ✅ Error tracking dengan Order
- ✅ Clean error handling

## 🔄 Request Flow

### Complete Request Lifecycle

```
1. Client Request
   │
   ├─→ 2. Fiber Router (router/router.go)
   │       │
   │       ├─→ 3. Middleware Chain
   │       │       ├─→ CORS
   │       │       ├─→ Rate Limiter
   │       │       ├─→ Logger
   │       │       ├─→ Auth Middleware (jika protected)
   │       │       └─→ Authz Middleware (jika butuh permission)
   │       │
   │       └─→ 4. Handler (handler/*.go)
   │               ├─→ Parse Request
   │               ├─→ Validate DTO
   │               └─→ Call Service
   │
   ├─→ 5. Service Layer (internal/service/*.go)
   │       ├─→ Business Logic
   │       ├─→ Call Repository
   │       └─→ Build Response DTO
   │
   ├─→ 6. Repository Layer (internal/repository/postgres/*.go)
   │       ├─→ Build Query
   │       ├─→ Execute via GORM
   │       └─→ Return Entity
   │
   ├─→ 7. Database (PostgreSQL/MySQL)
   │
   └─→ 8. Response to Client
```

### Example: User Registration Flow

```go
// 1. Client sends POST /api/v1/auth/register
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123"
}

// 2. Router matches route
router.Post("/register", authHandler.Register)

// 3. No middleware (public endpoint)

// 4. Handler validates and calls service
func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req user.RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return apierror.BadRequestError{Message: "Invalid request"}
    }
    
    resp, err := h.authService.Register(req)
    if err != nil {
        return apierror.ErrorHelper(c, err)
    }
    
    return response.Response(c, 201, "User registered", resp)
}

// 5. Service implements business logic
func (s *AuthService) Register(req user.RegisterRequest) (*user.LoginResponse, error) {
    // Check duplicate email
    if s.userRepo.ExistEmail(req.Email) {
        return nil, &apierror.BadRequestError{
            Message: "Email already exists",
            Order:   "S1",
        }
    }
    
    // Hash password
    hashedPassword, _ := crypto.HashPassword(req.Password)
    
    // Create user entity
    user := &user.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
        Role:     user.UserRoleUser,
    }
    
    // Save via repository
    if err := s.userRepo.Create(user); err != nil {
        return nil, &apierror.InternalServerError{
            Message: "Failed to create user",
            Order:   "S2",
        }
    }
    
    // Generate tokens
    accessToken, _ := crypto.GenerateJWT(user)
    refreshToken, _ := crypto.GenerateRefreshToken(user)
    
    return &user.LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User: pkg/dto.UserProfile{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
            Role:  string(user.Role),
        },
    }, nil
}

// 6. Repository executes database operation
func (r *UserRepository) Create(user *user.User) error {
    return r.DB.Create(user).Error
}

// 7. GORM generates SQL
INSERT INTO users (name, email, password, role, created_at, updated_at)
VALUES ('John Doe', 'john@example.com', '$2a$10$...', 'user', NOW(), NOW());

// 8. Response to client
{
    "success": true,
    "message": "User registered",
    "data": {
        "access_token": "eyJhbGc...",
        "refresh_token": "eyJhbGc...",
        "user": {
            "id": 1,
            "name": "John Doe",
            "email": "john@example.com",
            "role": "user"
        }
    }
}
```

## 🔐 Authentication & Authorization

### Authentication (JWT)

```
1. User Login
   ├─→ Service validates credentials
   ├─→ Generate Access Token (15 min expiry)
   ├─→ Generate Refresh Token (7 days expiry)
   └─→ Store Refresh Token in database

2. Access Protected Endpoint
   ├─→ Client sends: Authorization: Bearer <access_token>
   ├─→ Auth Middleware validates token
   ├─→ Extract user from token
   ├─→ Store user in c.Locals("user")
   └─→ Continue to handler

3. Token Expired
   ├─→ Client sends refresh token
   ├─→ Service validates refresh token
   ├─→ Generate new access token
   └─→ Optionally rotate refresh token
```

### Authorization (Casbin RBAC)

```
// assets/rbac/model.conf
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act

// assets/rbac/policy.csv
p, admin, /api/v1/posts, POST
p, admin, /api/v1/posts/:id, PUT
p, admin, /api/v1/posts/:id, DELETE
p, user, /api/v1/posts, POST
p, user, /api/v1/posts/:id, PUT  # Only own posts (checked in service)

g, user@example.com, user
g, admin@example.com, admin
```

## 📦 Dependency Management

### Main Dependencies

```go
// Framework
github.com/gofiber/fiber/v2          // Web framework
github.com/gofiber/contrib/fibercasbin // Casbin integration

// Database
gorm.io/gorm                          // ORM
gorm.io/driver/postgres               // PostgreSQL driver
gorm.io/driver/mysql                  // MySQL driver
gorm.io/driver/sqlite                 // SQLite (for testing)

// Security
github.com/golang-jwt/jwt/v5          // JWT tokens
golang.org/x/crypto/bcrypt            // Password hashing
github.com/casbin/casbin/v2           // RBAC authorization

// Testing
github.com/stretchr/testify           // Test assertions
```

## 🎨 Code Organization Principles

### 1. Single Responsibility
Setiap file/struct hanya punya satu tanggung jawab.

```go
// ✅ Good
// internal/handler/http/auth.go - Handle HTTP requests only
// service/auth.go - Business logic only
// repository/user.go - Database operations only

// ❌ Bad
// auth.go - Mix handler, service, repository
```

### 2. Dependency Direction
Dependencies harus searah: Handler → Service → Repository

```go
// ✅ Good
Handler depends on Service
Service depends on Repository
Repository depends on Database

// ❌ Bad
Service depends on Handler
Repository depends on Service
```

### 3. Interface Segregation
Gunakan interface untuk abstraksi.

```go
// ✅ Good
type UserRepository interface {
    Create(user *user.User) error
    FindByEmail(email string) (*user.User, error)
}

// Mudah di-mock untuk testing
type MockUserRepository struct {
    mock.Mock
}
```

### 4. Error Propagation
Error harus di-propagate dengan konteks yang jelas.

```go
// ✅ Good
if err := repo.Create(user); err != nil {
    return nil, &apierror.InternalServerError{
        Message: "Failed to create user",
        Order:   "S2", // Track error location
    }
}

// ❌ Bad
if err := repo.Create(user); err != nil {
    return nil, err // Lost context
}
```

## 🧪 Testing Strategy

### Unit Testing
Test individual components dengan mocked dependencies.

```go
// service/auth_test.go
func TestRegister_Success(t *testing.T) {
    mockUserRepo := new(MockUserRepository)
    mockUserRepo.On("ExistEmail", "test@example.com").Return(false)
    mockUserRepo.On("Create", mock.Anything).Return(nil)
    
    authService := NewAuthService(mockUserRepo, nil, nil)
    
    resp, err := authService.Register(user.RegisterRequest{
        Email: "test@example.com",
        // ...
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    mockUserRepo.AssertExpectations(t)
}
```

### Integration Testing
Test end-to-end dengan real HTTP requests.

```go
// tests/auth_test.go
func (s *AuthTestSuite) TestRegister_Success() {
    req := user.RegisterRequest{
        Name:     "Test User",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    resp := s.MakeRequest("POST", "/api/v1/auth/register", req)
    s.AssertSuccessResponse(resp, 201)
}
```

## 🚀 Performance Considerations

### Database Query Optimization
1. **Pagination**: Gunakan `LIMIT` dan `OFFSET`
2. **Eager Loading**: Preload relations untuk avoid N+1 query
3. **Indexing**: Index pada kolom yang sering di-query

```go
// Pagination
func (r *PostRepository) All(page, pageSize int) ([]post.Post, error) {
    var posts []post.Post
    offset := (page - 1) * pageSize
    
    err := r.DB.
        Preload("User").  // Eager loading
        Order("posts.id desc").
        Limit(pageSize).
        Offset(offset).
        Find(&posts).Error
        
    return posts, err
}
```

### Caching Strategy (Future)
1. Redis untuk session storage
2. Cache query results yang jarang berubah
3. Invalidate cache on data update

## 📚 Additional Resources

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
- [Casbin Documentation](https://casbin.org/docs/overview)
- [Testing Best Practices](https://github.com/stretchr/testify)

---

**Last Updated**: December 31, 2025
