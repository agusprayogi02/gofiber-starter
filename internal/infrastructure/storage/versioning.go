package storage

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"starter-gofiber/entity"
	"starter-gofiber/pkg/apierror"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileVersionConfig represents file versioning configuration
type FileVersionConfig struct {
	MaxVersions    int    // Maximum number of versions to keep (0 = unlimited)
	ChecksumType   string // md5 or sha256
	AutoCleanup    bool   // Automatically cleanup old versions
	StorageType    string // local or s3
	EnableMetadata bool   // Store additional metadata
}

// FileMetadata represents additional file metadata
type FileMetadata struct {
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Duration  int    `json:"duration,omitempty"`   // For videos
	PageCount int    `json:"page_count,omitempty"` // For PDFs
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Comments  string `json:"comments,omitempty"`
}

// DefaultFileVersionConfig returns default versioning configuration
func DefaultFileVersionConfig() FileVersionConfig {
	return FileVersionConfig{
		MaxVersions:    10,
		ChecksumType:   "md5",
		AutoCleanup:    true,
		StorageType:    "local",
		EnableMetadata: true,
	}
}

// CreateFile creates a new file record with initial version
func CreateFile(db *gorm.DB, fileHeader *multipart.FileHeader, entityType string, entityID uint, userID uint, userName string, dirPath string, config FileVersionConfig) (*entity.File, error) {
	// Generate unique file ID
	fileID := uuid.New().String()

	// Calculate checksum
	checksum, err := calculateFileChecksum(fileHeader, config.ChecksumType)
	if err != nil {
		return nil, err
	}

	// Save file physically
	fullPath := "./public" + dirPath
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to create directory",
			Order:   "H-FileVer-Create-1",
		}
	}

	ext := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s_v1%s", fileID, ext)
	filePath := fullPath + fileName

	// Save file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to open file",
			Order:   "H-FileVer-Create-2",
		}
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to create file",
			Order:   "H-FileVer-Create-3",
		}
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to save file",
			Order:   "H-FileVer-Create-4",
		}
	}

	// Create file version record
	fileVersion := entity.FileVersion{
		FileID:         fileID,
		Version:        1,
		FileName:       fileName,
		OriginalName:   fileHeader.Filename,
		FilePath:       dirPath + fileName,
		FileSize:       fileHeader.Size,
		MimeType:       fileHeader.Header.Get("Content-Type"),
		StorageType:    config.StorageType,
		Checksum:       checksum,
		UploadedBy:     userID,
		UploadedByName: userName,
		IsLatest:       true,
	}

	if config.EnableMetadata {
		metadata := FileMetadata{}
		metadataJSON, _ := json.Marshal(metadata)
		fileVersion.Metadata = string(metadataJSON)
	}

	if err := db.Create(&fileVersion).Error; err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to create file version",
			Order:   "H-FileVer-Create-5",
		}
	}

	// Create file record
	file := &entity.File{
		ID:              fileID,
		Name:            fileHeader.Filename,
		CurrentVersion:  1,
		TotalVersions:   1,
		LatestVersionID: fileVersion.ID,
		EntityType:      entityType,
		EntityID:        entityID,
		CreatedBy:       userID,
		CreatedByName:   userName,
	}

	if err := db.Create(file).Error; err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to create file record",
			Order:   "H-FileVer-Create-6",
		}
	}

	return file, nil
}

