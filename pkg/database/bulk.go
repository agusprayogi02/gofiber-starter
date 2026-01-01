package database

import (
	"fmt"

	"starter-gofiber/pkg/apierror"

	"gorm.io/gorm"
)

// BulkCreateResult represents result of bulk create operation
type BulkCreateResult struct {
	SuccessCount int         `json:"success_count"`
	FailedCount  int         `json:"failed_count"`
	Errors       []BulkError `json:"errors,omitempty"`
	CreatedIDs   []uint      `json:"created_ids,omitempty"`
}

// BulkUpdateResult represents result of bulk update operation
type BulkUpdateResult struct {
	UpdatedCount int         `json:"updated_count"`
	FailedCount  int         `json:"failed_count"`
	Errors       []BulkError `json:"errors,omitempty"`
}

// BulkDeleteResult represents result of bulk delete operation
type BulkDeleteResult struct {
	DeletedCount int         `json:"deleted_count"`
	FailedCount  int         `json:"failed_count"`
	Errors       []BulkError `json:"errors,omitempty"`
}

// BulkError represents an error in bulk operation
type BulkError struct {
	Index   int    `json:"index"`
	ID      uint   `json:"id,omitempty"`
	Message string `json:"message"`
}

// BulkCreate creates multiple records in a single transaction
func BulkCreate(db *gorm.DB, records interface{}) (*BulkCreateResult, error) {
	result := &BulkCreateResult{
		CreatedIDs: make([]uint, 0),
		Errors:     make([]BulkError, 0),
	}

	// Use transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(records).Error; err != nil {
			return &apierror.InternalServerError{
				Message: fmt.Sprintf("Failed to bulk create records: %v", err),
				Order:   "H-Bulk-Create-1",
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Count success (this is simplified, you might need reflection for actual count)
	result.SuccessCount = 1 // Placeholder

	return result, nil
}

// BulkCreateWithValidation creates multiple records with individual validation
func BulkCreateWithValidation(db *gorm.DB, records []interface{}) *BulkCreateResult {
	result := &BulkCreateResult{
		CreatedIDs: make([]uint, 0),
		Errors:     make([]BulkError, 0),
	}

	for i, record := range records {
		err := db.Create(record).Error
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, BulkError{
				Index:   i,
				Message: err.Error(),
			})
		} else {
			result.SuccessCount++
			// Get ID using reflection if needed
			// result.CreatedIDs = append(result.CreatedIDs, id)
		}
	}

	return result
}

// BulkUpdate updates multiple records by IDs
func BulkUpdate(db *gorm.DB, model interface{}, ids []uint, updates map[string]interface{}) (*BulkUpdateResult, error) {
	result := &BulkUpdateResult{
		Errors: make([]BulkError, 0),
	}

	if len(ids) == 0 {
		return nil, &apierror.BadRequestError{
			Message: "No IDs provided for bulk update",
			Order:   "H-Bulk-Update-1",
		}
	}

	// Update all records matching IDs
	tx := db.Model(model).Where("id IN ?", ids).Updates(updates)
	if tx.Error != nil {
		return nil, &apierror.InternalServerError{
			Message: fmt.Sprintf("Failed to bulk update records: %v", tx.Error),
			Order:   "H-Bulk-Update-2",
		}
	}

	result.UpdatedCount = int(tx.RowsAffected)

	return result, nil
}

// BulkUpdateWithValidation updates multiple records with individual validation
func BulkUpdateWithValidation(db *gorm.DB, model interface{}, updates []map[string]interface{}) *BulkUpdateResult {
	result := &BulkUpdateResult{
		Errors: make([]BulkError, 0),
	}

	for i, update := range updates {
		id, ok := update["id"].(uint)
		if !ok {
			result.FailedCount++
			result.Errors = append(result.Errors, BulkError{
				Index:   i,
				Message: "Invalid or missing ID",
			})
			continue
		}

		delete(update, "id") // Remove ID from update map

		tx := db.Model(model).Where("id = ?", id).Updates(update)
		if tx.Error != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, BulkError{
				Index:   i,
				ID:      id,
				Message: tx.Error.Error(),
			})
		} else if tx.RowsAffected == 0 {
			result.FailedCount++
			result.Errors = append(result.Errors, BulkError{
				Index:   i,
				ID:      id,
				Message: "Record not found",
			})
		} else {
			result.UpdatedCount++
		}
	}

	return result
}

