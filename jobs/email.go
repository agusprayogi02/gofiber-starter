package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"starter-gofiber/helper"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

const (
	// Email task types
	TypeEmailWelcome       = "email:welcome"
	TypeEmailPasswordReset = "email:password_reset"
	TypeEmailVerification  = "email:verification"
	TypeEmailCustom        = "email:custom"
)

// Email job payloads
type EmailWelcomePayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type EmailPasswordResetPayload struct {
	Email      string `json:"email"`
	ResetToken string `json:"reset_token"`
}

type EmailVerificationPayload struct {
	Email             string `json:"email"`
	VerificationToken string `json:"verification_token"`
}

type EmailCustomPayload struct {
	To           []string               `json:"to"`
	CC           []string               `json:"cc,omitempty"`
	BCC          []string               `json:"bcc,omitempty"`
	Subject      string                 `json:"subject"`
	HTMLBody     string                 `json:"html_body,omitempty"`
	TextBody     string                 `json:"text_body,omitempty"`
	TemplateName string                 `json:"template_name,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
	Attachments  []string               `json:"attachments,omitempty"`
}

// EnqueueEmailWelcome enqueues a welcome email task
func EnqueueEmailWelcome(email, name string) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(EmailWelcomePayload{
		Email: email,
		Name:  name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeEmailWelcome, payload, asynq.Queue("email"), asynq.MaxRetry(3))
	return helper.AsynqClientInstance.Enqueue(task)
}

// EnqueueEmailPasswordReset enqueues a password reset email task
func EnqueueEmailPasswordReset(email, resetToken string) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(EmailPasswordResetPayload{
		Email:      email,
		ResetToken: resetToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeEmailPasswordReset, payload, asynq.Queue("email"), asynq.MaxRetry(3))
	return helper.AsynqClientInstance.Enqueue(task)
}

// EnqueueEmailVerification enqueues an email verification task
func EnqueueEmailVerification(email, verificationToken string) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(EmailVerificationPayload{
		Email:             email,
		VerificationToken: verificationToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeEmailVerification, payload, asynq.Queue("email"), asynq.MaxRetry(3))
	return helper.AsynqClientInstance.Enqueue(task)
}

// EnqueueEmailCustom enqueues a custom email task
func EnqueueEmailCustom(opts *EmailCustomPayload) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeEmailCustom, payload, asynq.Queue("email"), asynq.MaxRetry(3))
	return helper.AsynqClientInstance.Enqueue(task)
}

// HandleEmailWelcome handles welcome email tasks
func HandleEmailWelcome(ctx context.Context, t *asynq.Task) error {
	var payload EmailWelcomePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	helper.Info("Processing welcome email job",
		zap.String("email", payload.Email),
		zap.String("name", payload.Name),
	)

	if err := helper.SendWelcomeEmail(payload.Email, payload.Name); err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	helper.Info("Welcome email sent successfully",
		zap.String("email", payload.Email),
	)

	return nil
}

// HandleEmailPasswordReset handles password reset email tasks
func HandleEmailPasswordReset(ctx context.Context, t *asynq.Task) error {
	var payload EmailPasswordResetPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	helper.Info("Processing password reset email job",
		zap.String("email", payload.Email),
	)

	if err := helper.SendPasswordResetEmail(payload.Email, payload.ResetToken); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	helper.Info("Password reset email sent successfully",
		zap.String("email", payload.Email),
	)

	return nil
}

// HandleEmailVerification handles email verification tasks
func HandleEmailVerification(ctx context.Context, t *asynq.Task) error {
	var payload EmailVerificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	helper.Info("Processing email verification job",
		zap.String("email", payload.Email),
	)

	if err := helper.SendVerificationEmail(payload.Email, payload.VerificationToken); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	helper.Info("Verification email sent successfully",
		zap.String("email", payload.Email),
	)

	return nil
}

// HandleEmailCustom handles custom email tasks
func HandleEmailCustom(ctx context.Context, t *asynq.Task) error {
	var payload EmailCustomPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	helper.Info("Processing custom email job",
		zap.Strings("to", payload.To),
		zap.String("subject", payload.Subject),
	)

	opts := &helper.EmailOptions{
		To:           payload.To,
		CC:           payload.CC,
		BCC:          payload.BCC,
		Subject:      payload.Subject,
		HTMLBody:     payload.HTMLBody,
		TextBody:     payload.TextBody,
		TemplateName: payload.TemplateName,
		TemplateData: payload.TemplateData,
		Attachments:  payload.Attachments,
	}

	if err := helper.SendEmail(opts); err != nil {
		return fmt.Errorf("failed to send custom email: %w", err)
	}

	helper.Info("Custom email sent successfully",
		zap.Strings("to", payload.To),
	)

	return nil
}
