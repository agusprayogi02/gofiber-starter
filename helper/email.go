package helper

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	mail "github.com/wneessen/go-mail"
	"go.uber.org/zap"
)

// EmailTemplate represents an email template
type EmailTemplate struct {
	Subject  string
	HTMLBody string
	TextBody string
}

// EmailOptions contains options for sending email
type EmailOptions struct {
	To           []string
	CC           []string
	BCC          []string
	Subject      string
	HTMLBody     string
	TextBody     string
	Attachments  []string
	TemplateDir  string
	TemplateName string
	TemplateData interface{}
}

// SendEmail sends an email using the configured SMTP settings
func SendEmail(opts *EmailOptions) error {
	// Create mail client
	client, err := Email.NewMailClient()
	if err != nil {
		Error("Failed to create mail client", zap.Error(err))
		return fmt.Errorf("failed to create mail client: %w", err)
	}
	defer client.Close()

	// Create new message
	msg := mail.NewMsg()

	// Set sender
	if err := msg.From(Email.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set from name if provided
	if Email.FromName != "" {
		if err := msg.FromFormat(Email.FromName, Email.From); err != nil {
			return fmt.Errorf("failed to set sender name: %w", err)
		}
	}

	// Set recipients
	if len(opts.To) > 0 {
		if err := msg.To(opts.To...); err != nil {
			return fmt.Errorf("failed to set recipients: %w", err)
		}
	}

	if len(opts.CC) > 0 {
		if err := msg.Cc(opts.CC...); err != nil {
			return fmt.Errorf("failed to set CC: %w", err)
		}
	}

	if len(opts.BCC) > 0 {
		if err := msg.Bcc(opts.BCC...); err != nil {
			return fmt.Errorf("failed to set BCC: %w", err)
		}
	}

	// Load template if specified
	if opts.TemplateName != "" {
		tmpl, err := LoadEmailTemplate(opts.TemplateDir, opts.TemplateName, opts.TemplateData)
		if err != nil {
			return fmt.Errorf("failed to load template: %w", err)
		}
		opts.Subject = tmpl.Subject
		opts.HTMLBody = tmpl.HTMLBody
		opts.TextBody = tmpl.TextBody
	}

	// Set subject
	msg.Subject(opts.Subject)

	// Set body
	if opts.HTMLBody != "" {
		msg.SetBodyString(mail.TypeTextHTML, opts.HTMLBody)
	}
	if opts.TextBody != "" {
		if opts.HTMLBody != "" {
			msg.AddAlternativeString(mail.TypeTextPlain, opts.TextBody)
		} else {
			msg.SetBodyString(mail.TypeTextPlain, opts.TextBody)
		}
	}

	// Add attachments
	for _, attachment := range opts.Attachments {
		msg.AttachFile(attachment)
	}

	// Send email
	if err := client.DialAndSend(msg); err != nil {
		Error("Failed to send email",
			zap.Error(err),
			zap.Strings("to", opts.To),
			zap.String("subject", opts.Subject),
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	Info("Email sent successfully",
		zap.Strings("to", opts.To),
		zap.String("subject", opts.Subject),
	)

	return nil
}

// LoadEmailTemplate loads and parses an email template
func LoadEmailTemplate(templateDir, templateName string, data interface{}) (*EmailTemplate, error) {
	if templateDir == "" {
		templateDir = "templates/email"
	}

	// Load HTML template
	htmlPath := filepath.Join(templateDir, templateName+".html")
	htmlTmpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err := htmlTmpl.Execute(&htmlBuf, data); err != nil {
		return nil, fmt.Errorf("failed to execute HTML template: %w", err)
	}

	// Load text template (optional)
	textPath := filepath.Join(templateDir, templateName+".txt")
	var textBody string
	textTmpl, err := template.ParseFiles(textPath)
	if err == nil {
		var textBuf bytes.Buffer
		if err := textTmpl.Execute(&textBuf, data); err == nil {
			textBody = textBuf.String()
		}
	}

	// Extract subject from template data if exists
	subject := templateName
	if dataMap, ok := data.(map[string]interface{}); ok {
		if subj, ok := dataMap["Subject"].(string); ok {
			subject = subj
		}
	}

	return &EmailTemplate{
		Subject:  subject,
		HTMLBody: htmlBuf.String(),
		TextBody: textBody,
	}, nil
}

// SendWelcomeEmail sends a welcome email to new users
func SendWelcomeEmail(email, name string) error {
	return SendEmail(&EmailOptions{
		To:           []string{email},
		TemplateName: "welcome",
		TemplateData: map[string]interface{}{
			"Subject": "Welcome to Our Platform!",
			"Name":    name,
			"Email":   email,
		},
	})
}

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(email, resetToken string) error {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:3000"
	}
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", appURL, resetToken)

	return SendEmail(&EmailOptions{
		To:           []string{email},
		TemplateName: "reset-password",
		TemplateData: map[string]interface{}{
			"Subject":  "Reset Your Password",
			"Email":    email,
			"ResetURL": resetURL,
			"Token":    resetToken,
		},
	})
}

// SendVerificationEmail sends an email verification email
func SendVerificationEmail(email, verificationToken string) error {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:3000"
	}
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", appURL, verificationToken)

	return SendEmail(&EmailOptions{
		To:           []string{email},
		TemplateName: "verify-email",
		TemplateData: map[string]interface{}{
			"Subject":   "Verify Your Email Address",
			"Email":     email,
			"VerifyURL": verifyURL,
			"Token":     verificationToken,
		},
	})
}
