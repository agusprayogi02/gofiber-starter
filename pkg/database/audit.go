package database

import (
	"encoding/json"
	"fmt"
	"reflect"
	"starter-gofiber/entity"
	"time"

	"gorm.io/gorm"
)

// AuditLogger handles automatic audit logging for GORM operations
type AuditLogger struct {
	db        *gorm.DB
	userID    *uint
	username  string
	ipAddress string
	userAgent string
	requestID string
}

// NewAuditLogger creates a new audit logger instance
func NewAuditLogger(db *gorm.DB) *AuditLogger {
	return &AuditLogger{
		db: db,
	}
}

// WithUser sets the user information for audit log
func (a *AuditLogger) WithUser(userID uint, username string) *AuditLogger {
	a.userID = &userID
	a.username = username
	return a
}

// WithRequest sets the request context for audit log
func (a *AuditLogger) WithRequest(ipAddress, userAgent, requestID string) *AuditLogger {
	a.ipAddress = ipAddress
	a.userAgent = userAgent
	a.requestID = requestID
	return a
}

// LogCreate logs a CREATE operation
func (a *AuditLogger) LogCreate(entityType string, entityID uint, newData interface{}) error {
	newValues, err := json.Marshal(newData)
	if err != nil {
		return err
	}

	auditLog := entity.AuditLog{
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      entity.AuditActionCreate,
		Description: fmt.Sprintf("Created %s #%d", entityType, entityID),
		NewValues:   string(newValues),
		UserID:      a.userID,
		Username:    a.username,
		IPAddress:   a.ipAddress,
		UserAgent:   a.userAgent,
		RequestID:   a.requestID,
	}

	return a.db.Create(&auditLog).Error
}

// LogUpdate logs an UPDATE operation
func (a *AuditLogger) LogUpdate(entityType string, entityID uint, oldData, newData interface{}) error {
	oldValues, err := json.Marshal(oldData)
	if err != nil {
		return err
	}

	newValues, err := json.Marshal(newData)
	if err != nil {
		return err
	}

	// Get changed fields
	changes := getChangedFields(oldData, newData)
	description := fmt.Sprintf("Updated %s #%d", entityType, entityID)
	if len(changes) > 0 {
		description = fmt.Sprintf("Updated %s #%d: %v", entityType, entityID, changes)
	}

	auditLog := entity.AuditLog{
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      entity.AuditActionUpdate,
		Description: description,
		OldValues:   string(oldValues),
		NewValues:   string(newValues),
		UserID:      a.userID,
		Username:    a.username,
		IPAddress:   a.ipAddress,
		UserAgent:   a.userAgent,
		RequestID:   a.requestID,
	}

	return a.db.Create(&auditLog).Error
}

// LogDelete logs a DELETE operation (soft or hard delete)
func (a *AuditLogger) LogDelete(entityType string, entityID uint, oldData interface{}, isSoftDelete bool) error {
	oldValues, err := json.Marshal(oldData)
	if err != nil {
		return err
	}

	action := entity.AuditActionDelete
	description := fmt.Sprintf("Deleted %s #%d", entityType, entityID)
	if isSoftDelete {
		description = fmt.Sprintf("Soft deleted %s #%d", entityType, entityID)
	}

	auditLog := entity.AuditLog{
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      action,
		Description: description,
		OldValues:   string(oldValues),
		UserID:      a.userID,
		Username:    a.username,
		IPAddress:   a.ipAddress,
		UserAgent:   a.userAgent,
		RequestID:   a.requestID,
	}

	return a.db.Create(&auditLog).Error
}

// LogRestore logs a RESTORE operation
func (a *AuditLogger) LogRestore(entityType string, entityID uint, data interface{}) error {
	newValues, err := json.Marshal(data)
	if err != nil {
		return err
	}

	auditLog := entity.AuditLog{
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      entity.AuditActionRestore,
		Description: fmt.Sprintf("Restored %s #%d", entityType, entityID),
		NewValues:   string(newValues),
		UserID:      a.userID,
		Username:    a.username,
		IPAddress:   a.ipAddress,
		UserAgent:   a.userAgent,
		RequestID:   a.requestID,
	}

	return a.db.Create(&auditLog).Error
}

// getChangedFields compares two structs and returns changed field names
func getChangedFields(oldData, newData interface{}) []string {
	var changes []string

	oldVal := reflect.ValueOf(oldData)
	newVal := reflect.ValueOf(newData)

	// Dereference pointers
	if oldVal.Kind() == reflect.Ptr {
		oldVal = oldVal.Elem()
	}
	if newVal.Kind() == reflect.Ptr {
		newVal = newVal.Elem()
	}

	// Only compare structs
	if oldVal.Kind() != reflect.Struct || newVal.Kind() != reflect.Struct {
		return changes
	}

	for i := 0; i < oldVal.NumField(); i++ {
		fieldName := oldVal.Type().Field(i).Name

		// Skip unexported fields and certain fields
		if !oldVal.Field(i).CanInterface() || fieldName == "UpdatedAt" || fieldName == "CreatedAt" {
			continue
		}

		oldFieldVal := oldVal.Field(i).Interface()
		newFieldVal := newVal.Field(i).Interface()

		if !reflect.DeepEqual(oldFieldVal, newFieldVal) {
			changes = append(changes, fieldName)
		}
	}

	return changes
}