// BulkDelete deletes multiple records by IDs
func BulkDelete(db *gorm.DB, model interface{}, ids []uint) (*BulkDeleteResult, error) {
	result := &BulkDeleteResult{
		Errors: make([]BulkError, 0),
	}

	if len(ids) == 0 {
		return nil, &apierror.BadRequestError{
			Message: "No IDs provided for bulk delete",
			Order:   "H-Bulk-Delete-1",
		}
	}

	// Delete all records matching IDs (soft delete if model has DeletedAt)
	tx := db.Where("id IN ?", ids).Delete(model)
	if tx.Error != nil {
		return nil, &apierror.InternalServerError{
			Message: fmt.Sprintf("Failed to bulk delete records: %v", tx.Error),
			Order:   "H-Bulk-Delete-2",
		}
	}

	result.DeletedCount = int(tx.RowsAffected)

	return result, nil
}

// BulkDeletePermanent permanently deletes multiple records by IDs
func BulkDeletePermanent(db *gorm.DB, model interface{}, ids []uint) (*BulkDeleteResult, error) {
	result := &BulkDeleteResult{
		Errors: make([]BulkError, 0),
	}

	if len(ids) == 0 {
		return nil, &apierror.BadRequestError{
			Message: "No IDs provided for bulk delete",
			Order:   "H-Bulk-DeletePerm-1",
		}
	}

	// Permanent delete (Unscoped)
	tx := db.Unscoped().Where("id IN ?", ids).Delete(model)
	if tx.Error != nil {
		return nil, &apierror.InternalServerError{
			Message: fmt.Sprintf("Failed to permanently delete records: %v", tx.Error),
			Order:   "H-Bulk-DeletePerm-2",
		}
	}

	result.DeletedCount = int(tx.RowsAffected)

	return result, nil
}

// BulkRestore restores multiple soft-deleted records by IDs
func BulkRestore(db *gorm.DB, model interface{}, ids []uint) (*BulkUpdateResult, error) {
	result := &BulkUpdateResult{
		Errors: make([]BulkError, 0),
	}

	if len(ids) == 0 {
		return nil, &apierror.BadRequestError{
			Message: "No IDs provided for bulk restore",
			Order:   "H-Bulk-Restore-1",
		}
	}

	// Restore soft-deleted records
	tx := db.Model(model).Unscoped().Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Update("deleted_at", nil)

	if tx.Error != nil {
		return nil, &apierror.InternalServerError{
			Message: fmt.Sprintf("Failed to bulk restore records: %v", tx.Error),
			Order:   "H-Bulk-Restore-2",
		}
	}

	result.UpdatedCount = int(tx.RowsAffected)

	return result, nil
}

// BulkUpsert creates or updates multiple records (insert if not exists, update if exists)
func BulkUpsert(db *gorm.DB, records interface{}, conflictColumns []string) (*BulkCreateResult, error) {
	result := &BulkCreateResult{
		CreatedIDs: make([]uint, 0),
		Errors:     make([]BulkError, 0),
	}

	// Use Clauses for upsert (ON CONFLICT DO UPDATE)
	err := db.Transaction(func(tx *gorm.DB) error {
		// This requires GORM v2 with proper driver support
		if err := tx.Save(records).Error; err != nil {
			return &apierror.InternalServerError{
				Message: fmt.Sprintf("Failed to bulk upsert records: %v", err),
				Order:   "H-Bulk-Upsert-1",
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result.SuccessCount = 1 // Placeholder

	return result, nil
}

// ValidateBulkIDs validates that all IDs exist in database
func ValidateBulkIDs(db *gorm.DB, model interface{}, ids []uint) error {
	var count int64
	if err := db.Model(model).Where("id IN ?", ids).Count(&count).Error; err != nil {
		return &apierror.InternalServerError{
			Message: "Failed to validate IDs",
			Order:   "H-Bulk-Validate-1",
		}
	}

	if int(count) != len(ids) {
		return &apierror.NotFoundError{
			Message: fmt.Sprintf("Some IDs not found. Expected %d, found %d", len(ids), count),
			Order:   "H-Bulk-Validate-2",
		}
	}

	return nil
}
