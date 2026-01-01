# Email Notification System

This application includes a comprehensive email notification system with template support, background job processing, and multi-provider SMTP compatibility.

## Features

- ✅ **Multiple SMTP Providers**: Mailpit (dev), Mailgun, SendGrid, AWS SES, Gmail, or any RFC-compliant SMTP server
- ✅ **Email Templates**: HTML + text alternatives for professional emails
- ✅ **Background Jobs**: Async email sending via Asynq with retries
- ✅ **TLS Support**: Configurable TLS/STARTTLS for secure connections
- ✅ **Template Engine**: Go html/template with data binding
- ✅ **Predefined Emails**: Welcome, password reset, email verification

## Configuration

### Environment Variables

Add these variables to your `.env` file:

```env
# SMTP Configuration
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM=noreply@example.com
SMTP_FROM_NAME=Your App Name
SMTP_USE_TLS=false

# Application URL (for email links)
APP_URL=http://localhost:3000
```

### Provider Examples

#### Mailpit (Development)

Mailpit is a local SMTP server perfect for testing. Run it with Docker:

```bash
docker run -d -p 1025:1025 -p 8025:8025 mailpit/mailpit
```

Then configure:

```env
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_USE_TLS=false
```

Access the web UI at: http://localhost:8025

#### Mailgun

```env
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USERNAME=postmaster@your-domain.mailgun.org
SMTP_PASSWORD=your-mailgun-smtp-password
SMTP_FROM=noreply@your-domain.com
SMTP_FROM_NAME=Your App
SMTP_USE_TLS=true
```

#### SendGrid

```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
SMTP_FROM=noreply@your-domain.com
SMTP_FROM_NAME=Your App
SMTP_USE_TLS=true
```

#### Gmail

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=your-email@gmail.com
SMTP_FROM_NAME=Your Name
SMTP_USE_TLS=true
```

**Note**: For Gmail, use an [App Password](https://support.google.com/accounts/answer/185833), not your regular password.

#### AWS SES

```env
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
SMTP_USERNAME=your-ses-smtp-username
SMTP_PASSWORD=your-ses-smtp-password
SMTP_FROM=noreply@verified-domain.com
SMTP_FROM_NAME=Your App
SMTP_USE_TLS=true
```

## Usage

### Sending Predefined Emails

#### Welcome Email

```go
import "starter-gofiber/jobs"

// Send welcome email (async via Asynq)
info, err := jobs.EnqueueEmailWelcome("user@example.com", "John Doe")
if err != nil {
    log.Printf("Failed to enqueue welcome email: %v", err)
}
log.Printf("Welcome email enqueued: ID=%s", info.ID)
```

#### Password Reset Email

```go
// Send password reset email with token
info, err := jobs.EnqueueEmailPasswordReset("user@example.com", "reset-token-123")
if err != nil {
    log.Printf("Failed to enqueue password reset: %v", err)
}
```

The email will include a link: `{APP_URL}/reset-password?token=reset-token-123`

#### Email Verification

```go
// Send email verification
info, err := jobs.EnqueueEmailVerification("user@example.com", "verify-token-456")
if err != nil {
    log.Printf("Failed to enqueue verification: %v", err)
}
```

The email will include a link: `{APP_URL}/verify-email?token=verify-token-456`

### Sending Custom Emails

#### Using EmailOptions

```go
import "starter-gofiber/helper"

err := helper.SendEmail(&helper.EmailOptions{
    To:       []string{"user@example.com"},
    Subject:  "Custom Email",
    HTMLBody: "<h1>Hello!</h1><p>This is a custom email.</p>",
    TextBody: "Hello!\n\nThis is a custom email.",
})
if err != nil {
    log.Printf("Failed to send email: %v", err)
}
```

#### With CC and BCC

```go
err := helper.SendEmail(&helper.EmailOptions{
    To:       []string{"user@example.com"},
    CC:       []string{"manager@example.com"},
    BCC:      []string{"admin@example.com"},
    Subject:  "Important Notice",
    HTMLBody: "<p>Important message...</p>",
    TextBody: "Important message...",
})
```

#### With Attachments

```go
err := helper.SendEmail(&helper.EmailOptions{
    To:          []string{"user@example.com"},
    Subject:     "Invoice #12345",
    HTMLBody:    "<p>Please find your invoice attached.</p>",
    TextBody:    "Please find your invoice attached.",
    Attachments: []string{"/path/to/invoice.pdf"},
})
```

### Using Custom Templates

#### Create Template Files

Create your template files in `templates/email/`:

**templates/email/custom-notification.html:**
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.Subject}}</title>
</head>
<body>
    <h1>Hello {{.UserName}}!</h1>
    <p>You have {{.Count}} new notifications.</p>
</body>
</html>
```

