# Test Coverage Report

Dokumentasi lengkap untuk sistem Test Coverage Report pada Starter Template Go Fiber.

## 📊 Overview

Test Coverage Report membantu mengukur seberapa banyak kode aplikasi yang tercovered oleh automated tests. Coverage yang baik menunjukkan kualitas testing yang baik dan mengurangi bug di production.

### Packages yang Diukur

Coverage mencakup package-package berikut:
- `./handler/...` - HTTP handlers
- `./service/...` - Business logic
- `./repository/...` - Database operations
- `./middleware/...` - Middleware functions
- `./helper/...` - Utility functions

## 🚀 Quick Start

### Menjalankan Test Coverage

```bash
# Menggunakan script yang sudah disediakan
bash scripts/test-coverage.sh

# Atau dengan make (jika ada Makefile)
make test-coverage

# Manual dengan go test
go test ./tests/... -v -coverprofile=coverage.out -covermode=atomic \
  -coverpkg=./handler/...,./service/...,./repository/...,./middleware/...,./helper/...
```

### Melihat Hasil Coverage

```bash
# 1. Terminal output
go tool cover -func=coverage.out

# 2. HTML Report (Lebih interaktif)
go tool cover -html=coverage.out -o coverage.html
firefox coverage.html  # atau browser lain
```

## 📁 Files yang Dihasilkan

### coverage.out
- **Format**: Binary coverage profile
- **Kegunaan**: Input untuk tools coverage
- **Lokasi**: Root project
- **Git**: Sudah di-ignore di `.gitignore`

### coverage.html
- **Format**: HTML interactive report
- **Kegunaan**: Visual representation dengan syntax highlighting
- **Fitur**: 
  - ✅ Green: Lines yang tercovered
  - ❌ Red: Lines yang tidak tercovered
  - 🔍 Klik function untuk detail

## 📖 Membaca Coverage Report

### Terminal Output Example

```bash
starter-gofiber/internal/handler/http/auth.go:22:     Register        77.8%
starter-gofiber/internal/handler/http/auth.go:43:     Login           77.8%
starter-gofiber/internal/handler/http/auth.go:68:     RefreshToken    0.0%
starter-gofiber/internal/service/auth.go:33:     Register        81.8%
total:                                  (statements)    23.2%
```

**Penjelasan:**
- `77.8%` - Persentase lines yang tereksekusi dalam test
- `0.0%` - Function tidak ada test coverage sama sekali
- `23.2%` - **Total coverage** dari semua package

### Coverage Levels

| Coverage | Status | Keterangan |
|----------|--------|------------|
| 0-30% | 🔴 Poor | Perlu banyak improvement |
| 30-50% | 🟡 Fair | Ada coverage tapi masih kurang |
| 50-70% | 🟢 Good | Coverage cukup baik |
| 70-85% | 🟢 Very Good | Coverage sangat baik |
| 85-100% | 🟢 Excellent | Coverage hampir sempurna |

## ⚙️ Konfigurasi

### Coverage Threshold

Script `test-coverage.sh` menggunakan threshold **60%**. Bisa disesuaikan:

```bash
# Edit scripts/test-coverage.sh
THRESHOLD=60  # Ubah sesuai kebutuhan

# Threshold terlalu tinggi (>80%): Sulit dicapai, bisa memperlambat development
# Threshold terlalu rendah (<40%): Tidak efektif mendeteksi code yang untested
```

### Coverage Mode

```bash
# atomic - Thread-safe, cocok untuk concurrent tests
-covermode=atomic

# set - Hanya track apakah line dieksekusi (lebih cepat)
-covermode=set

# count - Hitung berapa kali line dieksekusi
-covermode=count
```

## 🎯 Strategy Meningkatkan Coverage

### 1. Identifikasi Gap

```bash
# Lihat function yang 0% coverage
go tool cover -func=coverage.out | grep "0.0%"

# Prioritaskan:
# - Critical business logic (service layer)
# - Handler endpoints yang sering dipakai
# - Helper functions yang kompleks
```

### 2. Tambah Test Cases

**Example**: Meningkatkan coverage untuk RefreshToken

```go
// tests/auth_test.go
func (s *AuthTestSuite) TestRefreshToken_Success() {
    // 1. Login dulu untuk dapat refresh token
    loginReq := user.LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    }
    loginResp := s.MakeRequest("POST", "/api/v1/auth/login", loginReq)
    // ... assertions

    // 2. Extract refresh token
    var loginData map[string]interface{}
    json.Unmarshal(loginResp.Body.Bytes(), &loginData)
    refreshToken := loginData["data"].(map[string]interface{})["refresh_token"].(string)

    // 3. Test refresh endpoint
    refreshReq := dto.RefreshTokenRequest{RefreshToken: refreshToken}
    resp := s.MakeRequest("POST", "/api/v1/auth/refresh", refreshReq)
    
    s.AssertSuccessResponse(resp, 200)
    // Coverage RefreshToken handler naik dari 0% → ~70%
}
```

