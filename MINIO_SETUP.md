# MinIO Local Setup - Quick Start

MinIO is a high-performance, S3-compatible object storage server that's perfect for local development and testing.

## Quick Start (5 minutes)

### 1. Start MinIO

```bash
docker-compose -f docker-compose.minio.yml up -d
```

This will start:
- MinIO server on port **9000** (API)
- MinIO Console on port **9001** (Web UI)
- PostgreSQL on port **5432**
- Redis on port **6379**

### 2. Access MinIO Console

Open your browser: **http://localhost:9001**

- **Username**: `minioadmin`
- **Password**: `minioadmin`

### 3. Configure Your App

Update your `.env` file:

```env
AWS_ACCESS_KEY_ID=minioadmin
AWS_SECRET_ACCESS_KEY=minioadmin
AWS_REGION=us-east-1
AWS_S3_BUCKET=uploads
AWS_S3_ENDPOINT=http://localhost:9000
```

### 4. Test File Upload

The buckets are auto-created on startup:
- `uploads` - General file uploads
- `images` - Image files
- `documents` - Document files
- `videos` - Video files

### 5. Access Uploaded Files

Files are publicly accessible at:
```
http://localhost:9000/{bucket-name}/{file-path}
```

Example:
```
http://localhost:9000/uploads/my-file.jpg
```

## Verify Installation

Check if MinIO is running:

```bash
# Check running containers
docker ps | grep minio

# Check MinIO logs
docker logs minio

# Check bucket initialization
docker logs minio-init
```

## Using MinIO Client (mc)

Install MinIO Client for advanced operations:

```bash
# macOS
brew install minio/stable/mc

# Linux
wget https://dl.min.io/client/mc/release/linux-amd64/mc
chmod +x mc
sudo mv mc /usr/local/bin/

# Windows
# Download from https://dl.min.io/client/mc/release/windows-amd64/mc.exe
```

### Common mc Commands

```bash
# Add MinIO alias
mc alias set local http://localhost:9000 minioadmin minioadmin

# List buckets
mc ls local

# List files in bucket
mc ls local/uploads

# Copy file to bucket
mc cp myfile.jpg local/uploads/

# Remove file
mc rm local/uploads/myfile.jpg

# Make bucket public
mc anonymous set download local/uploads

# Make bucket private
mc anonymous set none local/uploads
```

## API Testing with curl

### Upload File

```bash
# Install AWS CLI
pip install awscli

# Configure credentials
aws configure --profile minio
# AWS Access Key ID: minioadmin
# AWS Secret Access Key: minioadmin
# Default region name: us-east-1
# Default output format: json

# Upload file
aws --profile minio --endpoint-url http://localhost:9000 \
    s3 cp myfile.jpg s3://uploads/

# List files
aws --profile minio --endpoint-url http://localhost:9000 \
    s3 ls s3://uploads/

# Download file
aws --profile minio --endpoint-url http://localhost:9000 \
    s3 cp s3://uploads/myfile.jpg ./downloaded.jpg

# Delete file
aws --profile minio --endpoint-url http://localhost:9000 \
    s3 rm s3://uploads/myfile.jpg
```

## Using in Your Go Code

```go
package main

import (
    "context"
    "log"
    "starter-gofiber/config"
)

func main() {
    // Initialize S3 client (automatically uses .env config)
    s3Client, err := config.InitS3Client()
    if err != nil {
        log.Fatal(err)
    }

    // Upload file
    file, _ := c.FormFile("file")
    result, err := s3Client.UploadFile(context.Background(), file, "uploads/")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("File uploaded: %s", result.Location)
}
```

## Switching Between Providers

MinIO for development, AWS S3 for production:

### Development (.env.development)
```env
AWS_ACCESS_KEY_ID=minioadmin
AWS_SECRET_ACCESS_KEY=minioadmin
AWS_REGION=us-east-1
AWS_S3_BUCKET=uploads
AWS_S3_ENDPOINT=http://localhost:9000
```

### Production (.env.production)
```env
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_REGION=us-east-1
AWS_S3_BUCKET=my-production-bucket
AWS_S3_ENDPOINT=
```

**No code changes needed!** Just switch the environment variables.

## Stop MinIO

```bash
# Stop containers
docker-compose -f docker-compose.minio.yml down

# Stop and remove volumes (delete all data)
docker-compose -f docker-compose.minio.yml down -v
```

## Troubleshooting

### Port Already in Use

If port 9000 or 9001 is already in use, edit `docker-compose.minio.yml`:

```yaml
ports:
  - "9002:9000"  # Change 9000 to 9002
  - "9003:9001"  # Change 9001 to 9003
```

Then update your `.env`:
```env
AWS_S3_ENDPOINT=http://localhost:9002
```

### Connection Refused

Make sure MinIO is running:
```bash
docker ps | grep minio
```

If not running:
```bash
docker-compose -f docker-compose.minio.yml up -d
```

### Bucket Not Found

Check if buckets were created:
```bash
docker logs minio-init
```

Should see: `MinIO setup completed successfully!`

If not, recreate buckets:
```bash
docker exec -it minio sh
mc alias set local http://localhost:9000 minioadmin minioadmin
mc mb local/uploads
mc anonymous set download local/uploads
```

## Advanced Configuration

### Custom Bucket Policy

```bash
# Read-only access
mc anonymous set download local/uploads

# Read-write access
mc anonymous set upload local/uploads

# Private (no public access)
mc anonymous set none local/uploads
```

### Enable Versioning

```bash
mc version enable local/uploads
```

### Set Lifecycle Rules

```bash
# Delete files older than 30 days
mc ilm add local/uploads --expiry-days 30
```

## Production Deployment

For production, use MinIO in cluster mode:

```yaml
# docker-compose.prod.yml
services:
  minio1:
    image: minio/minio
    command: server /data{1...4}
    volumes:
      - data1:/data1
      - data2:/data2
      - data3:/data3
      - data4:/data4
```

Or use cloud providers:
- AWS S3
- DigitalOcean Spaces
- Wasabi
- Cloudflare R2

See [docs/CLOUD_STORAGE_EXAMPLES.md](docs/CLOUD_STORAGE_EXAMPLES.md) for all options.

## Resources

- **MinIO Documentation**: https://min.io/docs
- **MinIO GitHub**: https://github.com/minio/minio
- **AWS S3 API Reference**: https://docs.aws.amazon.com/s3/
- **Project Documentation**: [docs/FILE_MANAGEMENT.md](docs/FILE_MANAGEMENT.md)

## Next Steps

1. ✅ MinIO is running
2. ✅ Buckets are created
3. ✅ App is configured
4. 🚀 Start uploading files!

See [docs/FILE_MANAGEMENT.md](docs/FILE_MANAGEMENT.md) for complete file management features including:
- File validation
- Image processing
- File versioning
- Multiple providers
