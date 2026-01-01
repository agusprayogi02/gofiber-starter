package storage

import (
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"starter-gofiber/pkg/apierror"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UploadResult represents the result of a file upload
type UploadResult struct {
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
	URL          string `json:"url"`
	Error        error  `json:"error,omitempty"`
}

// UploadMultipleFiles uploads multiple files and returns results
func UploadMultipleFiles(c *fiber.Ctx, formName string, dirPath string, config FileValidationConfig) ([]UploadResult, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, &apierror.BadRequestError{
			Message: "Failed to parse multipart form",
			Order:   "H-UploadMultiple-1",
		}
	}

	files := form.File[formName]
	if len(files) == 0 {
		return nil, &apierror.BadRequestError{
			Message: "No files uploaded",
			Order:   "H-UploadMultiple-2",
		}
	}

	// Create directory if not exists
	fullPath := "./public" + dirPath
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to create upload directory",
			Order:   "H-UploadMultiple-3",
		}
	}

	results := make([]UploadResult, len(files))
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Upload files concurrently
	for i, file := range files {
		wg.Add(1)
		go func(index int, fileHeader *multipart.FileHeader) {
			defer wg.Done()

			result := UploadResult{
				OriginalName: fileHeader.Filename,
				Size:         fileHeader.Size,
			}

			// Validate file
			if err := ValidateFile(fileHeader, config); err != nil {
				result.Error = err
				mu.Lock()
				results[index] = result
				mu.Unlock()
				return
			}

			// Generate unique filename
			ext := filepath.Ext(fileHeader.Filename)
			fileName := uuid.New().String() + ext
			filePath := fullPath + fileName

			// Save file
			if err := c.SaveFile(fileHeader, filePath); err != nil {
				result.Error = &apierror.InternalServerError{
					Message: "Failed to save file",
					Order:   "H-UploadMultiple-4",
				}
				mu.Lock()
				results[index] = result
				mu.Unlock()
				return
			}

			// Detect MIME type
			src, _ := fileHeader.Open()
			buffer := make([]byte, 512)
			src.Read(buffer)
			src.Close()
			mimeType := http.DetectContentType(buffer)

			// Fallback to extension-based MIME type
			if mimeType == "application/octet-stream" {
				ext := filepath.Ext(fileHeader.Filename)
				mimeType = mime.TypeByExtension(ext)
			}

			result.FileName = fileName
			result.MimeType = mimeType
			result.URL = dirPath + fileName

			mu.Lock()
			results[index] = result
			mu.Unlock()
		}(i, file)
	}

	wg.Wait()

	// Check if any uploads failed
	hasError := false
	for _, result := range results {
		if result.Error != nil {
			hasError = true
			break
		}
	}

	if hasError {
		// Cleanup successfully uploaded files if any failed
		for _, result := range results {
			if result.Error == nil && result.FileName != "" {
				DeleteFile(&result.FileName, dirPath)
			}
		}
		return results, &apierror.UnprocessableEntityError{
			Message: "Some files failed to upload",
			Order:   "H-UploadMultiple-5",
		}
	}

	return results, nil
}

// UploadFileWithValidation uploads a single file with validation
func UploadFileWithValidation(c *fiber.Ctx, file *multipart.FileHeader, dirPath string, config FileValidationConfig) (UploadResult, error) {
	result := UploadResult{
		OriginalName: file.Filename,
		Size:         file.Size,
	}

	// Validate file
	if err := ValidateFile(file, config); err != nil {
		return result, err
	}

	// Create directory if not exists
	fullPath := "./public" + dirPath
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return result, &apierror.InternalServerError{
			Message: "Failed to create upload directory",
			Order:   "H-UploadSingle-1",
		}
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	fileName := uuid.New().String() + ext
	filePath := fullPath + fileName

	// Save file
	if err := c.SaveFile(file, filePath); err != nil {
		return result, &apierror.InternalServerError{
			Message: "Failed to save file",
			Order:   "H-UploadSingle-2",
		}
	}

	// Detect MIME type
	src, _ := file.Open()
	buffer := make([]byte, 512)
	src.Read(buffer)
	src.Close()
	mimeType := http.DetectContentType(buffer)

	// Fallback to extension-based MIME type
	if mimeType == "application/octet-stream" {
		ext := filepath.Ext(file.Filename)
		mimeType = mime.TypeByExtension(ext)
	}

	result.FileName = fileName
	result.MimeType = mimeType
	result.URL = dirPath + fileName

	return result, nil
}

// DeleteMultipleFiles deletes multiple files
func DeleteMultipleFiles(fileNames []string, dirPath string) error {
	var errors []error

	for _, fileName := range fileNames {
		if err := DeleteFile(&fileName, dirPath); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return &apierror.InternalServerError{
			Message: fmt.Sprintf("Failed to delete %d files", len(errors)),
			Order:   "H-DeleteMultiple-1",
		}
	}

	return nil
}

// GetFileInfo returns information about an uploaded file
func GetFileInfo(fileName string, dirPath string) (*UploadResult, error) {
	fullPath := "./public" + dirPath + fileName
	fileInfo, err := os.Stat(fullPath)

	if os.IsNotExist(err) {
		return nil, &apierror.NotFoundError{
			Message: "File not found",
			Order:   "H-GetFileInfo-1",
		}
	}

	if err != nil {
		return nil, &apierror.InternalServerError{
			Message: "Failed to get file info",
			Order:   "H-GetFileInfo-2",
		}
	}

	ext := filepath.Ext(fileName)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return &UploadResult{
		FileName: fileName,
		Size:     fileInfo.Size(),
		MimeType: mimeType,
		URL:      dirPath + fileName,
	}, nil
}
