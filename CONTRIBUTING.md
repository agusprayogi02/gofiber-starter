# Contributing Guide

Terima kasih atas minat Anda untuk berkontribusi pada Starter Template Go Fiber! 🎉

Panduan ini akan membantu Anda memahami cara berkontribusi pada project ini dengan efektif.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)

## 🤝 Code of Conduct

### Our Pledge

Kami berkomitmen untuk menjadikan partisipasi dalam project ini bebas dari harassment untuk semua orang, terlepas dari:
- Usia, ukuran tubuh, disabilitas
- Etnis, identitas dan ekspresi gender
- Level pengalaman, kebangsaan
- Penampilan pribadi, ras, agama
- Identitas dan orientasi seksual

### Our Standards

**Perilaku yang Diharapkan** ✅:
- Menggunakan bahasa yang ramah dan inklusif
- Menghormati sudut pandang dan pengalaman yang berbeda
- Menerima kritik konstruktif dengan anggun
- Fokus pada yang terbaik untuk komunitas
- Menunjukkan empati terhadap anggota komunitas lainnya

**Perilaku yang Tidak Dapat Diterima** ❌:
- Penggunaan bahasa atau gambar seksual
- Trolling, komentar menghina/merendahkan
- Harassment publik atau pribadi
- Mempublikasikan informasi pribadi orang lain tanpa izin
- Perilaku lain yang tidak profesional atau tidak pantas

## 🚀 Getting Started

### Prerequisites

Pastikan Anda sudah install:
- Go 1.20 atau lebih tinggi
- Docker & Docker Compose
- Git
- Air (untuk live reload)
- Code editor (VS Code recommended)

### Fork dan Clone Repository

```bash
# 1. Fork repository via GitHub UI

# 2. Clone fork Anda
git clone https://github.com/YOUR_USERNAME/starter-gofiber.git
cd starter-gofiber

# 3. Add upstream remote
git remote add upstream https://github.com/ORIGINAL_OWNER/starter-gofiber.git

# 4. Verify remotes
git remote -v
```

### Setup Development Environment

```bash
# 1. Copy environment file
cp .env.example .env

# 2. Generate SSL certificate
cd assets/certs
openssl genpkey -algorithm RSA -out certificate.pem -pkeyopt rsa_keygen_bits:4096
cd ../..

# 3. Install dependencies
go mod download

# 4. Install air for hot reload
go install github.com/air-verse/air@latest

# 5. Run application
air

# Atau dengan Docker
docker compose up -d
```

### Verify Setup

```bash
# Test application
curl http://localhost:3000/health

# Expected response:
# {"success":true,"message":"Server is running"}

# Run tests
go test ./tests/...
```

## 🔄 Development Workflow

### 1. Create a Branch

Gunakan naming convention:
- `feature/nama-fitur` - untuk fitur baru
- `fix/nama-bug` - untuk bug fix
- `docs/topik` - untuk dokumentasi
- `refactor/area` - untuk refactoring
- `test/area` - untuk menambah tests

```bash
# Update main branch
git checkout main
git pull upstream main

# Create and checkout new branch
git checkout -b feature/add-user-profile

# Or
git checkout -b fix/login-validation-bug
```

### 2. Make Changes

```bash
# Edit files dengan code editor favorit Anda

# Check changed files
git status

# View changes
git diff
```

### 3. Test Your Changes

```bash
# Run tests
go test ./tests/... -v

# Run with coverage
bash scripts/test-coverage.sh

# Test specific package
go test ./service/... -v

# Manual testing
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com","password":"password123"}'
```

### 4. Commit Changes

```bash
# Stage changes
git add .

# Commit with meaningful message
git commit -m "feat: add user profile endpoint"

# Or stage specific files
git add handler/user.go service/user.go
git commit -m "feat: implement user profile logic"
```

### 5. Push to Your Fork

```bash
# Push branch
git push origin feature/add-user-profile

# If you updated commits
git push origin feature/add-user-profile --force
```

