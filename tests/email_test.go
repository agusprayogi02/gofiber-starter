package tests

import (
	"os"
	"testing"

	"starter-gofiber/helper"
	"starter-gofiber/jobs"

	"github.com/stretchr/testify/suite"
)

// EmailTestSuite tests the email notification system
// Covers:
// - Email configuration initialization
// - Email template loading (welcome, password-reset, verify-email)
// - Email payload structures for background jobs
// - Template data binding with variables
// - Full email workflow integration tests
//
// Note: These tests validate email system components without requiring
// an actual SMTP server. To test real email sending, set up Mailpit
// or another SMTP server and run integration tests.
type EmailTestSuite struct {
	suite.Suite
}

func TestEmailTestSuite(t *testing.T) {
	suite.Run(t, new(EmailTestSuite))
}

func (s *EmailTestSuite) SetupSuite() {
	// Set up test SMTP configuration
	os.Setenv("SMTP_HOST", "localhost")
	os.Setenv("SMTP_PORT", "1025")
	os.Setenv("SMTP_USERNAME", "")
	os.Setenv("SMTP_PASSWORD", "")
	os.Setenv("SMTP_FROM", "test@example.com")
	os.Setenv("SMTP_FROM_NAME", "Test App")
	os.Setenv("SMTP_USE_TLS", "false")
	os.Setenv("APP_URL", "http://localhost:3000")

	// Initialize email config
	helper.InitEmail()
}

func (s *EmailTestSuite) TearDownSuite() {
	// Clean up environment variables
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USERNAME")
	os.Unsetenv("SMTP_PASSWORD")
	os.Unsetenv("SMTP_FROM")
	os.Unsetenv("SMTP_FROM_NAME")
	os.Unsetenv("SMTP_USE_TLS")
	os.Unsetenv("APP_URL")
}

// Test Email Configuration
func (s *EmailTestSuite) TestEmailConfigInitialization() {
	s.NotNil(helper.Email, "Email config should be initialized")
	s.Equal("localhost", helper.Email.Host)
	s.Equal(1025, helper.Email.Port)
	s.Equal("test@example.com", helper.Email.From)
	s.Equal("Test App", helper.Email.FromName)
	s.False(helper.Email.UseTLS)
}

// Test Email Template Loading
func (s *EmailTestSuite) TestLoadEmailTemplate_Welcome() {
	data := map[string]interface{}{
		"Name":    "John Doe",
		"Email":   "john@example.com",
		"Subject": "Welcome to Our App",
	}

	template, err := helper.LoadEmailTemplate("../templates/email", "welcome", data)
	s.NoError(err, "Should load welcome template successfully")
	s.NotNil(template)
	s.Equal("Welcome to Our App", template.Subject)
	s.Contains(template.HTMLBody, "John Doe")
	s.Contains(template.HTMLBody, "john@example.com")
	s.Contains(template.TextBody, "John Doe")
	s.NotEmpty(template.TextBody)
}

func (s *EmailTestSuite) TestLoadEmailTemplate_PasswordReset() {
	data := map[string]interface{}{
		"Email":    "user@example.com",
		"ResetURL": "http://localhost:3000/reset-password?token=abc123",
		"Token":    "abc123",
		"Subject":  "Reset Your Password",
	}

	template, err := helper.LoadEmailTemplate("../templates/email", "reset-password", data)
	s.NoError(err, "Should load reset-password template successfully")
	s.NotNil(template)
	s.Equal("Reset Your Password", template.Subject)
	s.Contains(template.HTMLBody, "http://localhost:3000/reset-password?token=abc123")
	s.Contains(template.HTMLBody, "user@example.com")
	s.Contains(template.TextBody, "reset your password")
	s.NotEmpty(template.TextBody)
}

func (s *EmailTestSuite) TestLoadEmailTemplate_EmailVerification() {
	data := map[string]interface{}{
		"Email":     "user@example.com",
		"VerifyURL": "http://localhost:3000/verify-email?token=xyz789",
		"Token":     "xyz789",
		"Subject":   "Verify Your Email",
	}

	template, err := helper.LoadEmailTemplate("../templates/email", "verify-email", data)
	s.NoError(err, "Should load verify-email template successfully")
	s.NotNil(template)
	s.Equal("Verify Your Email", template.Subject)
	s.Contains(template.HTMLBody, "http://localhost:3000/verify-email?token=xyz789")
	s.Contains(template.HTMLBody, "user@example.com")
	s.Contains(template.TextBody, "verify your email")
	s.NotEmpty(template.TextBody)
}

