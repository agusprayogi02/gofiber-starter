package database

import (
	"gorm.io/gorm"
)

// SoftDeleteScope returns a scope that includes soft deleted records
func SoftDeleteScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}
}

// WithTrashed includes soft deleted records in query
func WithTrashed(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}

// OnlyTrashed returns only soft deleted records
func OnlyTrashed(db *gorm.DB) *gorm.DB {
	return db.Unscoped().Where("deleted_at IS NOT NULL")
}

// Restore restores a soft deleted record
func Restore(db *gorm.DB, model interface{}) error {
	return db.Unscoped().Model(model).Update("deleted_at", nil).Error
}

// ForceDelete permanently deletes a record (bypasses soft delete)
func ForceDelete(db *gorm.DB, model interface{}) error {
	return db.Unscoped().Delete(model).Error
}

// RestoreByID restores a soft deleted record by ID
func RestoreByID(db *gorm.DB, model interface{}, id interface{}) error {
	return db.Unscoped().Model(model).Where("id = ?", id).Update("deleted_at", nil).Error
}

// ForceDeleteByID permanently deletes a record by ID
func ForceDeleteByID(db *gorm.DB, model interface{}, id interface{}) error {
	return db.Unscoped().Where("id = ?", id).Delete(model).Error
}

// RestoreAll restores all soft deleted records for a model
func RestoreAll(db *gorm.DB, model interface{}) error {
	return db.Unscoped().Model(model).Where("deleted_at IS NOT NULL").Update("deleted_at", nil).Error
}

// CountTrashed counts soft deleted records
func CountTrashed(db *gorm.DB, model interface{}) (int64, error) {
	var count int64
	err := db.Unscoped().Model(model).Where("deleted_at IS NOT NULL").Count(&count).Error
	return count, err
}

// IsTrashed checks if a record is soft deleted
func IsTrashed(db *gorm.DB, model interface{}, id interface{}) (bool, error) {
	var count int64
	err := db.Unscoped().Model(model).Where("id = ? AND deleted_at IS NOT NULL", id).Count(&count).Error
	return count > 0, err
}
