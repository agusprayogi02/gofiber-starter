package entity

import (
	"time"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	AuditActionCreate AuditAction = "CREATE"
	AuditActionUpdate AuditAction = "UPDATE"
	AuditActionDelete AuditAction = "DELETE"
	AuditActionRestore AuditAction = "RESTORE"
)

// AuditLog tracks all data changes in the system
type AuditLog struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// Entity information
	EntityType string `gorm:"type:varchar(100);not null;index" json:"entity_type"` // e.g., "users", "posts"
	EntityID   uint   `gorm:"not null;index" json:"entity_id"`                     // ID of the affected record

	// Action details
	Action      AuditAction `gorm:"type:varchar(20);not null;index" json:"action"` // CREATE, UPDATE, DELETE, RESTORE
	Description string      `gorm:"type:text" json:"description"`                  // Human-readable description

	// Change tracking
	OldValues string `gorm:"type:json" json:"old_values,omitempty"` // JSON of old values (for UPDATE/DELETE)
	NewValues string `gorm:"type:json" json:"new_values,omitempty"` // JSON of new values (for CREATE/UPDATE)

	// User tracking
	UserID   *uint  `gorm:"index" json:"user_id,omitempty"`         // User who performed the action (nullable for system actions)
	Username string `gorm:"type:varchar(255)" json:"username"`      // Cached username for quick display
	IPAddress string `gorm:"type:varchar(45)" json:"ip_address"`    // IPv4 or IPv6
	UserAgent string `gorm:"type:text" json:"user_agent,omitempty"` // Browser/client info

	// Metadata
	RequestID string    `gorm:"type:varchar(100);index" json:"request_id,omitempty"` // Trace request chain
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName specifies the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogFilter for querying audit logs
type AuditLogFilter struct {
	EntityType string
	EntityID   *uint
	UserID     *uint
	Action     AuditAction
	StartDate  *time.Time
	EndDate    *time.Time
	RequestID  string
}