**templates/email/custom-notification.txt:**
```
Hello {{.UserName}}!

You have {{.Count}} new notifications.
```

#### Send Using Template

```go
import "starter-gofiber/helper"

template, err := helper.LoadEmailTemplate("custom-notification", map[string]interface{}{
    "Subject":  "New Notifications",
    "UserName": "John",
    "Count":    5,
})
if err != nil {
    log.Printf("Failed to load template: %v", err)
    return
}

err = helper.SendEmail(&helper.EmailOptions{
    To:       []string{"user@example.com"},
    Subject:  template.Subject,
    HTMLBody: template.HTMLBody,
    TextBody: template.TextBody,
})
```

### Background Jobs (Laravel-style Dispatch)

The email system integrates with Asynq for background processing:

```go
import "starter-gofiber/jobs"

// All emails are sent asynchronously with automatic retries

// 1. Welcome email
info, _ := jobs.EnqueueEmailWelcome("user@example.com", "John Doe")

// 2. Password reset
info, _ := jobs.EnqueueEmailPasswordReset("user@example.com", "token123")

// 3. Email verification
info, _ := jobs.EnqueueEmailVerification("user@example.com", "token456")

// 4. Custom email via job
info, _ := jobs.EnqueueEmailCustom(&jobs.EmailCustomPayload{
    To:       []string{"user@example.com"},
    Subject:  "Custom Message",
    HTMLBody: "<h1>Hello!</h1>",
    TextBody: "Hello!",
})
```

**Job Configuration:**
- Queue: `email`
- Max Retries: 3
- Exponential backoff on failures

## Email Templates

The system includes three built-in templates:

### 1. Welcome Email (`welcome`)

Sent when a user registers.

**Variables:**
- `{{.Name}}` - User's name
- `{{.Email}}` - User's email
- `{{.Subject}}` - Email subject

### 2. Password Reset (`reset-password`)

Sent when a user requests password reset.

**Variables:**
- `{{.Email}}` - User's email
- `{{.ResetURL}}` - Full reset URL with token
- `{{.Token}}` - Reset token (for manual URL construction)
- `{{.Subject}}` - Email subject

### 3. Email Verification (`verify-email`)

Sent for email address verification.

**Variables:**
- `{{.Email}}` - User's email
- `{{.VerifyURL}}` - Full verification URL with token
- `{{.Token}}` - Verification token
- `{{.Subject}}` - Email subject

## Template Customization

### Modify Existing Templates

Edit the HTML files in `templates/email/`:

```bash
templates/email/
├── welcome.html           # HTML version
├── welcome.txt            # Plain text version
├── reset-password.html
├── reset-password.txt
├── verify-email.html
└── verify-email.txt
```

**Important**: Always maintain both HTML and text versions for better email client compatibility.

### Create New Templates

1. Create HTML and text versions:
   ```bash
   touch templates/email/my-template.html
   touch templates/email/my-template.txt
   ```

2. Add your content with template variables:
   ```html
   <!-- my-template.html -->
   <!DOCTYPE html>
   <html>
   <body>
       <h1>Hello {{.Name}}!</h1>
       <p>{{.Message}}</p>
   </body>
   </html>
   ```

3. Load and use:
   ```go
   template, err := helper.LoadEmailTemplate("my-template", map[string]interface{}{
       "Name":    "John",
       "Message": "Welcome to our service!",
       "Subject": "Welcome",
   })
   ```

## Testing

### Using Mailpit

1. Start Mailpit:
   ```bash
   docker run -d -p 1025:1025 -p 8025:8025 mailpit/mailpit
   ```

2. Configure `.env`:
   ```env
   SMTP_HOST=localhost
   SMTP_PORT=1025
   SMTP_USERNAME=
   SMTP_PASSWORD=
   SMTP_USE_TLS=false
   ```

