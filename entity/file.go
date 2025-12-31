package entity

import (
	"time"

	"gorm.io/gorm"
)

// FileVersion represents a version of a file
type FileVersion struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	FileID         string         `gorm:"type:varchar(100);not null;index" json:"file_id"`
	Version        int            `gorm:"not null;default:1" json:"version"`
	FileName       string         `gorm:"type:varchar(255);not null" json:"file_name"`
	OriginalName   string         `gorm:"type:varchar(255);not null" json:"original_name"`
	FilePath       string         `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize       int64          `gorm:"not null" json:"file_size"`
	MimeType       string         `gorm:"type:varchar(100)" json:"mime_type"`
	StorageType    string         `gorm:"type:varchar(50);default:'local'" json:"storage_type"` // local, s3, etc.
	S3Key          string         `gorm:"type:varchar(500)" json:"s3_key,omitempty"`
	Checksum       string         `gorm:"type:varchar(64)" json:"checksum"` // MD5 or SHA256
	UploadedBy     uint           `gorm:"index" json:"uploaded_by"`
	UploadedByName string         `gorm:"type:varchar(100)" json:"uploaded_by_name"`
	IsLatest       bool           `gorm:"default:true;index" json:"is_latest"`
	Metadata       string         `gorm:"type:text" json:"metadata,omitempty"` // JSON metadata
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// File represents a file with versioning
type File struct {
	ID              string         `gorm:"primarykey;type:varchar(100)" json:"id"`
	Name            string         `gorm:"type:varchar(255);not null" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	CurrentVersion  int            `gorm:"not null;default:1" json:"current_version"`
	TotalVersions   int            `gorm:"not null;default:1" json:"total_versions"`
	LatestVersionID uint           `gorm:"index" json:"latest_version_id"`
	EntityType      string         `gorm:"type:varchar(50)" json:"entity_type"` // post, user, product, etc.
	EntityID        uint           `gorm:"index" json:"entity_id"`
	CreatedBy       uint           `gorm:"index" json:"created_by"`
	CreatedByName   string         `gorm:"type:varchar(100)" json:"created_by_name"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Versions []FileVersion `gorm:"foreignKey:FileID;references:ID" json:"versions,omitempty"`
}

// BeforeCreate hook
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		// ID will be set by the helper function
	}
	return nil
}
