package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"starter-gofiber/internal/infrastructure/email"
	"starter-gofiber/pkg/logger"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// EmailPayload struktur data untuk email task
type EmailPayload struct {
	To      string                 `json:"to"`
	Subject string                 `json:"subject"`
	Body    string                 `json:"body"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// VerificationEmailPayload payload for verification email
type VerificationEmailPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
	URL   string `json:"url"`
}

// PasswordResetPayload payload for password reset email
type PasswordResetPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
	URL   string `json:"url"`
}

// ExportPayload payload for export job
type ExportPayload struct {
	UserID     uint   `json:"user_id"`
	ExportType string `json:"export_type"`
	Query      string `json:"query"`
	Filename   string `json:"filename"`
}

// CleanupPayload payload for cleanup job
type CleanupPayload struct {
	Directory string `json:"directory"`
	OlderThan int    `json:"older_than"`
}

// ReportPayload payload for report generation
type ReportPayload struct {
	ReportType string                 `json:"report_type"`
	StartDate  string                 `json:"start_date"`
	EndDate    string                 `json:"end_date"`
	Filters    map[string]interface{} `json:"filters"`
}

// NotificationPayload payload for notification job
type NotificationPayload struct {
	UserID  uint                   `json:"user_id"`
	Type    string                 `json:"type"`
	Title   string                 `json:"title"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// HandleSendEmail handler untuk task send email (legacy - use email.go handlers instead)
func HandleSendEmail(ctx context.Context, t *asynq.Task) error {
	var payload EmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		logger.Error("Failed to unmarshal email payload", zap.Error(err))
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	logger.Info("Sending email (legacy handler)",
		zap.String("to", payload.To),
		zap.String("subject", payload.Subject),
	)

	// Use email service from config
	err := email.SendEmail(&email.EmailOptions{
		To:       []string{payload.To},
		Subject:  payload.Subject,
		HTMLBody: payload.Body,
		TextBody: payload.Body,
	})
	if err != nil {
		logger.Error("Failed to send email", zap.Error(err))
		return err
	}

	logger.Info("Email sent successfully", zap.String("to", payload.To))
	return nil
}

// HandleSendVerificationEmail handler untuk verification email (legacy - use email.go handlers instead)
func HandleSendVerificationEmail(ctx context.Context, t *asynq.Task) error {
	var payload VerificationEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	logger.Info("Sending verification email (legacy handler)", zap.String("email", payload.Email))

	// Use email service from config
	return email.SendVerificationEmail(payload.Email, payload.Token)
}

// HandleSendPasswordReset handler untuk password reset email (legacy - use email.go handlers instead)
func HandleSendPasswordReset(ctx context.Context, t *asynq.Task) error {
	var payload PasswordResetPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	logger.Info("Sending password reset email (legacy handler)", zap.String("email", payload.Email))

	// Use email service from config
	return email.SendPasswordResetEmail(payload.Email, payload.Token)
}

// HandleProcessExport handler untuk export data
func HandleProcessExport(ctx context.Context, t *asynq.Task) error {
	var payload ExportPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	logger.Info(fmt.Sprintf("Processing export for user %d: %s", payload.UserID, payload.ExportType))

	// TODO: Implement actual export logic
	// 1. Fetch data from database
	// 2. Generate file (CSV/Excel/PDF)
	// 3. Upload to storage (S3/local)
	// 4. Send notification to user with download link

	logger.Info(fmt.Sprintf("Export completed: %s", payload.Filename))
	return nil
}

// HandleCleanupOldFiles handler untuk cleanup old files
func HandleCleanupOldFiles(ctx context.Context, t *asynq.Task) error {
	var payload CleanupPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	logger.Info(fmt.Sprintf("Cleaning up files in %s older than %d days", payload.Directory, payload.OlderThan))

	// TODO: Implement actual cleanup logic
	// 1. Scan directory
	// 2. Find files older than X days
	// 3. Delete files
	// 4. Log results

	logger.Info("Cleanup completed")
	return nil
}

// HandleGenerateReport handler untuk generate report
func HandleGenerateReport(ctx context.Context, t *asynq.Task) error {
	var payload ReportPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	logger.Info(fmt.Sprintf("Generating report: %s from %s to %s", payload.ReportType, payload.StartDate, payload.EndDate))

	// TODO: Implement report generation logic
	// 1. Fetch data based on filters
	// 2. Process data
	// 3. Generate report file
	// 4. Store report
	// 5. Notify user

	logger.Info("Report generated successfully")
	return nil
}

// HandleSendNotification handler untuk send notification
func HandleSendNotification(ctx context.Context, t *asynq.Task) error {
	var payload NotificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	logger.Info(fmt.Sprintf("Sending %s notification to user %d: %s", payload.Type, payload.UserID, payload.Title))

	// TODO: Implement notification logic based on type
	// switch payload.Type {
	// case "email":
	//     // Send email notification
	// case "push":
	//     // Send push notification (Firebase, OneSignal, etc)
	// case "sms":
	//     // Send SMS (Twilio, Nexmo, etc)
	// }

	logger.Info("Notification sent successfully")
	return nil
}
