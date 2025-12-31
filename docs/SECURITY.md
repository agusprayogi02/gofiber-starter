# Security Features

Dokumentasi lengkap tentang fitur keamanan yang tersedia di starter kit ini.

## Daftar Isi

- [Rate Limiting](#rate-limiting)
- [CSRF Protection](#csrf-protection)
- [Security Headers](#security-headers)
- [HTTPS Enforcement](#https-enforcement)
- [Input Sanitization](#input-sanitization)
- [IP Filtering](#ip-filtering)
- [API Key Authentication](#api-key-authentication)
- [Encryption at Rest](#encryption-at-rest)
- [SQL Injection Prevention](#sql-injection-prevention)
- [XSS Protection](#xss-protection)

---

## Rate Limiting

Mencegah penyalahgunaan API dengan membatasi jumlah request per pengguna/IP.

### Konfigurasi

```go
// config/app.go
rateLimiter := middleware.NewUserRateLimiter(100, time.Minute)
app.Use(rateLimiter.Middleware())
```

### Parameter

- `max`: Jumlah maksimum request (default: 100)
- `window`: Durasi window time (default: 1 menit)

### Response Headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640000000
```

### Error Response (429)

```json
{
  "success": false,
  "message": "Too many requests. Please try again later.",
  "order": "TM1"
}
```

---

## CSRF Protection

Melindungi dari Cross-Site Request Forgery attacks.

### Setup

```go
// config/app.go
app.Use(middleware.CSRFMiddleware())
```

### Frontend Usage

```javascript
// Get CSRF token from cookie
const csrfToken = document.cookie
  .split('; ')
  .find(row => row.startsWith('csrf_'))
  .split('=')[1];

// Include in request headers
fetch('/api/posts', {
  method: 'POST',
  headers: {
    'X-CSRF-Token': csrfToken,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(data)
});
```

### Konfigurasi

- Token expiration: 24 jam
- Cookie settings: Strict SameSite, Secure (HTTPS), HTTPOnly
- Header: `X-CSRF-Token`

---

## Security Headers

Menambahkan HTTP security headers untuk melindungi dari berbagai serangan.

### Setup

```go
// config/app.go
app.Use(middleware.SecurityHeadersMiddleware())
```

### Headers Included

| Header | Value | Fungsi |
|--------|-------|--------|
| X-Frame-Options | SAMEORIGIN | Prevent clickjacking |
| X-Content-Type-Options | nosniff | Prevent MIME sniffing |
| X-XSS-Protection | 1; mode=block | Legacy XSS protection |
| Referrer-Policy | strict-origin-when-cross-origin | Control referrer info |
| Permissions-Policy | Various | Control browser features |
| Content-Security-Policy | default-src 'self' | Control resource loading |
| Strict-Transport-Security | max-age=31536000 | Force HTTPS (production only) |

### Custom CSP

```go
customCSP := "default-src 'self'; script-src 'self' 'unsafe-inline' cdn.example.com"
app.Use(middleware.CustomCSPMiddleware(customCSP))
```

---

## HTTPS Enforcement

Memaksa semua koneksi menggunakan HTTPS di production.

### Smart Redirect (Recommended)

```go
// config/app.go
app.Use(middleware.HTTPSRedirectMiddleware())
```

- Otomatis skip di development/test environment
- Exempt untuk endpoint health check dan metrics
- HTTP 301 permanent redirect

### Force HTTPS (Strict)

```go
app.Use(middleware.ForceHTTPSMiddleware())
```

- Selalu enforce HTTPS tanpa exception
- Gunakan hanya di production

---

## Input Sanitization

Membersihkan input dari HTML, XSS, dan SQL injection attempts.

### Strip All HTML

```go
import "starter-gofiber/helper"

clean := helper.SanitizeInput(userInput)
```

### Allow Safe HTML (UGC)

```go
clean := helper.SanitizeHTML(userInput)
// Allows: <b>, <i>, <em>, <strong>, <a>, <p>, etc.
```

### Validate No XSS

```go
if err := helper.ValidateNoXSS(userInput); err != nil {
    return err // "Potentially dangerous content detected"
}
```

### Validate No SQL Injection

```go
if err := helper.ValidateNoSQLInjection(userInput); err != nil {
    return err // "Potentially dangerous SQL pattern detected"
}
```

### Bulk Sanitization

```go
// Sanitize map
sanitized := helper.SanitizeMap(map[string]string{
    "name": "<script>alert('xss')</script>John",
    "bio": "<b>Developer</b>",
})

// Sanitize slice
sanitized := helper.SanitizeSlice([]string{
    "<script>evil</script>",
    "Safe text",
})
```

---

## IP Filtering

Kontrol akses berdasarkan IP address dengan whitelist/blacklist.

### Static Whitelist

```go
allowedIPs := []string{
    "192.168.1.100",
    "10.0.0.*",        // Wildcard support
    "203.0.113.0/24",  // CIDR notation (planned)
}
app.Use(middleware.IPWhitelistMiddleware(allowedIPs))
```

### Static Blacklist

```go
blockedIPs := []string{
    "123.45.67.89",
    "198.51.100.*",
}
app.Use(middleware.IPBlacklistMiddleware(blockedIPs))
```

### Dynamic Filter (Runtime Control)

```go
filter := middleware.DynamicIPFilterMiddleware()

// Add to whitelist at runtime
middleware.AddToWhitelist("192.168.1.200")

// Add to blacklist
middleware.AddToBlacklist("malicious.ip.address")

// Remove from whitelist
middleware.RemoveFromWhitelist("192.168.1.200")

// Remove from blacklist
middleware.RemoveFromBlacklist("malicious.ip.address")

app.Use(filter)
```

### Error Response (403)

```json
{
  "success": false,
  "message": "Access denied from your IP address",
  "order": "IF1"
}
```

---

## API Key Authentication

Alternative authentication menggunakan API key untuk machine-to-machine atau CLI access.

### Setup

```go
// main.go
middleware.InitAPIKeyMiddleware(config.DB)
```

### Generate API Key

```go
import "starter-gofiber/helper"

apiKey, err := helper.GenerateAPIKey(db, userID, "CLI Tool")
// Returns: "randomly-generated-32-character-key"
// ⚠️ Save this key! It's only shown once.
```

### Middleware Usage

```go
// Require API key
app.Use("/api/admin", middleware.APIKeyAuth())

// Optional: API key OR JWT
app.Use("/api/data", middleware.OptionalAPIKeyAuth())

// Either API key OR JWT (strict)
app.Use("/api/resource", middleware.APIKeyOrJWT())
```

### Client Request

```bash
curl -H "X-API-Key: your-api-key-here" \
     https://api.example.com/admin/users
```

### Management Functions

```go
// List user's API keys
keys, err := helper.ListAPIKeys(db, userID)

// Revoke API key
err := helper.RevokeAPIKey(db, keyID)
```

### Database Schema

```sql
CREATE TABLE api_keys (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(100),
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    last_used_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Security Notes

- Keys are hashed using SHA256 before storage
- Plain key only shown once during generation
- Automatic `last_used_at` tracking
- Can be revoked without deleting user account

---

## Encryption at Rest

Enkripsi data sensitif sebelum disimpan ke database.

### Setup

```go
// main.go
if err := helper.InitEncryption(config.ENV.ENCRYPTION_KEY); err != nil {
    log.Fatal(err)
}
```

### Environment Variable

```bash
# .env
ENCRYPTION_KEY="your-32-character-secret-key-here!!"
```

### Encrypt/Decrypt Data

```go
import "starter-gofiber/helper"

// Encrypt
encrypted, err := helper.Encrypt("sensitive data")

// Decrypt
decrypted, err := helper.Decrypt(encrypted)
```

### Database Field Encryption

```go
// Before saving
user.SSN, err = helper.EncryptField(user.SSN)
user.CreditCard, err = helper.EncryptField(user.CreditCard)
db.Save(&user)

// After loading
user.SSN, err = helper.DecryptField(user.SSN)
user.CreditCard, err = helper.DecryptField(user.CreditCard)
```

### Utility Functions

```go
// Hash string (one-way)
hash := helper.HashString("password123")

// Generate random string
random := helper.GenerateRandomString(32)

// Generate secure token
token := helper.GenerateSecureToken(16) // base64 URL-safe
```

### Algorithm

- **Cipher**: AES-256-GCM (Authenticated Encryption)
- **Key Derivation**: SHA256
- **Encoding**: Base64
- **Nonce**: 12 bytes random per encryption

---

## SQL Injection Prevention

### GORM Protection (Built-in)

```go
// ✅ Safe: Parameterized query
db.Where("email = ?", userInput).First(&user)

// ❌ Unsafe: String concatenation
db.Where("email = '" + userInput + "'").First(&user)
```

### Validation Helper

```go
if err := helper.ValidateNoSQLInjection(userInput); err != nil {
    return err
}
```

### Detected Patterns

- `UNION SELECT`
- `DROP TABLE`
- `--` (SQL comments)
- `/*` `*/` (block comments)
- `xp_cmdshell`
- Boolean operators in suspicious contexts

---

## XSS Protection

### Content Security Policy

Set via Security Headers middleware.

### Input Sanitization

```go
// Strip all HTML tags
clean := helper.SanitizeInput(userInput)

// Allow safe HTML for rich text
clean := helper.SanitizeHTML(userInput)
```

### Validation

```go
if err := helper.ValidateNoXSS(userInput); err != nil {
    return &helper.ValidationError{
        Message: "Dangerous content detected",
        Order:   "XS1",
    }
}
```

### Detected Patterns

- `<script>` tags
- `javascript:` protocol
- `onerror=`, `onload=`, `onclick=` handlers
- `<iframe>`, `<embed>`, `<object>` tags
- Data URIs with scripts

### Template Protection

```go
// Fiber automatically escapes HTML in templates
c.Render("template", fiber.Map{
    "UserInput": "<script>alert('xss')</script>",
})
// Rendered as: &lt;script&gt;alert('xss')&lt;/script&gt;
```

---

## Best Practices

### 1. Defense in Depth

Gunakan multiple layers of security:

```go
app.Use(middleware.HTTPSRedirectMiddleware())
app.Use(middleware.SecurityHeadersMiddleware())
app.Use(middleware.RateLimiter.Middleware())
app.Use(middleware.IPWhitelistMiddleware(trustedIPs))

// Then your routes
app.Post("/api/data", handler.CreateData)
```

### 2. Environment-Specific

```go
if config.ENV.ENV_TYPE == "production" {
    app.Use(middleware.ForceHTTPSMiddleware())
    app.Use(middleware.IPWhitelistMiddleware(productionIPs))
}
```

### 3. Validate All Inputs

```go
func CreateUser(c *fiber.Ctx) error {
    var req dto.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return err
    }

    // Sanitize
    req.Name = helper.SanitizeInput(req.Name)
    req.Bio = helper.SanitizeHTML(req.Bio)

    // Validate
    if err := helper.ValidateNoXSS(req.Name); err != nil {
        return err
    }

    // ... proceed
}
```

### 4. Encrypt Sensitive Data

```go
type User struct {
    ID        uint
    Email     string
    SSN       string `gorm:"type:text"` // Encrypted in DB
    Password  string `gorm:"type:varchar(255)"` // Hashed
}

// Before save
user.SSN, _ = helper.EncryptField(user.SSN)
user.Password = helper.HashPassword(user.Password)
```

### 5. Use API Keys for Service Accounts

```go
// For CI/CD, cron jobs, etc.
apiKey, _ := helper.GenerateAPIKey(db, serviceUserID, "GitHub Actions")

// Revoke when compromised
helper.RevokeAPIKey(db, keyID)
```

### 6. Monitor Rate Limits

```go
rateLimiter := middleware.NewUserRateLimiter(100, time.Minute)
app.Use(rateLimiter.Middleware())

// Monitor via metrics
// Rate limit hits will appear in application logs
```

---

## Testing

### Test HTTPS Redirect

```bash
curl -I http://localhost:3000/api/health
# Should return 301 redirect in production
```

### Test Rate Limiting

```bash
for i in {1..150}; do
  curl http://localhost:3000/api/posts
done
# After 100 requests: 429 Too Many Requests
```

### Test CSRF Protection

```bash
# Without token: 403 Forbidden
curl -X POST http://localhost:3000/api/posts

# With token: 200 OK
curl -X POST http://localhost:3000/api/posts \
  -H "X-CSRF-Token: your-token"
```

### Test API Key Auth

```bash
# Without key: 401 Unauthorized
curl http://localhost:3000/api/admin/users

# With key: 200 OK
curl -H "X-API-Key: your-api-key" \
     http://localhost:3000/api/admin/users
```

### Test IP Filtering

```bash
# From allowed IP: 200 OK
curl http://localhost:3000/api/data

# From blocked IP: 403 Forbidden
curl --interface 123.45.67.89 http://localhost:3000/api/data
```

---

## Security Checklist

- [x] Rate limiting per user/IP
- [x] CSRF protection for state-changing requests
- [x] Security headers (CSP, HSTS, X-Frame-Options, etc.)
- [x] HTTPS enforcement in production
- [x] Input sanitization (HTML, XSS, SQL injection)
- [x] IP whitelist/blacklist filtering
- [x] API key authentication for machine access
- [x] Encryption at rest for sensitive data
- [x] SQL injection prevention via GORM
- [x] XSS protection via CSP and sanitization

---

## Referensi

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [OWASP Cheat Sheet Series](https://cheatsheetseries.owasp.org/)
- [Fiber Security Middleware](https://docs.gofiber.io/api/middleware/)
- [Go Crypto Package](https://pkg.go.dev/crypto)
- [Content Security Policy](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)

---

## Troubleshooting

### CSRF Token Mismatch

**Problem**: 403 Forbidden pada POST request

**Solution**:
```javascript
// Pastikan token di-include dalam header
const token = getCookieValue('csrf_');
headers['X-CSRF-Token'] = token;
```

### Rate Limit False Positives

**Problem**: User legitimate ter-block

**Solution**:
```go
// Tingkatkan limit
rateLimiter := middleware.NewUserRateLimiter(200, time.Minute)
```

### Encryption Key Error

**Problem**: "Invalid encryption key length"

**Solution**:
```bash
# Generate 32-character key
openssl rand -base64 32 | cut -c1-32
```

### IP Filter Not Working

**Problem**: IP masih bisa akses setelah di-blacklist

**Solution**:
```go
// Pastikan middleware diinisialisasi sebelum routes
app.Use(middleware.IPBlacklistMiddleware(blockedIPs))
router.SetupRoutes(app) // After middleware
```

---

## Lisensi

MIT License - Silakan gunakan dengan bijak untuk melindungi aplikasi Anda.
