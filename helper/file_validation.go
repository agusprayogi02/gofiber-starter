package helper

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// FileValidationConfig holds validation configuration
type FileValidationConfig struct {
	MaxSize          int64    // Maximum file size in bytes
	AllowedTypes     []string // Allowed MIME types (e.g., "image/jpeg", "application/pdf")
	AllowedExts      []string // Allowed file extensions (e.g., ".jpg", ".pdf")
	CheckMagicNumber bool     // Validate file content by magic number
}

// DefaultImageConfig returns default config for images
func DefaultImageConfig() FileValidationConfig {
	return FileValidationConfig{
		MaxSize: 5 * 1024 * 1024, // 5MB
		AllowedTypes: []string{
			"image/jpeg",
			"image/jpg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedExts:      []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		CheckMagicNumber: true,
	}
}

// DefaultDocumentConfig returns default config for documents
func DefaultDocumentConfig() FileValidationConfig {
	return FileValidationConfig{
		MaxSize: 10 * 1024 * 1024, // 10MB
		AllowedTypes: []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		},
		AllowedExts:      []string{".pdf", ".doc", ".docx", ".xls", ".xlsx"},
		CheckMagicNumber: true,
	}
}

// DefaultVideoConfig returns default config for videos
func DefaultVideoConfig() FileValidationConfig {
	return FileValidationConfig{
		MaxSize: 100 * 1024 * 1024, // 100MB
		AllowedTypes: []string{
			"video/mp4",
			"video/mpeg",
			"video/quicktime",
			"video/x-msvideo",
			"video/webm",
		},
		AllowedExts:      []string{".mp4", ".mpeg", ".mov", ".avi", ".webm"},
		CheckMagicNumber: true,
	}
}

// ValidateFile validates uploaded file against config
func ValidateFile(file *multipart.FileHeader, config FileValidationConfig) error {
	// Check file size
	if file.Size > config.MaxSize {
		return &BadRequestError{
			Message: fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", config.MaxSize),
			Order:   "H-ValidateFile-1",
		}
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Check file extension
	if len(config.AllowedExts) > 0 {
		allowed := false
		for _, allowedExt := range config.AllowedExts {
			if ext == allowedExt {
				allowed = true
				break
			}
		}
		if !allowed {
			return &BadRequestError{
				Message: fmt.Sprintf("File extension %s is not allowed. Allowed: %v", ext, config.AllowedExts),
				Order:   "H-ValidateFile-2",
			}
		}
	}

	// Open file to check content
	src, err := file.Open()
	if err != nil {
		return &InternalServerError{
			Message: "Failed to open file",
			Order:   "H-ValidateFile-3",
		}
	}
	defer src.Close()

	// Read first 512 bytes for MIME type detection
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return &InternalServerError{
			Message: "Failed to read file content",
			Order:   "H-ValidateFile-4",
		}
	}

	// Detect MIME type from content
	contentType := http.DetectContentType(buffer)

	// Check MIME type
	if len(config.AllowedTypes) > 0 {
		allowed := false
		for _, allowedType := range config.AllowedTypes {
			if strings.HasPrefix(contentType, allowedType) {
				allowed = true
				break
			}
		}
		if !allowed {
			return &BadRequestError{
				Message: fmt.Sprintf("File type %s is not allowed. Allowed: %v", contentType, config.AllowedTypes),
				Order:   "H-ValidateFile-5",
			}
		}
	}

	// Check magic number if enabled
	if config.CheckMagicNumber {
		if err := validateMagicNumber(buffer, ext); err != nil {
			return err
		}
	}

	return nil
}

// validateMagicNumber checks file signature (magic number)
func validateMagicNumber(buffer []byte, ext string) error {
	magicNumbers := map[string][]byte{
		".jpg":  {0xFF, 0xD8, 0xFF},
		".jpeg": {0xFF, 0xD8, 0xFF},
		".png":  {0x89, 0x50, 0x4E, 0x47},
		".gif":  {0x47, 0x49, 0x46, 0x38},
		".pdf":  {0x25, 0x50, 0x44, 0x46},
		".zip":  {0x50, 0x4B, 0x03, 0x04},
		".mp4":  {0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70},
		".webp": {0x52, 0x49, 0x46, 0x46},
	}

	expected, exists := magicNumbers[ext]
	if !exists {
		// No magic number to check for this extension
		return nil
	}

	if len(buffer) < len(expected) {
		return &BadRequestError{
			Message: "File content is too short to validate",
			Order:   "H-ValidateMagicNumber-1",
		}
	}

	// Special case for WEBP (need to check further)
	if ext == ".webp" {
		if !bytes.HasPrefix(buffer, expected) {
			return &BadRequestError{
				Message: "File content does not match expected format (invalid WEBP)",
				Order:   "H-ValidateMagicNumber-2",
			}
		}
		// Check for WEBP marker at offset 8
		if len(buffer) >= 12 && string(buffer[8:12]) != "WEBP" {
			return &BadRequestError{
				Message: "File content does not match expected format (invalid WEBP marker)",
				Order:   "H-ValidateMagicNumber-3",
			}
		}
		return nil
	}

	// Check magic number
	if !bytes.HasPrefix(buffer, expected) {
		return &BadRequestError{
			Message: fmt.Sprintf("File content does not match extension %s (possible file spoofing)", ext),
			Order:   "H-ValidateMagicNumber-4",
		}
	}

	return nil
}

// IsImage checks if file is an image
func IsImage(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

// IsVideo checks if file is a video
func IsVideo(contentType string) bool {
	return strings.HasPrefix(contentType, "video/")
}

// IsDocument checks if file is a document
func IsDocument(contentType string) bool {
	documentTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument",
		"text/plain",
	}

	for _, docType := range documentTypes {
		if strings.HasPrefix(contentType, docType) {
			return true
		}
	}
	return false
}

// FormatFileSize formats bytes to human-readable string
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