3. Send a test email:
   ```go
   jobs.EnqueueEmailWelcome("test@example.com", "Test User")
   ```

4. View in Mailpit UI: http://localhost:8025

### Testing Email Templates

Send test emails with sample data:

```go
// Test welcome email
helper.SendWelcomeEmail("test@example.com", "Test User")

// Test password reset
helper.SendPasswordResetEmail("test@example.com", "test-token-123")

// Test verification
helper.SendVerificationEmail("test@example.com", "verify-token-456")
```

## Error Handling

### Email Sending Errors

```go
err := helper.SendEmail(&helper.EmailOptions{
    To:      []string{"user@example.com"},
    Subject: "Test",
    HTMLBody: "<p>Test</p>",
})
if err != nil {
    // Log error - email will be retried by Asynq if sent via job
    log.Printf("Email send failed: %v", err)
    
    // Check specific errors
    if strings.Contains(err.Error(), "connection refused") {
        log.Println("SMTP server not available")
    } else if strings.Contains(err.Error(), "authentication failed") {
        log.Println("Invalid SMTP credentials")
    }
}
```

### Template Loading Errors

```go
template, err := helper.LoadEmailTemplate("my-template", data)
if err != nil {
    if os.IsNotExist(err) {
        log.Println("Template file not found")
    } else {
        log.Printf("Template parsing error: %v", err)
    }
}
```

## Troubleshooting

### Email Not Sending

1. **Check SMTP configuration:**
   ```bash
   # Verify environment variables are set
   echo $SMTP_HOST
   echo $SMTP_PORT
   ```

2. **Test SMTP connection:**
   ```bash
   # Using telnet
   telnet $SMTP_HOST $SMTP_PORT
   ```

3. **Check Asynq worker logs:**
   ```bash
   # Look for email task processing
   tail -f logs/app.log | grep email
   ```

### Template Not Found

Ensure template files exist:
```bash
ls -la templates/email/
```

Path should be relative to application root.

### TLS/SSL Errors

If you get TLS handshake errors:

1. Try disabling TLS for testing:
   ```env
   SMTP_USE_TLS=false
   ```

2. Check SMTP port:
   - Port 25: Usually no TLS
   - Port 587: STARTTLS
   - Port 465: SSL/TLS (not commonly supported)

3. Verify SMTP server supports TLS:
   ```bash
   openssl s_client -starttls smtp -connect smtp.example.com:587
   ```

### Authentication Failed

1. **Gmail**: Use App Password, not account password
2. **SendGrid**: Username must be `apikey`, password is your API key
3. **Mailgun**: Use full SMTP username from Mailgun dashboard
4. **AWS SES**: Use SMTP credentials, not IAM credentials

### Emails Go to Spam

1. Add SPF record:
   ```
   v=spf1 include:_spf.mailgun.org ~all
   ```

2. Add DKIM record (provided by your SMTP service)

3. Add DMARC record:
   ```
   v=DMARC1; p=none; rua=mailto:admin@yourdomain.com
   ```

4. Use verified sender domain

5. Avoid spam trigger words in subject/body

## Performance Tips

### Batch Sending

For bulk emails, use background jobs:

```go
for _, user := range users {
    jobs.EnqueueEmailWelcome(user.Email, user.Name)
}
```

Asynq will process them concurrently based on worker configuration.

### Rate Limiting

Configure Asynq worker concurrency in `helper/jobs.go`:

```go
srv := asynq.NewServer(
    asynq.RedisClientOpt{Addr: redisAddr},
    asynq.Config{
        Concurrency: 10,  // Process 10 emails concurrently
        Queues: map[string]int{
            "email":    6,  // Higher priority for emails
            "critical": 8,
            "default":  4,
        },
    },
)
```

### Template Caching

Templates are loaded fresh each time. For production, consider caching:

```go
var templateCache = make(map[string]*helper.EmailTemplate)

func GetCachedTemplate(name string, data map[string]interface{}) (*helper.EmailTemplate, error) {
    cacheKey := fmt.Sprintf("%s-%v", name, data)
    if cached, ok := templateCache[cacheKey]; ok {
        return cached, nil
    }
    
    template, err := helper.LoadEmailTemplate(name, data)
    if err == nil {
        templateCache[cacheKey] = template
    }
    return template, err
}
```

