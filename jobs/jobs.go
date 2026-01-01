package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"starter-gofiber/helper"
)

// SendEmailJob dispatch email job to queue (Laravel style)
func SendEmailJob(to, subject, body string, data map[string]interface{}) error {
	payload := EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
		Data:    data,
	}

	return helper.EnqueueTaskToQueue(
		helper.TaskSendEmail,
		payload,
		helper.QueueDefault,
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

	return helper.EnqueueTaskToQueue(
		helper.TaskSendVerificationCode,
		payload,
		helper.QueueCritical, // High priority
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

	return helper.EnqueueTaskToQueue(
		helper.TaskSendPasswordReset,
		payload,
		helper.QueueCritical, // High priority
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

	return helper.EnqueueTaskToQueue(
		helper.TaskProcessExport,
		payload,
		helper.QueueLow, // Low priority - heavy task
	)
}

// CleanupOldFilesJob dispatch cleanup job
func CleanupOldFilesJob(directory string, olderThanDays int) error {
	payload := CleanupPayload{
		Directory: directory,
		OlderThan: olderThanDays,
	}

	return helper.EnqueueTask(
		helper.TaskCleanupOldFiles,
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

	return helper.EnqueueTaskToQueue(
		helper.TaskGenerateReport,
		payload,
		helper.QueueDefault,
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

	return helper.EnqueueTaskToQueue(
		helper.TaskSendNotification,
		payload,
		helper.QueueDefault,
	)
}

// SendDelayedEmailJob send email with delay (Laravel style: dispatch()->delay())
func SendDelayedEmailJob(to, subject, body string, delay time.Duration) error {
	payload := EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	return helper.EnqueueTaskWithDelay(
		helper.TaskSendEmail,
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

	return helper.EnqueueTaskAt(
		helper.TaskSendEmail,
		payload,
		processAt,
	)
}

// SendBulkEmails example for bulk email sending
func SendBulkEmails(recipients []string, subject, body string) error {
	for _, email := range recipients {
		// Dispatch each email as separate job
		if err := SendEmailJob(email, subject, body, nil); err != nil {
			helper.Logger.Error(fmt.Sprintf("Failed to queue email for %s: %v", email, err))
			continue
		}
	}

	helper.Logger.Info(fmt.Sprintf("Queued %d emails successfully", len(recipients)))
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
	return helper.EnqueueTask("chain:execute", string(data))
}