func (s *EmailTestSuite) TestLoadEmailTemplate_NonExistent() {
	data := map[string]interface{}{
		"Subject": "Test",
	}

	template, err := helper.LoadEmailTemplate("../templates/email", "non-existent-template", data)
	s.Error(err, "Should return error for non-existent template")
	s.Nil(template)
}

// Test Email Task Type Constants
func (s *EmailTestSuite) TestEmailTaskTypeConstants() {
	s.Equal("email:welcome", jobs.TypeEmailWelcome)
	s.Equal("email:password_reset", jobs.TypeEmailPasswordReset)
	s.Equal("email:verification", jobs.TypeEmailVerification)
	s.Equal("email:custom", jobs.TypeEmailCustom)
}

// Test Email Payload Structures
func (s *EmailTestSuite) TestEmailWelcomePayload() {
	payload := jobs.EmailWelcomePayload{
		Email: "test@example.com",
		Name:  "Test User",
	}
	s.Equal("test@example.com", payload.Email)
	s.Equal("Test User", payload.Name)
}

func (s *EmailTestSuite) TestEmailPasswordResetPayload() {
	payload := jobs.EmailPasswordResetPayload{
		Email:      "test@example.com",
		ResetToken: "reset-token-123",
	}
	s.Equal("test@example.com", payload.Email)
	s.Equal("reset-token-123", payload.ResetToken)
}

func (s *EmailTestSuite) TestEmailVerificationPayload() {
	payload := jobs.EmailVerificationPayload{
		Email:             "test@example.com",
		VerificationToken: "verify-token-456",
	}
	s.Equal("test@example.com", payload.Email)
	s.Equal("verify-token-456", payload.VerificationToken)
}

func (s *EmailTestSuite) TestEmailCustomPayload() {
	payload := &jobs.EmailCustomPayload{
		To:       []string{"recipient@example.com"},
		Subject:  "Custom Email",
		HTMLBody: "<h1>Custom</h1>",
		TextBody: "Custom",
	}

	s.Len(payload.To, 1)
	s.Equal("Custom Email", payload.Subject)
	s.Contains(payload.HTMLBody, "<h1>Custom</h1>")
	s.NotEmpty(payload.TextBody)
}

// Test Email Client Initialization
func (s *EmailTestSuite) TestEmailClientInitialized() {
	// Email client is created internally when sending emails
	// We test that Email config is properly initialized
	s.NotNil(helper.Email)
	s.NotEmpty(helper.Email.Host)
	s.NotEmpty(helper.Email.From)
}

// Test Template Data Binding
func (s *EmailTestSuite) TestTemplateDataBinding_MultipleVariables() {
	data := map[string]interface{}{
		"Name":      "John Doe",
		"Email":     "john@example.com",
		"Subject":   "Test Subject",
		"CustomVar": "Custom Value",
	}

	template, err := helper.LoadEmailTemplate("../templates/email", "welcome", data)
	s.NoError(err)
	s.Contains(template.HTMLBody, "John Doe")
	s.Contains(template.HTMLBody, "john@example.com")
}

// Integration Test: Full Email Flow
func (s *EmailTestSuite) TestFullEmailFlow_WelcomeEmail() {
	// 1. Load template
	data := map[string]interface{}{
		"Name":    "Integration Test User",
		"Email":   "integration@example.com",
		"Subject": "Welcome!",
	}

	template, err := helper.LoadEmailTemplate("../templates/email", "welcome", data)
	s.NoError(err)
	s.NotNil(template)

	// 2. Prepare email options
	opts := &helper.EmailOptions{
		To:       []string{"integration@example.com"},
		Subject:  template.Subject,
		HTMLBody: template.HTMLBody,
		TextBody: template.TextBody,
	}

	// 3. Validate options
	s.NotEmpty(opts.To)
	s.NotEmpty(opts.Subject)
	s.NotEmpty(opts.HTMLBody)
	s.NotEmpty(opts.TextBody)
}

func (s *EmailTestSuite) TestFullEmailFlow_PasswordReset() {
	// 1. Generate reset token
	token := "generated-reset-token-123"

	// 2. Load template
	data := map[string]interface{}{
		"Email":    "user@example.com",
		"ResetURL": "http://localhost:3000/reset-password?token=" + token,
		"Token":    token,
		"Subject":  "Reset Your Password",
	}

	template, err := helper.LoadEmailTemplate("../templates/email", "reset-password", data)
	s.NoError(err)

	// 3. Validate template contains token
	s.Contains(template.HTMLBody, token)
	s.Contains(template.TextBody, token)
}
