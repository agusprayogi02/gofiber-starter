package helper

import (
	"fmt"
	"image"
	"image/color"
	"mime/multipart"
	"os"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

// ImageProcessConfig represents image processing configuration
type ImageProcessConfig struct {
	Quality      int  // JPEG quality (1-100)
	CreateThumb  bool // Create thumbnail
	ThumbWidth   int  // Thumbnail width
	ThumbHeight  int  // Thumbnail height
	CreateMedium bool // Create medium size
	MediumWidth  int  // Medium width
	MediumHeight int  // Medium height
	Format       string // Output format: jpeg, png, webp
}

// ImageProcessResult represents the result of image processing
type ImageProcessResult struct {
	OriginalPath  string `json:"original_path"`
	OriginalURL   string `json:"original_url"`
	ThumbnailPath string `json:"thumbnail_path,omitempty"`
	ThumbnailURL  string `json:"thumbnail_url,omitempty"`
	MediumPath    string `json:"medium_path,omitempty"`
	MediumURL     string `json:"medium_url,omitempty"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
}

// DefaultImageProcessConfig returns default image processing configuration
func DefaultImageProcessConfig() ImageProcessConfig {
	return ImageProcessConfig{
		Quality:      85,
		CreateThumb:  true,
		ThumbWidth:   150,
		ThumbHeight:  150,
		CreateMedium: true,
		MediumWidth:  800,
		MediumHeight: 600,
		Format:       "jpeg",
	}
}

// ProcessImage processes an uploaded image (resize, compress, create thumbnails)
func ProcessImage(file *multipart.FileHeader, dirPath string, config ImageProcessConfig) (*ImageProcessResult, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, &InternalServerError{
			Message: "Failed to open uploaded file",
			Order:   "H-ImgProcess-1",
		}
	}
	defer src.Close()

	// Decode image
	img, err := imaging.Decode(src)
	if err != nil {
		return nil, &UnprocessableEntityError{
			Message: "Failed to decode image",
			Order:   "H-ImgProcess-2",
		}
	}

	// Create directory if not exists
	fullPath := "./public" + dirPath
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return nil, &InternalServerError{
			Message: "Failed to create directory",
			Order:   "H-ImgProcess-3",
		}
	}

	result := &ImageProcessResult{
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
	}

	// Generate base filename
	baseFileName := uuid.New().String()
	ext := getImageExt(config.Format)

	// Save original (or compressed original)
	originalFileName := baseFileName + ext
	originalPath := fullPath + originalFileName
	if err := saveImage(img, originalPath, config); err != nil {
		return nil, err
	}
	result.OriginalPath = originalPath
	result.OriginalURL = dirPath + originalFileName

	// Create thumbnail
	if config.CreateThumb {
		thumb := imaging.Thumbnail(img, config.ThumbWidth, config.ThumbHeight, imaging.Lanczos)
		thumbFileName := baseFileName + "_thumb" + ext
		thumbPath := fullPath + thumbFileName
		if err := saveImage(thumb, thumbPath, config); err != nil {
			return nil, err
		}
		result.ThumbnailPath = thumbPath
		result.ThumbnailURL = dirPath + thumbFileName
	}

	// Create medium size
	if config.CreateMedium {
		medium := imaging.Fit(img, config.MediumWidth, config.MediumHeight, imaging.Lanczos)
		mediumFileName := baseFileName + "_medium" + ext
		mediumPath := fullPath + mediumFileName
		if err := saveImage(medium, mediumPath, config); err != nil {
			return nil, err
		}
		result.MediumPath = mediumPath
		result.MediumURL = dirPath + mediumFileName
	}

	return result, nil
}

// ResizeImage resizes an image to specific dimensions
func ResizeImage(imagePath string, width, height int, outputPath string) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return &NotFoundError{
			Message: "Image not found",
			Order:   "H-Resize-1",
		}
	}

	resized := imaging.Resize(img, width, height, imaging.Lanczos)

	if err := imaging.Save(resized, outputPath); err != nil {
		return &InternalServerError{
			Message: "Failed to save resized image",
			Order:   "H-Resize-2",
		}
	}

	return nil
}

// CropImage crops an image to specific dimensions
func CropImage(imagePath string, x, y, width, height int, outputPath string) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return &NotFoundError{
			Message: "Image not found",
			Order:   "H-Crop-1",
		}
	}

	rect := image.Rect(x, y, x+width, y+height)
	cropped := imaging.Crop(img, rect)

	if err := imaging.Save(cropped, outputPath); err != nil {
		return &InternalServerError{
			Message: "Failed to save cropped image",
			Order:   "H-Crop-2",
		}
	}

	return nil
}

// RotateImage rotates an image by specified degrees
func RotateImage(imagePath string, degrees float64, outputPath string) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return &NotFoundError{
			Message: "Image not found",
			Order:   "H-Rotate-1",
		}
	}

	rotated := imaging.Rotate(img, degrees, color.Black)

	if err := imaging.Save(rotated, outputPath); err != nil {
		return &InternalServerError{
			Message: "Failed to save rotated image",
			Order:   "H-Rotate-2",
		}
	}

	return nil
}

// FlipImage flips an image horizontally or vertically
func FlipImage(imagePath string, horizontal bool, outputPath string) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return &NotFoundError{
			Message: "Image not found",
			Order:   "H-Flip-1",
		}
	}

	var flipped *image.NRGBA
	if horizontal {
		flipped = imaging.FlipH(img)
	} else {
		flipped = imaging.FlipV(img)
	}

	if err := imaging.Save(flipped, outputPath); err != nil {
		return &InternalServerError{
			Message: "Failed to save flipped image",
			Order:   "H-Flip-2",
		}
	}

	return nil
}

// ApplyBlur applies Gaussian blur to an image
func ApplyBlur(imagePath string, sigma float64, outputPath string) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return &NotFoundError{
			Message: "Image not found",
			Order:   "H-Blur-1",
		}
	}

	blurred := imaging.Blur(img, sigma)

	if err := imaging.Save(blurred, outputPath); err != nil {
		return &InternalServerError{
			Message: "Failed to save blurred image",
			Order:   "H-Blur-2",
		}
	}

	return nil
}

// ApplySharpen applies sharpening to an image
func ApplySharpen(imagePath string, sigma float64, outputPath string) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return &NotFoundError{
			Message: "Image not found",
			Order:   "H-Sharpen-1",
		}
	}

	sharpened := imaging.Sharpen(img, sigma)

	if err := imaging.Save(sharpened, outputPath); err != nil {
		return &InternalServerError{
			Message: "Failed to save sharpened image",
			Order:   "H-Sharpen-2",
		}
	}

	return nil
}

// AdjustBrightness adjusts the brightness of an image
func AdjustBrightness(imagePath string, percentage float64, outputPath string) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return &NotFoundError{
			Message: "Image not found",
			Order:   "H-Brightness-1",
		}
	}

	adjusted := imaging.AdjustBrightness(img, percentage)

	if err := imaging.Save(adjusted, outputPath); err != nil {
		return &InternalServerError{
			Message: "Failed to save adjusted image",
			Order:   "H-Brightness-2",
		}
	}

	return nil
}

// AdjustContrast adjusts the contrast of an image
func AdjustContrast(imagePath string, percentage float64, outputPath string) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return &NotFoundError{
			Message: "Image not found",
			Order:   "H-Contrast-1",
		}
	}

	adjusted := imaging.AdjustContrast(img, percentage)

	if err := imaging.Save(adjusted, outputPath); err != nil {
		return &InternalServerError{
			Message: "Failed to save adjusted image",
			Order:   "H-Contrast-2",
		}
	}

	return nil
}

// ConvertToGrayscale converts an image to grayscale
func ConvertToGrayscale(imagePath string, outputPath string) error {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return &NotFoundError{
			Message: "Image not found",
			Order:   "H-Grayscale-1",
		}
	}

	grayscale := imaging.Grayscale(img)

	if err := imaging.Save(grayscale, outputPath); err != nil {
		return &InternalServerError{
			Message: "Failed to save grayscale image",
			Order:   "H-Grayscale-2",
		}
	}

	return nil
}

// GetImageDimensions returns the dimensions of an image
func GetImageDimensions(imagePath string) (width int, height int, err error) {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return 0, 0, &NotFoundError{
			Message: "Image not found",
			Order:   "H-Dimensions-1",
		}
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}

// Helper functions

func saveImage(img image.Image, path string, config ImageProcessConfig) error {
	var err error

	switch config.Format {
	case "png":
		err = imaging.Save(img, path, imaging.PNGCompressionLevel(9))
	case "webp":
		// Note: imaging library doesn't support WebP directly
		// You might need to use another library for WebP
		err = imaging.Save(img, path)
	default: // jpeg
		err = imaging.Save(img, path, imaging.JPEGQuality(config.Quality))
	}

	if err != nil {
		return &InternalServerError{
			Message: fmt.Sprintf("Failed to save image: %v", err),
			Order:   "H-SaveImage-1",
		}
	}

	return nil
}

func getImageExt(format string) string {
	switch format {
	case "png":
		return ".png"
	case "webp":
		return ".webp"
	default:
		return ".jpg"
	}
}
