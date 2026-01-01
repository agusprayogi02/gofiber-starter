package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"starter-gofiber/pkg/logger"
)

// SendEmailJob dispatch email job to queue (Laravel style)
func SendEmailJob(to, subject, body string, data map[string]interface{}) error {
	payload := EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
		Data:    data,
	}

	return EnqueueTaskToQueue(
		TaskSendEmail,
		payload,
		QueueDefault,
	)
}

// SendVerificationEmailJob dispatch verification email
func SendVerificationEmailJob(email, name, token, baseURL string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", baseURL, token)

	payload := VerificationEmailPayload{
		Email: email,
		Name:  name,
		Token: token,
		URL:   verifyURL,
	}

	return EnqueueTaskToQueue(
		TaskSendVerificationCode,
		payload,
		QueueCritical, // High priority
	)
}

// SendPasswordResetEmailJob dispatch password reset email
func SendPasswordResetEmailJob(email, name, token, baseURL string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)

	payload := PasswordResetPayload{
		Email: email,
		Name:  name,
		Token: token,
		URL:   resetURL,
	}

	return EnqueueTaskToQueue(
		TaskSendPasswordReset,
		payload,
		QueueCritical, // High priority
	)
}

// ProcessExportJob dispatch export job (async)
func ProcessExportJob(userID uint, exportType, query, filename string) error {
	payload := ExportPayload{
		UserID:     userID,
		ExportType: exportType,
		Query:      query,
		Filename:   filename,
	}

	return EnqueueTaskToQueue(
		TaskProcessExport,
		payload,
		QueueLow, // Low priority - heavy task
	)
}

// CleanupOldFilesJob dispatch cleanup job
func CleanupOldFilesJob(directory string, olderThanDays int) error {
	payload := CleanupPayload{
		Directory: directory,
		OlderThan: olderThanDays,
	}

	return EnqueueTask(
		TaskCleanupOldFiles,
		payload,
	)
}

// GenerateReportJob dispatch report generation job
func GenerateReportJob(reportType, startDate, endDate string, filters map[string]interface{}) error {
	payload := ReportPayload{
		ReportType: reportType,
		StartDate:  startDate,
		EndDate:    endDate,
		Filters:    filters,
	}

	return EnqueueTaskToQueue(
		TaskGenerateReport,
		payload,
		QueueDefault,
	)
}

// SendNotificationJob dispatch notification job
func SendNotificationJob(userID uint, notifType, title, message string, data map[string]interface{}) error {
	payload := NotificationPayload{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
		Data:    data,
	}

	return EnqueueTaskToQueue(
		TaskSendNotification,
		payload,
		QueueDefault,
	)
}

// SendDelayedEmailJob send email with delay (Laravel style: dispatch()->delay())
func SendDelayedEmailJob(to, subject, body string, delay time.Duration) error {
	payload := EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	return EnqueueTaskWithDelay(
		TaskSendEmail,
		payload,
		delay,
	)
}

// SendEmailAtJob send email at specific time
func SendEmailAtJob(to, subject, body string, processAt time.Time) error {
	payload := EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	return EnqueueTaskAt(
		TaskSendEmail,
		payload,
		processAt,
	)
}

// === New Email Job Functions (Using new email service) ===

// SendWelcomeEmailJob dispatch welcome email using new email service
func SendWelcomeEmailJob(email, name string) error {
	_, err := EnqueueEmailWelcome(email, name)
	return err
}

// SendPasswordResetJob dispatch password reset email using new email service
func SendPasswordResetJob(email, resetToken string) error {
	_, err := EnqueueEmailPasswordReset(email, resetToken)
	return err
}

// SendVerificationJob dispatch verification email using new email service
func SendVerificationJob(email, verificationToken string) error {
	_, err := EnqueueEmailVerification(email, verificationToken)
	return err
}

// SendCustomEmailJob dispatch custom email using new email service
func SendCustomEmailJob(to []string, subject, htmlBody, textBody string) error {
	_, err := EnqueueEmailCustom(&EmailCustomPayload{
		To:       to,
		Subject:  subject,
		HTMLBody: htmlBody,
		TextBody: textBody,
	})
	return err
}

// SendTemplatedEmailJob dispatch templated email
func SendTemplatedEmailJob(to []string, templateName string, data map[string]interface{}) error {
	_, err := EnqueueEmailCustom(&EmailCustomPayload{
		To:           to,
		TemplateName: templateName,
		TemplateData: data,
	})
	return err
}

// SendBulkEmails example for bulk email sending
func SendBulkEmails(recipients []string, subject, body string) error {
	for _, email := range recipients {
		// Dispatch each email as separate job
		if err := SendEmailJob(email, subject, body, nil); err != nil {
			logger.Error(fmt.Sprintf("Failed to queue email for %s", email))
			continue
		}
	}

	logger.Info(fmt.Sprintf("Queued %d emails successfully", len(recipients)))
	return nil
}

// ChainJobs example for job chaining (sequential execution)
func ChainJobs(jobs []map[string]interface{}) error {
	// Convert jobs to JSON for payload
	data, err := json.Marshal(jobs)
	if err != nil {
		return err
	}

	// Create a custom chain task
	// TODO: Implement custom handler for job chains
	return EnqueueTask("chain:execute", string(data))
}
