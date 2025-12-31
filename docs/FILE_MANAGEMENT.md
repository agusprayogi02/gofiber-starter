# File Management Documentation

This document explains the comprehensive file management system including file validation, upload, cloud storage, image processing, and versioning features.

## Table of Contents

- [File Validation](#file-validation)
- [File Upload](#file-upload)
- [Cloud Storage (AWS S3)](#cloud-storage-aws-s3)
- [Image Processing](#image-processing)
- [File Versioning](#file-versioning)

---

## File Validation

The file validation system provides secure file upload with magic number checking, size limits, and MIME type detection.

### Features

- **File size validation**
- **Extension whitelist**
- **MIME type detection from content**
- **Magic number verification** (prevents file spoofing)
- **Preset configurations** for common file types

### Supported File Types

- **Images**: JPG, JPEG, PNG, GIF, WEBP
- **Documents**: PDF, ZIP
- **Videos**: MP4

### Usage

#### Basic Validation

```go
import "starter-gofiber/helper"

// Using preset configuration for images
config := helper.DefaultImageConfig() // Max 5MB
err := helper.ValidateFile(file, config)
if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error": err,
    })
}
```

#### Custom Configuration

```go
config := helper.FileValidationConfig{
    MaxSize:          10 * 1024 * 1024, // 10 MB
    AllowedTypes:     []string{"image/jpeg", "image/png", "application/pdf"},
    AllowedExts:      []string{".jpg", ".jpeg", ".png", ".pdf"},
    CheckMagicNumber: true,
}

err := helper.ValidateFile(file, config)
```

#### Preset Configurations

```go
// For images (Max 5MB)
imageConfig := helper.DefaultImageConfig()

// For documents (Max 10MB)
docConfig := helper.DefaultDocumentConfig()

// For videos (Max 100MB)
videoConfig := helper.DefaultVideoConfig()
```

#### File Type Checkers

```go
// Check if file is an image
isImage := helper.IsImage(ext)

// Check if file is a video
isVideo := helper.IsVideo(ext)

// Check if file is a document
isDoc := helper.IsDocument(ext)
```

#### File Size Formatting

```go
// Convert bytes to human-readable format
sizeStr := helper.FormatFileSize(1024 * 1024) // Returns "1.00 MB"
```

### Magic Number Verification

The validator checks file headers (magic numbers) to prevent file type spoofing:

```go
// Supported magic numbers:
// - JPG/JPEG: FF D8 FF
// - PNG: 89 50 4E 47
// - GIF: 47 49 46 38
// - PDF: 25 50 44 46
// - ZIP: 50 4B 03 04 or 50 4B 05 06
// - MP4: Various (ftyp...)
// - WEBP: 52 49 46 46...57 45 42 50
```

---

## File Upload

### Single File Upload with Validation

```go
func UploadFile(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }

    config := helper.DefaultImageConfig()
    result, err := helper.UploadFileWithValidation(c, file, "/uploads/", config)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "File uploaded successfully",
        "file":    result,
    })
}
```

### Multiple File Upload

```go
func UploadMultipleFiles(c *fiber.Ctx) error {
    config := helper.DefaultImageConfig()
    results, err := helper.UploadMultipleFiles(c, "files", "/uploads/", config)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Files uploaded successfully",
        "files":   results,
    })
}
```

### Get File Information

```go
func GetFileInfo(c *fiber.Ctx) error {
    fileName := c.Params("filename")
    
    fileInfo, err := helper.GetFileInfo(fileName, "/uploads/")
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fileInfo)
}
```

### Delete Multiple Files

```go
func DeleteFiles(c *fiber.Ctx) error {
    var req struct {
        FileNames []string `json:"file_names"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }

    err := helper.DeleteMultipleFiles(req.FileNames, "/uploads/")
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Files deleted successfully",
    })
}
```

---

## Cloud Storage (AWS S3)

The cloud storage implementation supports **AWS S3 and all S3-compatible services** including MinIO, DigitalOcean Spaces, Wasabi, Cloudflare R2, and Backblaze B2.

### Supported Providers

- ✅ **AWS S3** (official)
- ✅ **MinIO** (self-hosted, perfect for local development)
- ✅ **DigitalOcean Spaces**
- ✅ **Wasabi Cloud Storage**
- ✅ **Cloudflare R2**
- ✅ **Backblaze B2**
- ✅ **Any S3-compatible service**

### Configuration

Add to `.env`:

```env
# Cloud Storage Configuration
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_REGION=us-east-1
AWS_S3_BUCKET=your-bucket-name
AWS_S3_ENDPOINT= # Leave empty for AWS S3, set for S3-compatible services
```

### Provider-Specific Examples

```env
# AWS S3 (Official)
AWS_S3_ENDPOINT=

# MinIO (Local Development)
AWS_S3_ENDPOINT=http://localhost:9000

# DigitalOcean Spaces
AWS_S3_ENDPOINT=https://nyc3.digitaloceanspaces.com

# Wasabi
AWS_S3_ENDPOINT=https://s3.wasabisys.com

# Cloudflare R2
AWS_S3_ENDPOINT=https://<account_id>.r2.cloudflarestorage.com

# Backblaze B2
AWS_S3_ENDPOINT=https://s3.us-west-004.backblazeb2.com
```

See [CLOUD_STORAGE_EXAMPLES.md](CLOUD_STORAGE_EXAMPLES.md) for complete configuration examples and local MinIO setup with Docker.

### Initialize S3 Client

```go
import "starter-gofiber/config"

s3Client, err := config.InitS3Client()
if err != nil {
    log.Fatal(err)
}
```

### Upload File to S3

```go
func UploadToS3(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }

    s3Client, _ := config.InitS3Client()
    result, err := s3Client.UploadFile(c.Context(), file, "uploads/")
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "File uploaded to S3 successfully",
        "file":    result,
    })
}
```

### Upload Multiple Files to S3

```go
func UploadMultipleToS3(c *fiber.Ctx) error {
    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to parse form",
        })
    }

    files := form.File["files"]
    s3Client, _ := config.InitS3Client()
    
    results, err := s3Client.UploadMultipleFiles(c.Context(), files, "uploads/")
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Files uploaded to S3 successfully",
        "files":   results,
    })
}
```

### Delete File from S3

```go
func DeleteFromS3(c *fiber.Ctx) error {
    key := c.Query("key")
    
    s3Client, _ := config.InitS3Client()
    err := s3Client.DeleteFile(c.Context(), key)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "File deleted from S3 successfully",
    })
}
```

### Generate Presigned URL

```go
func GetPresignedURL(c *fiber.Ctx) error {
    key := c.Query("key")
    expiration := 15 * time.Minute
    
    s3Client, _ := config.InitS3Client()
    url, err := s3Client.GetPresignedURL(c.Context(), key, expiration)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "url": url,
    })
}
```

### Generate Presigned Upload URL

```go
func GetPresignedUploadURL(c *fiber.Ctx) error {
    key := c.Query("key")
    expiration := 15 * time.Minute
    
    s3Client, _ := config.InitS3Client()
    url, err := s3Client.GetPresignedUploadURL(c.Context(), key, expiration)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "upload_url": url,
    })
}
```

### List Files in S3

```go
func ListS3Files(c *fiber.Ctx) error {
    prefix := c.Query("prefix", "uploads/")
    
    s3Client, _ := config.InitS3Client()
    files, err := s3Client.ListFiles(c.Context(), prefix, 100)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "files": files,
    })
}
```

### Copy File in S3

```go
func CopyS3File(c *fiber.Ctx) error {
    var req struct {
        SourceKey string `json:"source_key"`
        DestKey   string `json:"dest_key"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }

    s3Client, _ := config.InitS3Client()
    err := s3Client.CopyFile(c.Context(), req.SourceKey, req.DestKey)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "File copied successfully",
    })
}
```

### Get File Metadata

```go
func GetS3FileMetadata(c *fiber.Ctx) error {
    key := c.Query("key")
    
    s3Client, _ := config.InitS3Client()
    metadata, err := s3Client.GetFileMetadata(c.Context(), key)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "metadata": metadata,
    })
}
```

### S3-Compatible Services

The S3 client supports S3-compatible services like MinIO, DigitalOcean Spaces, etc:

```env
# For MinIO
AWS_S3_ENDPOINT=http://localhost:9000

# For DigitalOcean Spaces
AWS_S3_ENDPOINT=https://nyc3.digitaloceanspaces.com
```

### Local Development with MinIO

For local development, use MinIO with Docker:

```bash
# Start MinIO with docker-compose
docker-compose -f docker-compose.minio.yml up -d

# Access MinIO Console at: http://localhost:9001
# Username: minioadmin
# Password: minioadmin
```

See [CLOUD_STORAGE_EXAMPLES.md](CLOUD_STORAGE_EXAMPLES.md) for:
- Complete MinIO setup guide
- Configuration examples for all providers
- Performance comparison
- Migration between providers
- Best practices

---

## Image Processing

### Features

- **Resize** images
- **Crop** images
- **Rotate** images
- **Flip** (horizontal/vertical)
- **Blur** and **sharpen**
- **Adjust brightness** and **contrast**
- **Convert to grayscale**
- **Automatic thumbnail** and **medium size** generation
- **JPEG quality** control

### Process Image with Thumbnails

```go
func ProcessUploadedImage(c *fiber.Ctx) error {
    file, err := c.FormFile("image")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "No image uploaded",
        })
    }

    // Configure image processing
    config := helper.DefaultImageProcessConfig()
    // config.Quality = 90
    // config.ThumbWidth = 200
    // config.ThumbHeight = 200

    result, err := helper.ProcessImage(file, "/uploads/images/", config)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Image processed successfully",
        "image":   result,
    })
}
```

### Custom Processing Configuration

```go
config := helper.ImageProcessConfig{
    Quality:      85,     // JPEG quality (1-100)
    CreateThumb:  true,   // Create thumbnail
    ThumbWidth:   150,    // Thumbnail width
    ThumbHeight:  150,    // Thumbnail height
    CreateMedium: true,   // Create medium size
    MediumWidth:  800,    // Medium width
    MediumHeight: 600,    // Medium height
    Format:       "jpeg", // Output format: jpeg, png, webp
}
```

### Resize Image

```go
func ResizeImage(c *fiber.Ctx) error {
    imagePath := "./public/uploads/original.jpg"
    outputPath := "./public/uploads/resized.jpg"
    
    err := helper.ResizeImage(imagePath, 800, 600, outputPath)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Image resized successfully",
    })
}
```

### Crop Image

```go
func CropImage(c *fiber.Ctx) error {
    imagePath := "./public/uploads/original.jpg"
    outputPath := "./public/uploads/cropped.jpg"
    
    // Crop from (x:100, y:100) with width 400 and height 300
    err := helper.CropImage(imagePath, 100, 100, 400, 300, outputPath)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Image cropped successfully",
    })
}
```

### Rotate Image

```go
func RotateImage(c *fiber.Ctx) error {
    imagePath := "./public/uploads/original.jpg"
    outputPath := "./public/uploads/rotated.jpg"
    
    // Rotate 90 degrees
    err := helper.RotateImage(imagePath, 90, outputPath)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Image rotated successfully",
    })
}
```

### Flip Image

```go
func FlipImage(c *fiber.Ctx) error {
    imagePath := "./public/uploads/original.jpg"
    outputPath := "./public/uploads/flipped.jpg"
    
    // Flip horizontally (true) or vertically (false)
    err := helper.FlipImage(imagePath, true, outputPath)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Image flipped successfully",
    })
}
```

### Apply Blur

```go
func ApplyBlur(c *fiber.Ctx) error {
    imagePath := "./public/uploads/original.jpg"
    outputPath := "./public/uploads/blurred.jpg"
    
    // Apply Gaussian blur with sigma 2.0
    err := helper.ApplyBlur(imagePath, 2.0, outputPath)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Blur applied successfully",
    })
}
```

### Apply Sharpen

```go
func ApplySharpen(c *fiber.Ctx) error {
    imagePath := "./public/uploads/original.jpg"
    outputPath := "./public/uploads/sharpened.jpg"
    
    // Apply sharpening with sigma 1.0
    err := helper.ApplySharpen(imagePath, 1.0, outputPath)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Sharpening applied successfully",
    })
}
```

### Adjust Brightness

```go
func AdjustBrightness(c *fiber.Ctx) error {
    imagePath := "./public/uploads/original.jpg"
    outputPath := "./public/uploads/bright.jpg"
    
    // Adjust brightness: -100 to 100 (negative = darker, positive = brighter)
    err := helper.AdjustBrightness(imagePath, 20, outputPath)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Brightness adjusted successfully",
    })
}
```

### Adjust Contrast

```go
func AdjustContrast(c *fiber.Ctx) error {
    imagePath := "./public/uploads/original.jpg"
    outputPath := "./public/uploads/contrast.jpg"
    
    // Adjust contrast: -100 to 100
    err := helper.AdjustContrast(imagePath, 20, outputPath)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Contrast adjusted successfully",
    })
}
```

### Convert to Grayscale

```go
func ConvertToGrayscale(c *fiber.Ctx) error {
    imagePath := "./public/uploads/original.jpg"
    outputPath := "./public/uploads/grayscale.jpg"
    
    err := helper.ConvertToGrayscale(imagePath, outputPath)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Image converted to grayscale successfully",
    })
}
```

### Get Image Dimensions

```go
func GetImageDimensions(c *fiber.Ctx) error {
    imagePath := "./public/uploads/original.jpg"
    
    width, height, err := helper.GetImageDimensions(imagePath)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "width":  width,
        "height": height,
    })
}
```

---

## File Versioning

The file versioning system tracks all changes to files with complete history and metadata.

### Features

- **Automatic version tracking**
- **Checksum verification** (MD5/SHA256)
- **Version comparison**
- **Restore previous versions**
- **Auto-cleanup old versions**
- **Metadata storage**
- **Duplicate detection**

### Database Entities

The system uses two tables:

1. **files** - Main file records
2. **file_versions** - Version history

### Create File with Versioning

```go
func CreateVersionedFile(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }

    // Get user info from context
    userID := c.Locals("user_id").(uint)
    userName := c.Locals("username").(string)
    
    config := helper.DefaultFileVersionConfig()
    db := config.DB // Your database instance
    
    fileRecord, err := helper.CreateFile(
        db,
        file,
        "post",     // entity type
        1,          // entity ID
        userID,
        userName,
        "/uploads/",
        config,
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "File created successfully",
        "file":    fileRecord,
    })
}
```

### Add New Version

```go
func AddFileVersion(c *fiber.Ctx) error {
    fileID := c.Params("id")
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }

    userID := c.Locals("user_id").(uint)
    userName := c.Locals("username").(string)
    
    config := helper.DefaultFileVersionConfig()
    db := config.DB
    
    version, err := helper.AddFileVersion(
        db,
        fileID,
        file,
        userID,
        userName,
        "/uploads/",
        config,
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "New version added successfully",
        "version": version,
    })
}
```

### Get All Versions

```go
func GetFileVersions(c *fiber.Ctx) error {
    fileID := c.Params("id")
    db := config.DB
    
    versions, err := helper.GetFileVersions(db, fileID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "versions": versions,
    })
}
```

### Get Specific Version

```go
func GetFileVersion(c *fiber.Ctx) error {
    fileID := c.Params("id")
    version, _ := strconv.Atoi(c.Params("version"))
    db := config.DB
    
    fileVersion, err := helper.GetFileVersion(db, fileID, version)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fileVersion)
}
```

### Get Latest Version

```go
func GetLatestVersion(c *fiber.Ctx) error {
    fileID := c.Params("id")
    db := config.DB
    
    version, err := helper.GetLatestVersion(db, fileID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(version)
}
```

### Restore Previous Version

```go
func RestoreFileVersion(c *fiber.Ctx) error {
    fileID := c.Params("id")
    version, _ := strconv.Atoi(c.Params("version"))
    
    userID := c.Locals("user_id").(uint)
    userName := c.Locals("username").(string)
    db := config.DB
    
    restoredVersion, err := helper.RestoreVersion(
        db,
        fileID,
        version,
        userID,
        userName,
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Version restored successfully",
        "version": restoredVersion,
    })
}
```

### Delete Version

```go
func DeleteFileVersion(c *fiber.Ctx) error {
    fileID := c.Params("id")
    version, _ := strconv.Atoi(c.Params("version"))
    db := config.DB
    
    err := helper.DeleteFileVersion(db, fileID, version)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Version deleted successfully",
    })
}
```

### Compare Versions

```go
func CompareFileVersions(c *fiber.Ctx) error {
    fileID := c.Params("id")
    version1, _ := strconv.Atoi(c.Query("v1"))
    version2, _ := strconv.Atoi(c.Query("v2"))
    db := config.DB
    
    comparison, err := helper.CompareVersions(db, fileID, version1, version2)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(comparison)
}
```

### Get File History

```go
func GetFileHistory(c *fiber.Ctx) error {
    fileID := c.Params("id")
    db := config.DB
    
    file, err := helper.GetFileHistory(db, fileID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(file)
}
```

### Configuration Options

```go
config := helper.FileVersionConfig{
    MaxVersions:    10,        // Keep last 10 versions
    ChecksumType:   "md5",     // or "sha256"
    AutoCleanup:    true,      // Auto-delete old versions
    StorageType:    "local",   // or "s3"
    EnableMetadata: true,      // Store additional metadata
}
```

### Metadata

Store additional information with each version:

```go
metadata := helper.FileMetadata{
    Width:     1920,
    Height:    1080,
    Duration:  120, // For videos
    PageCount: 10,  // For PDFs
    IPAddress: c.IP(),
    UserAgent: c.Get("User-Agent"),
    Comments:  "Updated logo",
}
```

---

## Complete Example: File Upload Handler

```go
package handler

import (
    "starter-gofiber/config"
    "starter-gofiber/helper"
    
    "github.com/gofiber/fiber/v2"
)

func UploadFile(c *fiber.Ctx) error {
    // Get file from form
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }

    // Validate file
    validationConfig := helper.DefaultImageConfig()
    if err := helper.ValidateFile(file, validationConfig); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err,
        })
    }

    // Process image (resize, create thumbnails)
    processConfig := helper.DefaultImageProcessConfig()
    result, err := helper.ProcessImage(file, "/uploads/images/", processConfig)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    // Create file with versioning
    userID := c.Locals("user_id").(uint)
    userName := c.Locals("username").(string)
    versionConfig := helper.DefaultFileVersionConfig()
    
    fileRecord, err := helper.CreateFile(
        config.DB,
        file,
        "post",
        1,
        userID,
        userName,
        "/uploads/images/",
        versionConfig,
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err,
        })
    }

    return c.JSON(fiber.Map{
        "message": "File uploaded and processed successfully",
        "file":    fileRecord,
        "images":  result,
    })
}
```

---

## Error Handling

All file management functions return custom error types:

- `BadRequestError` - Invalid input
- `NotFoundError` - File not found
- `UnprocessableEntityError` - Validation failed
- `InternalServerError` - Server error

Example error handling:

```go
if err != nil {
    switch err.(type) {
    case *helper.BadRequestError:
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
    case *helper.NotFoundError:
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
    case *helper.UnprocessableEntityError:
        return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err})
    default:
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
    }
}
```

---

## Best Practices

1. **Always validate files** before processing or uploading
2. **Use checksums** to detect duplicate files
3. **Enable auto-cleanup** to prevent storage bloat
4. **Set appropriate size limits** based on your use case
5. **Use S3 presigned URLs** for large file uploads (client-side direct upload)
6. **Create thumbnails** for images to improve performance
7. **Track file versions** for important documents
8. **Use magic number checking** to prevent file type spoofing
9. **Store file metadata** for better organization and search
10. **Implement access control** for sensitive files

---

## Security Considerations

1. **Magic number verification** prevents file extension spoofing
2. **File size limits** prevent DoS attacks
3. **MIME type validation** ensures only allowed file types
4. **Checksum verification** detects file tampering
5. **Access control** via presigned URLs (time-limited)
6. **Soft delete** for file versions (can be recovered)

---

## Performance Tips

1. **Use concurrent upload** for multiple files
2. **Enable S3 transfer acceleration** for faster uploads
3. **Use CloudFront CDN** for file delivery
4. **Compress images** to reduce storage and bandwidth
5. **Lazy load thumbnails** instead of full images
6. **Clean up old versions** automatically
7. **Use read replicas** for high-traffic file metadata queries

---

## Dependencies

- `github.com/aws/aws-sdk-go-v2/*` - AWS S3 SDK
- `github.com/disintegration/imaging` - Image processing
- `github.com/google/uuid` - UUID generation

Install all dependencies:

```bash
go get github.com/aws/aws-sdk-go-v2/aws
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/credentials
go get github.com/aws/aws-sdk-go-v2/service/s3
go get github.com/aws/aws-sdk-go-v2/feature/s3/manager
go get github.com/disintegration/imaging
```