### 3. Test Edge Cases

```go
// Test success case (happy path)
func (s *AuthTestSuite) TestLogin_Success() { /* ... */ }

// Test error cases (edge cases)
func (s *AuthTestSuite) TestLogin_InvalidCredentials() { /* ... */ }
func (s *AuthTestSuite) TestLogin_EmptyEmail() { /* ... */ }
func (s *AuthTestSuite) TestLogin_EmptyPassword() { /* ... */ }
func (s *AuthTestSuite) TestLogin_EmailNotFound() { /* ... */ }
func (s *AuthTestSuite) TestLogin_UnverifiedAccount() { /* ... */ }

// Coverage Login handler: 0% → 90%+
```

## 📈 Coverage Report di CI/CD

### GitHub Actions Example

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.20'
      
      - name: Run tests with coverage
        run: bash scripts/test-coverage.sh
      
      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
          fail_ci_if_error: true
```

### GitLab CI Example

```yaml
# .gitlab-ci.yml
test:
  stage: test
  script:
    - bash scripts/test-coverage.sh
  coverage: '/total:\s+\(statements\)\s+(\d+\.\d+)%/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.xml
```

## 🔍 Best Practices

### ✅ DO

1. **Test Business Logic Dulu**
   ```
   Priority: Service > Repository > Handler > Helper > Middleware
   ```

2. **Mock External Dependencies**
   ```go
   // Gunakan in-memory database untuk testing
   db := setupTestDB() // SQLite in-memory
   
   // Mock email service
   emailService := &MockEmailService{}
   ```

3. **Test Happy Path & Error Cases**
   ```
   Setiap function minimal 2 tests:
   - Success case (expected behavior)
   - Error case (error handling)
   ```

4. **Gunakan Table-Driven Tests untuk Variasi Input**
   ```go
   tests := []struct {
       name    string
       input   user.LoginRequest
       wantErr bool
   }{
       {"Valid credentials", validReq, false},
       {"Invalid email", invalidEmailReq, true},
       {"Invalid password", invalidPassReq, true},
   }
   ```

### ❌ DON'T

1. **Jangan Kejar 100% Coverage Buta**
   - 100% coverage ≠ bug-free code
   - Fokus pada quality test, bukan quantity

2. **Jangan Test Generated Code**
   ```go
   // Skip coverage untuk:
   // - Getter/Setter sederhana
   // - Auto-generated code
   // - Simple constructors
   ```

3. **Jangan Skip Integration Tests**
   ```
   Unit tests: Fast, isolated
   Integration tests: Slower, tapi test real flow
   
   Keduanya penting!
   ```

## 🛠️ Troubleshooting

### Coverage 0% padahal ada tests?

```bash
# ❌ Salah: Tidak specify -coverpkg
go test ./tests/... -coverprofile=coverage.out

# ✅ Benar: Specify package yang mau diukur
go test ./tests/... -coverprofile=coverage.out \
  -coverpkg=./handler/...,./service/...
```

### Coverage tidak update setelah tambah test?

```bash
# Clear cache dan re-run
go clean -testcache
bash scripts/test-coverage.sh
```

### Test timeout?

```bash
# Increase timeout untuk integration tests
go test ./tests/... -timeout=5m -coverprofile=coverage.out
```

## 📚 Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Go Coverage Guide](https://go.dev/blog/cover)
- [Effective Go Testing](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Testify Suite Pattern](https://github.com/stretchr/testify#suite-package)

## 🎯 Current Coverage Status

**Project Coverage**: 23.2%

### Coverage Breakdown

| Package | Coverage | Status |
|---------|----------|--------|
| Handler | ~50% | 🟡 Fair |
| Service | ~35% | 🟡 Fair |
| Repository | ~40% | 🟡 Fair |
| Middleware | ~50% | 🟡 Fair |
| Helper | ~45% | 🟡 Fair |

### Next Steps

1. ✅ **Test RefreshToken flow** (+5% coverage)
2. ✅ **Test Logout endpoint** (+3% coverage)
3. ✅ **Test Post CRUD with auth** (+8% coverage)
4. ✅ **Test error handling paths** (+10% coverage)

**Target**: 60% coverage dalam 2 minggu

---

**Last Updated**: December 31, 2025
