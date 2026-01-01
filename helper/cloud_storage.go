package helper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

// S3Config represents S3 configuration
type S3Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
	Endpoint  string // For S3-compatible services like MinIO
}

// S3Client wraps the S3 client and uploader
type S3Client struct {
	client   *s3.Client
	uploader *manager.Uploader
	config   S3Config
}

// S3UploadResult represents the result of S3 upload
type S3UploadResult struct {
	Key      string `json:"key"`
	Location string `json:"location"`
	Bucket   string `json:"bucket"`
	Size     int64  `json:"size"`
	ETag     string `json:"etag,omitempty"`
}

// NewS3Client creates a new S3 client
func NewS3Client(cfg S3Config) (*S3Client, error) {
	ctx := context.Background()

	// Create custom credentials provider
	credProvider := credentials.NewStaticCredentialsProvider(
		cfg.AccessKey,
		cfg.SecretKey,
		"",
	)

	// Load AWS config
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credProvider),
	)
	if err != nil {
		return nil, &InternalServerError{
			Message: "Failed to load AWS config",
			Order:   "H-S3-NewClient-1",
		}
	}

	// Create S3 client
	var client *s3.Client
	if cfg.Endpoint != "" {
		// For S3-compatible services (MinIO, DigitalOcean Spaces, etc.)
		client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true
		})
	} else {
		client = s3.NewFromConfig(awsCfg)
	}

	// Create uploader
	uploader := manager.NewUploader(client)

	return &S3Client{
		client:   client,
		uploader: uploader,
		config:   cfg,
	}, nil
}

// UploadFile uploads a file to S3
func (s *S3Client) UploadFile(ctx context.Context, file *multipart.FileHeader, prefix string) (*S3UploadResult, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return nil, &InternalServerError{
			Message: "Failed to open file",
			Order:   "H-S3-Upload-1",
		}
	}
	defer src.Close()

	// Generate unique key
	ext := filepath.Ext(file.Filename)
	key := fmt.Sprintf("%s%s%s", prefix, uuid.New().String(), ext)

	// Read file content
	buffer := bytes.NewBuffer(nil)
	size, err := io.Copy(buffer, src)
	if err != nil {
		return nil, &InternalServerError{
			Message: "Failed to read file content",
			Order:   "H-S3-Upload-2",
		}
	}

	// Upload to S3
	result, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.config.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(file.Header.Get("Content-Type")),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return nil, &InternalServerError{
			Message: "Failed to upload to S3",
			Order:   "H-S3-Upload-3",
		}
	}

	return &S3UploadResult{
		Key:      key,
		Location: result.Location,
		Bucket:   s.config.Bucket,
		Size:     size,
		ETag:     aws.ToString(result.ETag),
	}, nil
}

// UploadMultipleFiles uploads multiple files to S3
func (s *S3Client) UploadMultipleFiles(ctx context.Context, files []*multipart.FileHeader, prefix string) ([]S3UploadResult, error) {
	results := make([]S3UploadResult, 0, len(files))

	for _, file := range files {
		result, err := s.UploadFile(ctx, file, prefix)
		if err != nil {
			return results, err
		}
		results = append(results, *result)
	}

	return results, nil
}

// DeleteFile deletes a file from S3
func (s *S3Client) DeleteFile(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return &InternalServerError{
			Message: "Failed to delete file from S3",
			Order:   "H-S3-Delete-1",
		}
	}

	return nil
}

// DeleteMultipleFiles deletes multiple files from S3
func (s *S3Client) DeleteMultipleFiles(ctx context.Context, keys []string) error {
	objects := make([]types.ObjectIdentifier, len(keys))
	for i, key := range keys {
		objects[i] = types.ObjectIdentifier{
			Key: aws.String(key),
		}
	}

	_, err := s.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(s.config.Bucket),
		Delete: &types.Delete{
			Objects: objects,
		},
	})
	if err != nil {
		return &InternalServerError{
			Message: "Failed to delete files from S3",
			Order:   "H-S3-DeleteMultiple-1",
		}
	}

	return nil
}

// GetPresignedURL generates a presigned URL for temporary access
func (s *S3Client) GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	presignedReq, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", &InternalServerError{
			Message: "Failed to generate presigned URL",
			Order:   "H-S3-Presign-1",
		}
	}

	return presignedReq.URL, nil
}

// GetPresignedUploadURL generates a presigned URL for direct upload
func (s *S3Client) GetPresignedUploadURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	presignedReq, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", &InternalServerError{
			Message: "Failed to generate presigned upload URL",
			Order:   "H-S3-PresignUpload-1",
		}
	}

	return presignedReq.URL, nil
}

// ListFiles lists files in S3 bucket with prefix
func (s *S3Client) ListFiles(ctx context.Context, prefix string, maxKeys int32) ([]string, error) {
	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.config.Bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int32(maxKeys),
	})
	if err != nil {
		return nil, &InternalServerError{
			Message: "Failed to list files from S3",
			Order:   "H-S3-List-1",
		}
	}

	keys := make([]string, 0, len(result.Contents))
	for _, obj := range result.Contents {
		keys = append(keys, aws.ToString(obj.Key))
	}

	return keys, nil
}

// CopyFile copies a file within S3
func (s *S3Client) CopyFile(ctx context.Context, sourceKey, destKey string) error {
	copySource := fmt.Sprintf("%s/%s", s.config.Bucket, sourceKey)

	_, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.config.Bucket),
		CopySource: aws.String(copySource),
		Key:        aws.String(destKey),
	})
	if err != nil {
		return &InternalServerError{
			Message: "Failed to copy file in S3",
			Order:   "H-S3-Copy-1",
		}
	}

	return nil
}

// GetFileMetadata gets metadata of a file in S3
func (s *S3Client) GetFileMetadata(ctx context.Context, key string) (*s3.HeadObjectOutput, error) {
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, &NotFoundError{
			Message: "File not found in S3",
			Order:   "H-S3-Metadata-1",
		}
	}

	return result, nil
}