## Security Best Practices

1. **Never commit credentials**: Use environment variables
2. **Use App Passwords**: Don't use main account passwords
3. **Enable TLS**: Always use TLS in production
4. **Validate email addresses**: Before sending
5. **Rate limiting**: Prevent abuse
6. **Sanitize user input**: In templates to prevent XSS
7. **Use verified domains**: For better deliverability

## API Reference

### helper.SendEmail()

```go
func SendEmail(opts *EmailOptions) error
```

Sends an email with the given options.

**Parameters:**
- `opts.To` ([]string): Recipient email addresses
- `opts.CC` ([]string): CC recipients (optional)
- `opts.BCC` ([]string): BCC recipients (optional)
- `opts.Subject` (string): Email subject
- `opts.HTMLBody` (string): HTML email body
- `opts.TextBody` (string): Plain text body
- `opts.Attachments` ([]string): File paths to attach (optional)

### helper.LoadEmailTemplate()

```go
func LoadEmailTemplate(templateName string, data map[string]interface{}) (*EmailTemplate, error)
```

Loads and parses an email template.

**Parameters:**
- `templateName`: Name of template (without .html/.txt extension)
- `data`: Template variables

**Returns:**
- `EmailTemplate` with Subject, HTMLBody, TextBody

### jobs.EnqueueEmailWelcome()

```go
func EnqueueEmailWelcome(email, name string) (*asynq.TaskInfo, error)
```

Enqueues a welcome email background job.

### jobs.EnqueueEmailPasswordReset()

```go
func EnqueueEmailPasswordReset(email, resetToken string) (*asynq.TaskInfo, error)
```

Enqueues a password reset email.

### jobs.EnqueueEmailVerification()

```go
func EnqueueEmailVerification(email, verificationToken string) (*asynq.TaskInfo, error)
```

Enqueues an email verification job.

### jobs.EnqueueEmailCustom()

```go
func EnqueueEmailCustom(opts *EmailCustomPayload) (*asynq.TaskInfo, error)
```

Enqueues a custom email with full control over content.

## Examples

### Complete User Registration Flow

```go
func RegisterUser(c *fiber.Ctx) error {
    // ... validate and create user ...
    
    user := &entity.User{
        Email: "newuser@example.com",
        Name:  "John Doe",
    }
    
    // Save user to database
    if err := db.Create(user).Error; err != nil {
        return err
    }
    
    // Send welcome email asynchronously
    _, err := jobs.EnqueueEmailWelcome(user.Email, user.Name)
    if err != nil {
        log.Printf("Failed to send welcome email: %v", err)
        // Don't fail registration if email fails
    }
    
    return c.JSON(fiber.Map{
        "message": "Registration successful! Check your email.",
        "user":    user,
    })
}
```

### Password Reset Flow

```go
func RequestPasswordReset(c *fiber.Ctx) error {
    email := c.FormValue("email")
    
    // Find user
    var user entity.User
    if err := db.Where("email = ?", email).First(&user).Error; err != nil {
        // Return success even if user not found (security)
        return c.JSON(fiber.Map{"message": "If email exists, reset link sent"})
    }
    
    // Generate reset token
    token := helper.GenerateRandomString(32)
    user.ResetToken = token
    user.ResetTokenExpiry = time.Now().Add(1 * time.Hour)
    db.Save(&user)
    
    // Send reset email
    _, err := jobs.EnqueueEmailPasswordReset(user.Email, token)
    if err != nil {
        log.Printf("Failed to send reset email: %v", err)
    }
    
    return c.JSON(fiber.Map{
        "message": "Password reset email sent",
    })
}
```

### Email Verification Flow

```go
func SendVerification(user *entity.User) error {
    // Generate verification token
    token := helper.GenerateRandomString(32)
    user.VerificationToken = token
    user.VerificationTokenExpiry = time.Now().Add(24 * time.Hour)
    
    if err := db.Save(user).Error; err != nil {
        return err
    }
    
    // Send verification email
    _, err := jobs.EnqueueEmailVerification(user.Email, token)
    return err
}
```

## License

This email system is part of the starter-gofiber project.
