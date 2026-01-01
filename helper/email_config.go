package helper

import (
	"fmt"
	"os"
	"strconv"

	mail "github.com/wneessen/go-mail"
)

// EmailConfig holds email configuration
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	UseTLS   bool
}

var Email *EmailConfig

// InitEmail initializes email configuration from environment variables
func InitEmail() error {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		port = 587 // Default SMTP port
	}

	useTLS := os.Getenv("SMTP_USE_TLS") == "true"

	Email = &EmailConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		FromName: os.Getenv("SMTP_FROM_NAME"),
		UseTLS:   useTLS,
	}

	// Validate required fields
	if Email.Host == "" {
		return fmt.Errorf("SMTP_HOST is required")
	}

	if Email.From == "" {
		return fmt.Errorf("SMTP_FROM is required")
	}

	return nil
}

// NewMailClient creates a new mail client with current configuration
func (e *EmailConfig) NewMailClient() (*mail.Client, error) {
	var opts []mail.Option

	// Set TLS or STARTTLS
	if e.UseTLS {
		opts = append(opts, mail.WithTLSPortPolicy(mail.TLSMandatory))
	} else {
		opts = append(opts, mail.WithTLSPortPolicy(mail.TLSOpportunistic))
	}

	// Set authentication if credentials are provided
	if e.Username != "" && e.Password != "" {
		opts = append(opts, mail.WithSMTPAuth(mail.SMTPAuthPlain))
		opts = append(opts, mail.WithUsername(e.Username))
		opts = append(opts, mail.WithPassword(e.Password))
	}

	// Set custom port if specified
	if e.Port > 0 {
		opts = append(opts, mail.WithPort(e.Port))
	}

	client, err := mail.NewClient(e.Host, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	return client, nil
}