### 6. Create Pull Request

1. Buka repository Anda di GitHub
2. Klik "Compare & pull request"
3. Isi template PR dengan lengkap
4. Klik "Create pull request"

## 📝 Coding Standards

### Go Code Style

Ikuti [Effective Go](https://golang.org/doc/effective_go.html) dan [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

**Format Code**:
```bash
# Format all files
go fmt ./...

# Or use gofmt
gofmt -w .

# Check with golangci-lint (recommended)
golangci-lint run
```

### Naming Conventions

**Variables & Functions**:
```go
// ✅ Good - camelCase untuk private
var userRepository UserRepository
func getUserByEmail(email string) (*entity.User, error) {}

// ✅ Good - PascalCase untuk public
var UserService *AuthService
func GetUserProfile(id uint) (*dto.UserProfile, error) {}

// ❌ Bad - snake_case tidak digunakan di Go
var user_repository UserRepository
func get_user_by_email(email string) {}
```

**Structs & Interfaces**:
```go
// ✅ Good - singular, descriptive
type User struct {}
type UserRepository interface {}
type AuthService struct {}

// ❌ Bad - plural, abbreviated
type Users struct {}
type UsrRepo interface {}
type AuthSvc struct {}
```

**Constants**:
```go
// ✅ Good - PascalCase atau ALL_CAPS untuk exported
const (
    UserRoleAdmin UserRole = "admin"
    UserRoleUser  UserRole = "user"
)

const DefaultPageSize = 10

// ✅ Good - camelCase untuk private
const maxRetries = 3
```

### File Organization

```go
// ✅ Good - group imports
import (
    // Standard library
    "fmt"
    "time"
    
    // External packages
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    
    // Internal packages
    "starter-gofiber/dto"
    "starter-gofiber/entity"
    "starter-gofiber/helper"
)

// ✅ Good - struct definition
type UserService struct {
    repo repository.UserRepository
}

// ✅ Good - constructor
func NewUserService(repo repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

// ✅ Good - methods grouped by functionality
func (s *UserService) GetProfile(id uint) (*dto.UserProfile, error) {
    // Implementation
}

func (s *UserService) UpdateProfile(id uint, req dto.UpdateProfileRequest) error {
    // Implementation
}
```

### Error Handling

```go
// ✅ Good - use custom error types
if err := s.repo.Create(user); err != nil {
    return nil, &helper.InternalServerError{
        Message: "Failed to create user",
        Order:   "S1",
    }
}

// ✅ Good - check error immediately
user, err := s.repo.FindByEmail(email)
if err != nil {
    return nil, err
}

// ❌ Bad - ignore errors
user, _ := s.repo.FindByEmail(email)

// ❌ Bad - generic error without context
if err != nil {
    return nil, err
}
```

### Comments

```go
// ✅ Good - public functions have comments
// GetUserProfile retrieves user profile information by user ID.
// Returns UserProfile DTO or error if user not found.
func GetUserProfile(id uint) (*dto.UserProfile, error) {
    // Implementation
}

// ✅ Good - complex logic explained
// Calculate discount based on user tier and purchase amount.
// Tier 1: 5%, Tier 2: 10%, Tier 3: 15%
discount := calculateDiscount(user.Tier, amount)

// ❌ Bad - stating the obvious
// Get user by ID
user := getUserByID(id)

// ❌ Bad - commented out code (delete instead)
// oldImplementation()
newImplementation()
```

## 🧪 Testing Guidelines

### Test File Structure

```go
// tests/user_test.go
package tests

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
    suite.Suite
    app *fiber.App
    db  *gorm.DB
}

func (s *UserTestSuite) SetupTest() {
    // Setup before each test
    s.db = SetupTestDB()
    s.app = SetupTestApp(s.db)
}

func (s *UserTestSuite) TearDownTest() {
    // Cleanup after each test
    CleanupTestDB(s.db)
}

func (s *UserTestSuite) TestGetProfile_Success() {
    // Arrange
    user := createTestUser(s.db, "test@example.com")
    
    // Act
    resp := s.MakeRequest("GET", "/api/v1/users/profile", nil)
    
    // Assert
    s.AssertSuccessResponse(resp, 200)
    // More assertions...
}

func TestUserTestSuite(t *testing.T) {
    suite.Run(t, new(UserTestSuite))
}
```

### Test Coverage Requirements

- **Minimum coverage**: 60% overall
- **Critical paths**: 80%+ (auth, payment, etc.)
- **New features**: Must include tests
- **Bug fixes**: Add regression tests

```bash
# Check coverage
bash scripts/test-coverage.sh

# Target breakdown:
# - Handler: 70%+
# - Service: 80%+
# - Repository: 70%+
# - Helper: 75%+
```

### Test Naming

```go
// Pattern: Test{FunctionName}_{Scenario}_{ExpectedResult}

func TestRegister_ValidInput_Success() {}
func TestRegister_DuplicateEmail_ReturnsError() {}
func TestRegister_InvalidEmail_ReturnsValidationError() {}

func TestLogin_ValidCredentials_ReturnsTokens() {}
func TestLogin_InvalidPassword_ReturnsUnauthorized() {}
func TestLogin_UserNotFound_ReturnsUnauthorized() {}
```

## 📝 Commit Guidelines

### Commit Message Format

Gunakan [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: Fitur baru
- `fix`: Bug fix
- `docs`: Perubahan dokumentasi
- `style`: Format code (tidak mengubah logic)
- `refactor`: Refactoring code
- `test`: Menambah atau update tests
- `chore`: Maintenance tasks (dependencies, config, dll)
- `perf`: Performance improvements

### Examples

```bash
# Good commits ✅
git commit -m "feat(auth): add email verification endpoint"
git commit -m "fix(user): resolve null pointer on profile update"
git commit -m "docs(api): update authentication documentation"
git commit -m "test(auth): add integration tests for login flow"
git commit -m "refactor(service): extract validation logic to helper"
git commit -m "chore(deps): update fiber to v2.52.10"

# Bad commits ❌
git commit -m "update code"
git commit -m "fix bug"
git commit -m "changes"
git commit -m "WIP"
```

### Commit Body (Optional tapi Recommended)

```bash
git commit -m "feat(auth): add two-factor authentication

Implement 2FA using TOTP (Time-based One-Time Password).
Users can enable/disable 2FA in account settings.
QR code generated for authenticator app setup.

Closes #123"
```

## 🔀 Pull Request Process

### Before Creating PR

**Checklist**:
- [ ] Code mengikuti style guide
- [ ] Tests pass (`go test ./tests/...`)
- [ ] Coverage tidak menurun
- [ ] Dokumentasi updated (jika perlu)
- [ ] No merge conflicts dengan main
- [ ] Commits mengikuti conventional commits

```bash
# Update branch dengan main
git checkout main
git pull upstream main
git checkout feature/your-feature
git rebase main

# Resolve conflicts jika ada
# Test setelah rebase
go test ./tests/...
```

### PR Title

Gunakan format yang sama dengan commit messages:

```
feat(auth): add two-factor authentication support
fix(user): resolve profile image upload issue
docs(deployment): add Docker Compose production guide
```

### PR Description Template

```markdown
## Description
Deskripsi singkat tentang perubahan yang dilakukan.

## Type of Change
- [ ] Bug fix (non-breaking change yang memperbaiki issue)
- [ ] New feature (non-breaking change yang menambah functionality)
- [ ] Breaking change (fix/feature yang membuat existing functionality tidak kompatibel)
- [ ] Documentation update

## Related Issues
Closes #123
Fixes #456

## Changes Made
- Implement 2FA endpoint (`POST /api/v1/auth/2fa/setup`)
- Add TOTP secret generation
- Create QR code response
- Update user entity with 2FA fields
- Add integration tests for 2FA flow

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed
- [ ] Coverage: 75% (+5% from previous)

## Screenshots (jika ada UI changes)
[Attach screenshots]

## Checklist
- [ ] My code follows the style guidelines
- [ ] I have performed a self-review
- [ ] I have commented my code where necessary
- [ ] I have updated the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix/feature works
- [ ] New and existing tests pass locally
- [ ] Any dependent changes have been merged
```

### Review Process

1. **Automated Checks**: Tests, linting, coverage
2. **Code Review**: Minimal 1 approval dari maintainer
3. **Address Feedback**: Update PR berdasarkan review
4. **Final Check**: Maintainer akan final review
5. **Merge**: Squash and merge ke main

### Responding to Reviews

```bash
# Make requested changes
git add .
git commit -m "refactor: extract validation to helper function"
git push origin feature/your-feature

# Or fixup commits (akan di-squash saat merge)
git commit --fixup <commit-hash>
git push origin feature/your-feature
```

## 🐛 Issue Reporting

### Before Creating Issue

1. **Search existing issues** - Mungkin sudah ada yang report
2. **Check documentation** - Mungkin sudah dijelaskan
3. **Update dependencies** - Pastikan menggunakan versi terbaru
4. **Reproduce** - Pastikan bug bisa di-reproduce

### Bug Report Template

```markdown
**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Send POST request to '/api/v1/auth/login'
2. With payload: { "email": "test@example.com", "password": "wrong" }
3. Expect 401 response
4. See error (500 instead)

**Expected behavior**
Should return 401 Unauthorized with error message.

**Actual behavior**
Returns 500 Internal Server Error.

**Screenshots/Logs**
```
Error: pq: null value in column "last_login"
```

**Environment:**
 - OS: Ubuntu 22.04
 - Go version: 1.20.5
 - Fiber version: v2.52.10
 - Database: PostgreSQL 15

**Additional context**
This happens only when user has never logged in before.
```

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
I'm always frustrated when I need to manually verify emails.

**Describe the solution you'd like**
Add email verification feature:
- Send verification email on registration
- Verify email endpoint
- Resend verification email option

**Describe alternatives you've considered**
- Manual verification by admin
- No verification (current state)

**Additional context**
Many modern apps require email verification for security.
This would improve trust and reduce fake accounts.

**Would you like to implement this feature?**
Yes, I can submit a PR if approved.
```

## 🎯 Areas for Contribution

### High Priority

- [ ] Email service integration (SendGrid/Mailgun)
- [ ] Redis caching layer
- [ ] Database migration system
- [ ] Structured logging (zap/logrus)
- [ ] Health check endpoint
- [ ] Swagger/OpenAPI documentation
- [ ] More comprehensive tests

### Medium Priority

- [ ] File upload improvements
- [ ] User profile management
- [ ] Admin dashboard endpoints
- [ ] Webhook system
- [ ] Background job queue
- [ ] Performance optimizations

### Good First Issues

Look for issues labeled with:
- `good first issue` - Easy untuk pemula
- `help wanted` - Butuh kontributor
- `documentation` - Documentation improvements

## 💡 Tips for Success

1. **Start Small** - Mulai dari bug fixes atau documentation
2. **Ask Questions** - Jangan ragu bertanya di issues/discussions
3. **Read Code** - Pahami codebase sebelum contribute
4. **Follow Standards** - Ikuti conventions yang ada
5. **Be Patient** - Review process butuh waktu
6. **Stay Updated** - Follow repository untuk updates

## 📞 Getting Help

- **GitHub Issues**: Technical questions dan bug reports
- **GitHub Discussions**: General questions dan ideas
- **Email**: maintainer@example.com (untuk hal privat)

## 🙏 Recognition

Semua contributors akan dicantumkan di:
- README.md Contributors section
- Release notes (untuk significant contributions)
- Special thanks untuk major features

---

**Thank you for contributing!** 🎉

Your contributions make this project better for everyone.

**Last Updated**: December 31, 2025