// AddFileVersion adds a new version to an existing file
func AddFileVersion(db *gorm.DB, fileID string, fileHeader *multipart.FileHeader, userID uint, userName string, dirPath string, config FileVersionConfig) (*entity.FileVersion, error) {
	// Get file record
	var file entity.File
	if err := db.First(&file, "id = ?", fileID).Error; err != nil {
		return nil, &apierror.NotFoundError{
			Message: "File not found",
			Order:   "H-FileVer-Add-1",
		}
	}

	// Calculate checksum
	checksum, err := calculateFileChecksum(fileHeader, config.ChecksumType)
	if err != nil {
		return nil, err
	}

	// Check if checksum already exists (duplicate version)
	var existingVersion entity.FileVersion
	if err := db.Where("file_id = ? AND checksum = ?", fileID, checksum).First(&existingVersion).Error; err == nil {
		return &existingVersion, &apierror.UnprocessableEntityError{
			Message: "Identical version already exists",
			Order:   "H-FileVer-Add-2",
		}
	}

	// Increment version
	newVersion := file.CurrentVersion + 1

	// Save file physically
	fullPath := "./public" + dirPath
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to create directory",
			Order:   "H-FileVer-Add-3",
		}
	}

	ext := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s_v%d%s", fileID, newVersion, ext)
	filePath := fullPath + fileName

	// Save file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to open file",
			Order:   "H-FileVer-Add-4",
		}
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to create file",
			Order:   "H-FileVer-Add-5",
		}
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to save file",
			Order:   "H-FileVer-Add-6",
		}
	}

	// Mark previous latest version as not latest
	if err := db.Model(&entity.FileVersion{}).Where("file_id = ? AND is_latest = ?", fileID, true).Update("is_latest", false).Error; err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to update previous version",
			Order:   "H-FileVer-Add-7",
		}
	}

	// Create new file version record
	fileVersion := entity.FileVersion{
		FileID:         fileID,
		Version:        newVersion,
		FileName:       fileName,
		OriginalName:   fileHeader.Filename,
		FilePath:       dirPath + fileName,
		FileSize:       fileHeader.Size,
		MimeType:       fileHeader.Header.Get("Content-Type"),
		StorageType:    config.StorageType,
		Checksum:       checksum,
		UploadedBy:     userID,
		UploadedByName: userName,
		IsLatest:       true,
	}

	if config.EnableMetadata {
		metadata := FileMetadata{}
		metadataJSON, _ := json.Marshal(metadata)
		fileVersion.Metadata = string(metadataJSON)
	}

	if err := db.Create(&fileVersion).Error; err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to create file version",
			Order:   "H-FileVer-Add-8",
		}
	}

	// Update file record
	file.CurrentVersion = newVersion
	file.TotalVersions = newVersion
	file.LatestVersionID = fileVersion.ID

	if err := db.Save(&file).Error; err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to update file record",
			Order:   "H-FileVer-Add-9",
		}
	}

	// Auto cleanup old versions if enabled
	if config.AutoCleanup && config.MaxVersions > 0 {
		if err := CleanupOldVersions(db, fileID, config.MaxVersions); err != nil {
			// Log error but don't fail the operation
		}
	}

	return &fileVersion, nil
}

// GetFileVersions gets all versions of a file
func GetFileVersions(db *gorm.DB, fileID string) ([]entity.FileVersion, error) {
	var versions []entity.FileVersion
	if err := db.Where("file_id = ?", fileID).Order("version DESC").Find(&versions).Error; err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to get file versions",
			Order:   "H-FileVer-GetVersions-1",
		}
	}
	return versions, nil
}

// GetFileVersion gets a specific version of a file
func GetFileVersion(db *gorm.DB, fileID string, version int) (*entity.FileVersion, error) {
	var fileVersion entity.FileVersion
	if err := db.Where("file_id = ? AND version = ?", fileID, version).First(&fileVersion).Error; err != nil {
		return nil, &apierror.NotFoundError{
			Message: "File version not found",
			Order:   "H-FileVer-GetVersion-1",
		}
	}
	return &fileVersion, nil
}

// GetLatestVersion gets the latest version of a file
func GetLatestVersion(db *gorm.DB, fileID string) (*entity.FileVersion, error) {
	var fileVersion entity.FileVersion
	if err := db.Where("file_id = ? AND is_latest = ?", fileID, true).First(&fileVersion).Error; err != nil {
		return nil, &apierror.NotFoundError{
			Message: "Latest file version not found",
			Order:   "H-FileVer-GetLatest-1",
		}
	}
	return &fileVersion, nil
}

// RestoreVersion restores a previous version as the latest
func RestoreVersion(db *gorm.DB, fileID string, version int, userID uint, userName string) (*entity.FileVersion, error) {
	// Get the version to restore
	oldVersion, err := GetFileVersion(db, fileID, version)
	if err != nil {
		return nil, err
	}

	// Get file record
	var file entity.File
	if err := db.First(&file, "id = ?", fileID).Error; err != nil {
		return nil, &apierror.NotFoundError{
			Message: "File not found",
			Order:   "H-FileVer-Restore-1",
		}
	}

	// Create a new version from the old one
	newVersion := file.CurrentVersion + 1

	// Copy the old file to new version
	ext := filepath.Ext(oldVersion.FileName)
	newFileName := fmt.Sprintf("%s_v%d%s", fileID, newVersion, ext)
	fullPath := "./public" + filepath.Dir(oldVersion.FilePath)
	newFilePath := fullPath + "/" + newFileName

	// Copy file
	if err := copyFile("./public"+oldVersion.FilePath, newFilePath); err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to copy file",
			Order:   "H-FileVer-Restore-2",
		}
	}

	// Mark previous latest version as not latest
	if err := db.Model(&entity.FileVersion{}).Where("file_id = ? AND is_latest = ?", fileID, true).Update("is_latest", false).Error; err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to update previous version",
			Order:   "H-FileVer-Restore-3",
		}
	}

	// Create new version record
	newFileVersion := entity.FileVersion{
		FileID:         fileID,
		Version:        newVersion,
		FileName:       newFileName,
		OriginalName:   oldVersion.OriginalName,
		FilePath:       filepath.Dir(oldVersion.FilePath) + "/" + newFileName,
		FileSize:       oldVersion.FileSize,
		MimeType:       oldVersion.MimeType,
		StorageType:    oldVersion.StorageType,
		Checksum:       oldVersion.Checksum,
		UploadedBy:     userID,
		UploadedByName: userName,
		IsLatest:       true,
		Metadata:       oldVersion.Metadata,
	}

	if err := db.Create(&newFileVersion).Error; err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to create restored version",
			Order:   "H-FileVer-Restore-4",
		}
	}

	// Update file record
	file.CurrentVersion = newVersion
	file.TotalVersions = newVersion
	file.LatestVersionID = newFileVersion.ID

	if err := db.Save(&file).Error; err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to update file record",
			Order:   "H-FileVer-Restore-5",
		}
	}

	return &newFileVersion, nil
}