// RegisterAuditCallbacks registers GORM callbacks for automatic audit logging
func RegisterAuditCallbacks(db *gorm.DB) {
	// After Create callback
	db.Callback().Create().After("gorm:create").Register("audit:create", func(db *gorm.DB) {
		if db.Statement.Schema != nil && db.Error == nil {
			// Skip audit logs to prevent recursion
			if db.Statement.Schema.Table == "audit_logs" {
				return
			}

			// Get created record ID
			if db.Statement.Schema.PrioritizedPrimaryField != nil {
				entityID := getEntityID(db.Statement.ReflectValue)
				if entityID > 0 {
					logger := NewAuditLogger(db.Session(&gorm.Session{NewDB: true}))
					// Get user context from statement context if available
					if userID, ok := db.Statement.Context.Value("user_id").(uint); ok {
						username, _ := db.Statement.Context.Value("username").(string)
						logger.WithUser(userID, username)
					}
					if ipAddress, ok := db.Statement.Context.Value("ip_address").(string); ok {
						userAgent, _ := db.Statement.Context.Value("user_agent").(string)
						requestID, _ := db.Statement.Context.Value("request_id").(string)
						logger.WithRequest(ipAddress, userAgent, requestID)
					}

					_ = logger.LogCreate(db.Statement.Schema.Table, entityID, db.Statement.Dest)
				}
			}
		}
	})

	// After Update callback
	db.Callback().Update().After("gorm:update").Register("audit:update", func(db *gorm.DB) {
		if db.Statement.Schema != nil && db.Error == nil {
			if db.Statement.Schema.Table == "audit_logs" {
				return
			}

			// Note: To track old values, you need to query before update
			// This is a simplified version that only logs the action
			if db.Statement.Schema.PrioritizedPrimaryField != nil {
				entityID := getEntityID(db.Statement.ReflectValue)
				if entityID > 0 {
					logger := NewAuditLogger(db.Session(&gorm.Session{NewDB: true}))
					if userID, ok := db.Statement.Context.Value("user_id").(uint); ok {
						username, _ := db.Statement.Context.Value("username").(string)
						logger.WithUser(userID, username)
					}
					if ipAddress, ok := db.Statement.Context.Value("ip_address").(string); ok {
						userAgent, _ := db.Statement.Context.Value("user_agent").(string)
						requestID, _ := db.Statement.Context.Value("request_id").(string)
						logger.WithRequest(ipAddress, userAgent, requestID)
					}

					// Log with minimal info (oldData would need to be fetched before update)
					_ = logger.LogUpdate(db.Statement.Schema.Table, entityID, nil, db.Statement.Dest)
				}
			}
		}
	})

	// After Delete callback
	db.Callback().Delete().After("gorm:delete").Register("audit:delete", func(db *gorm.DB) {
		if db.Statement.Schema != nil && db.Error == nil {
			if db.Statement.Schema.Table == "audit_logs" {
				return
			}

			if db.Statement.Schema.PrioritizedPrimaryField != nil {
				entityID := getEntityID(db.Statement.ReflectValue)
				if entityID > 0 {
					logger := NewAuditLogger(db.Session(&gorm.Session{NewDB: true}))
					if userID, ok := db.Statement.Context.Value("user_id").(uint); ok {
						username, _ := db.Statement.Context.Value("username").(string)
						logger.WithUser(userID, username)
					}
					if ipAddress, ok := db.Statement.Context.Value("ip_address").(string); ok {
						userAgent, _ := db.Statement.Context.Value("user_agent").(string)
						requestID, _ := db.Statement.Context.Value("request_id").(string)
						logger.WithRequest(ipAddress, userAgent, requestID)
					}

					// Check if it's a soft delete
					isSoftDelete := db.Statement.Unscoped == false
					_ = logger.LogDelete(db.Statement.Schema.Table, entityID, db.Statement.Dest, isSoftDelete)
				}
			}
		}
	})
}

// getEntityID extracts the ID from a reflect.Value
func getEntityID(value reflect.Value) uint {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return 0
	}

	idField := value.FieldByName("ID")
	if !idField.IsValid() {
		return 0
	}

	if idField.CanUint() {
		return uint(idField.Uint())
	}

	return 0
}

// GetAuditLogs retrieves audit logs with filters
func GetAuditLogs(db *gorm.DB, filter entity.AuditLogFilter, page, limit int) ([]entity.AuditLog, int64, error) {
	var logs []entity.AuditLog
	var total int64

	query := db.Model(&entity.AuditLog{})

	// Apply filters
	if filter.EntityType != "" {
		query = query.Where("entity_type = ?", filter.EntityType)
	}
	if filter.EntityID != nil {
		query = query.Where("entity_id = ?", *filter.EntityID)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	if filter.RequestID != "" {
		query = query.Where("request_id = ?", filter.RequestID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}

// GetEntityAuditHistory retrieves full history for a specific entity
func GetEntityAuditHistory(db *gorm.DB, entityType string, entityID uint) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at ASC").
		Find(&logs).Error
	return logs, err
}

// GetUserActivity retrieves all activities by a user
func GetUserActivity(db *gorm.DB, userID uint, startDate, endDate *time.Time) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	query := db.Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// CleanupOldAuditLogs removes audit logs older than specified days
func CleanupOldAuditLogs(db *gorm.DB, daysToKeep int) error {
	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)
	return db.Where("created_at < ?", cutoffDate).Delete(&entity.AuditLog{}).Error
}