// DeleteFileVersion deletes a specific version (soft delete)
func DeleteFileVersion(db *gorm.DB, fileID string, version int) error {
	var fileVersion entity.FileVersion
	if err := db.Where("file_id = ? AND version = ?", fileID, version).First(&fileVersion).Error; err != nil {
		return &apierror.NotFoundError{
			Message: "File version not found",
			Order:   "H-FileVer-Delete-1",
		}
	}

	// Cannot delete the latest version
	if fileVersion.IsLatest {
		return &apierror.UnprocessableEntityError{
			Message: "Cannot delete the latest version",
			Order:   "H-FileVer-Delete-2",
		}
	}

	if err := db.Delete(&fileVersion).Error; err != nil {
		return &apierror.InternalServerError{
			Message: "Failed to delete file version",
			Order:   "H-FileVer-Delete-3",
		}
	}

	return nil
}

// CleanupOldVersions removes old versions beyond the limit
func CleanupOldVersions(db *gorm.DB, fileID string, maxVersions int) error {
	var versions []entity.FileVersion
	if err := db.Where("file_id = ?", fileID).Order("version DESC").Find(&versions).Error; err != nil {
		return err
	}

	if len(versions) <= maxVersions {
		return nil
	}

	// Delete old versions beyond the limit
	versionsToDelete := versions[maxVersions:]
	for _, version := range versionsToDelete {
		// Delete physical file
		filePath := "./public" + version.FilePath
		os.Remove(filePath)

		// Delete record
		db.Unscoped().Delete(&version)
	}

	return nil
}

// CompareVersions compares two versions and returns differences
func CompareVersions(db *gorm.DB, fileID string, version1, version2 int) (map[string]interface{}, error) {
	v1, err := GetFileVersion(db, fileID, version1)
	if err != nil {
		return nil, err
	}

	v2, err := GetFileVersion(db, fileID, version2)
	if err != nil {
		return nil, err
	}

	comparison := map[string]interface{}{
		"version1": map[string]interface{}{
			"version":     v1.Version,
			"file_name":   v1.FileName,
			"file_size":   v1.FileSize,
			"checksum":    v1.Checksum,
			"uploaded_by": v1.UploadedByName,
			"uploaded_at": v1.CreatedAt,
		},
		"version2": map[string]interface{}{
			"version":     v2.Version,
			"file_name":   v2.FileName,
			"file_size":   v2.FileSize,
			"checksum":    v2.Checksum,
			"uploaded_by": v2.UploadedByName,
			"uploaded_at": v2.CreatedAt,
		},
		"differences": map[string]interface{}{
			"size_diff":     v2.FileSize - v1.FileSize,
			"same_checksum": v1.Checksum == v2.Checksum,
			"time_diff":     v2.CreatedAt.Sub(v1.CreatedAt).String(),
		},
	}

	return comparison, nil
}

// GetFileHistory gets the complete history of a file
func GetFileHistory(db *gorm.DB, fileID string) (*entity.File, error) {
	var file entity.File
	if err := db.Preload("Versions", func(db *gorm.DB) *gorm.DB {
		return db.Order("version DESC")
	}).First(&file, "id = ?", fileID).Error; err != nil {
		return nil, &apierror.NotFoundError{
			Message: "File not found",
			Order:   "H-FileVer-History-1",
		}
	}
	return &file, nil
}

// Helper functions

func calculateFileChecksum(fileHeader *multipart.FileHeader, checksumType string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", &apierror.InternalServerError{
			Message: "Failed to open file for checksum",
			Order:   "H-Checksum-1",
		}
	}
	defer src.Close()

	if checksumType == "sha256" {
		hash := sha256.New()
		if _, err := io.Copy(hash, src); err != nil {
			return "", &apierror.InternalServerError{
				Message: "Failed to calculate checksum",
				Order:   "H-Checksum-2",
			}
		}
		return hex.EncodeToString(hash.Sum(nil)), nil
	}

	// Default to MD5
	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", &apierror.InternalServerError{
			Message: "Failed to calculate checksum",
			Order:   "H-Checksum-3",
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
