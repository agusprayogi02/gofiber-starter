# Cloud Storage Configuration Examples

This document provides configuration examples for various S3-compatible storage services.

## Table of Contents

- [AWS S3 (Official)](#aws-s3-official)
- [MinIO (Self-hosted)](#minio-self-hosted)
- [DigitalOcean Spaces](#digitalocean-spaces)
- [Wasabi Cloud Storage](#wasabi-cloud-storage)
- [Cloudflare R2](#cloudflare-r2)
- [Backblaze B2](#backblaze-b2)
- [Local Development with Docker](#local-development-with-docker)

---

## AWS S3 (Official)

### Configuration

```env
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_REGION=us-east-1
AWS_S3_BUCKET=my-bucket
AWS_S3_ENDPOINT=
```

Leave `AWS_S3_ENDPOINT` empty for official AWS S3.

### Usage

```go
s3Client, err := config.InitS3Client()
if err != nil {
    log.Fatal(err)
}

result, err := s3Client.UploadFile(ctx, file, "uploads/")
```

---

## MinIO (Self-hosted)

MinIO is a high-performance, S3-compatible object storage server that can run locally or in production.

### Configuration

```env
AWS_ACCESS_KEY_ID=minioadmin
AWS_SECRET_ACCESS_KEY=minioadmin
AWS_REGION=us-east-1
AWS_S3_BUCKET=my-bucket
AWS_S3_ENDPOINT=http://localhost:9000
```

### Docker Setup

See [docker-compose.minio.yml](#docker-composeminioyml) below for complete setup.

### Create Bucket (First Time)

```bash
# Using mc (MinIO Client)
mc alias set local http://localhost:9000 minioadmin minioadmin
mc mb local/my-bucket
mc policy set public local/my-bucket

# Or using web UI at http://localhost:9001
```

### Features

- ✅ 100% S3 API compatible
- ✅ High performance (faster than S3)
- ✅ Self-hosted (full control)
- ✅ Free and open source
- ✅ Multi-cloud gateway support
- ✅ Perfect for local development

---

## DigitalOcean Spaces

DigitalOcean Spaces is an S3-compatible object storage service.

### Configuration

```env
AWS_ACCESS_KEY_ID=DO00EXAMPLE9KBQJHTAEQ
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_REGION=nyc3
AWS_S3_BUCKET=my-space
AWS_S3_ENDPOINT=https://nyc3.digitaloceanspaces.com
```

### Available Regions

- `nyc3` - New York 3
- `sfo3` - San Francisco 3
- `sgp1` - Singapore 1
- `ams3` - Amsterdam 3
- `fra1` - Frankfurt 1
- `syd1` - Sydney 1
- `blr1` - Bangalore 1

### CDN Support

DigitalOcean Spaces includes free CDN. Your files will be available at:
- Direct: `https://my-space.nyc3.digitaloceanspaces.com/uploads/file.jpg`
- CDN: `https://my-space.nyc3.cdn.digitaloceanspaces.com/uploads/file.jpg`

---

## Wasabi Cloud Storage

Wasabi is a low-cost S3-compatible cloud storage service.

### Configuration

```env
AWS_ACCESS_KEY_ID=WASABI_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY=WASABI_SECRET_ACCESS_KEY
AWS_REGION=us-east-1
AWS_S3_BUCKET=my-bucket
AWS_S3_ENDPOINT=https://s3.wasabisys.com
```

### Available Regions

- `https://s3.wasabisys.com` - US East 1 (Ashburn)
- `https://s3.us-east-2.wasabisys.com` - US East 2 (Ashburn)
- `https://s3.us-central-1.wasabisys.com` - US Central 1 (Texas)
- `https://s3.us-west-1.wasabisys.com` - US West 1 (Oregon)
- `https://s3.eu-central-1.wasabisys.com` - EU Central 1 (Amsterdam)
- `https://s3.eu-west-1.wasabisys.com` - EU West 1 (London)
- `https://s3.eu-west-2.wasabisys.com` - EU West 2 (Paris)
- `https://s3.ap-northeast-1.wasabisys.com` - AP Northeast 1 (Tokyo)
- `https://s3.ap-northeast-2.wasabisys.com` - AP Northeast 2 (Osaka)
- `https://s3.ap-southeast-1.wasabisys.com` - AP Southeast 1 (Singapore)
- `https://s3.ap-southeast-2.wasabisys.com` - AP Southeast 2 (Sydney)

### Benefits

- 💰 **80% cheaper** than AWS S3
- 🚀 Fast performance
- ✅ S3 compatible
- 📦 No egress fees

---

## Cloudflare R2

Cloudflare R2 is zero-egress-fee S3-compatible storage.

### Configuration

```env
AWS_ACCESS_KEY_ID=your_cloudflare_access_key_id
AWS_SECRET_ACCESS_KEY=your_cloudflare_secret_access_key
AWS_REGION=auto
AWS_S3_BUCKET=my-bucket
AWS_S3_ENDPOINT=https://<account_id>.r2.cloudflarestorage.com
```

Replace `<account_id>` with your Cloudflare Account ID.

### Benefits

- 💰 **Zero egress fees** (free data transfer out)
- 🌐 Cloudflare's global network
- ✅ S3 compatible
- 🔒 Integrated with Cloudflare security

### Public Access

For public buckets, files are available at:
```
https://pub-<hash>.r2.dev/uploads/file.jpg
```

---

## Backblaze B2

Backblaze B2 with S3-compatible API.

### Configuration

```env
AWS_ACCESS_KEY_ID=your_application_key_id
AWS_SECRET_ACCESS_KEY=your_application_key
AWS_REGION=us-west-004
AWS_S3_BUCKET=my-bucket
AWS_S3_ENDPOINT=https://s3.us-west-004.backblazeb2.com
```

### Available Regions

- `s3.us-west-000.backblazeb2.com` - US West (Phoenix, AZ)
- `s3.us-west-001.backblazeb2.com` - US West (Sacramento, CA)
- `s3.us-west-002.backblazeb2.com` - US West (Phoenix, AZ)
- `s3.us-west-004.backblazeb2.com` - US West (Sacramento, CA)
- `s3.eu-central-003.backblazeb2.com` - EU Central (Amsterdam)

### Benefits

- 💰 Very affordable pricing
- ✅ S3 compatible
- 📦 Free egress up to 3x storage
- 🔒 Built-in encryption

---

## Local Development with Docker

### docker-compose.minio.yml

Create this file for local MinIO development:

```yaml
version: '3.8'

services:
  minio:
    image: minio/minio:latest
    container_name: minio
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"  # API
      - "9001:9001"  # Console (Web UI)
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio_data:/data
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Create default bucket on startup
  minio-init:
    image: minio/mc:latest
    depends_on:
      - minio
    entrypoint: >
      /bin/sh -c "
      until /usr/bin/mc alias set minio http://minio:9000 minioadmin minioadmin; do
        echo 'Waiting for MinIO...'
        sleep 1
      done;
      /usr/bin/mc mb minio/my-bucket --ignore-existing;
      /usr/bin/mc anonymous set download minio/my-bucket;
      echo 'MinIO bucket created successfully'
      "
    networks:
      - app-network

volumes:
  minio_data:
    driver: local

networks:
  app-network:
    driver: bridge
```

### Start MinIO

```bash
docker-compose -f docker-compose.minio.yml up -d
```

### Access MinIO

- **API Endpoint**: http://localhost:9000
- **Web Console**: http://localhost:9001
- **Username**: minioadmin
- **Password**: minioadmin

### Environment Configuration for Local MinIO

```env
AWS_ACCESS_KEY_ID=minioadmin
AWS_SECRET_ACCESS_KEY=minioadmin
AWS_REGION=us-east-1
AWS_S3_BUCKET=my-bucket
AWS_S3_ENDPOINT=http://localhost:9000
```

---

## Complete Example: Switching Between Providers

### Configuration Manager

```go
package config

import (
    "os"
    "starter-gofiber/pkg/apierror"
)

type StorageProvider string

const (
    ProviderAWS         StorageProvider = "aws"
    ProviderMinIO       StorageProvider = "minio"
    ProviderSpaces      StorageProvider = "spaces"
    ProviderWasabi      StorageProvider = "wasabi"
    ProviderR2          StorageProvider = "r2"
    ProviderBackblaze   StorageProvider = "backblaze"
)

func GetStorageConfig() helper.S3Config {
    provider := StorageProvider(os.Getenv("STORAGE_PROVIDER"))
    
    config := helper.S3Config{
        AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
        SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
        Region:    os.Getenv("AWS_REGION"),
        Bucket:    os.Getenv("AWS_S3_BUCKET"),
    }
    
    // Set endpoint based on provider
    switch provider {
    case ProviderMinIO:
        config.Endpoint = os.Getenv("AWS_S3_ENDPOINT") // http://localhost:9000
    case ProviderSpaces:
        config.Endpoint = os.Getenv("AWS_S3_ENDPOINT") // https://nyc3.digitaloceanspaces.com
    case ProviderWasabi:
        config.Endpoint = os.Getenv("AWS_S3_ENDPOINT") // https://s3.wasabisys.com
    case ProviderR2:
        config.Endpoint = os.Getenv("AWS_S3_ENDPOINT") // https://<account>.r2.cloudflarestorage.com
    case ProviderBackblaze:
        config.Endpoint = os.Getenv("AWS_S3_ENDPOINT") // https://s3.us-west-004.backblazeb2.com
    case ProviderAWS:
        // Leave endpoint empty for AWS S3
        config.Endpoint = ""
    default:
        config.Endpoint = os.Getenv("AWS_S3_ENDPOINT")
    }
    
    return config
}
```

### Updated .env

```env
# Storage Configuration
STORAGE_PROVIDER=minio  # aws|minio|spaces|wasabi|r2|backblaze

# Credentials (same for all S3-compatible services)
AWS_ACCESS_KEY_ID=minioadmin
AWS_SECRET_ACCESS_KEY=minioadmin
AWS_REGION=us-east-1
AWS_S3_BUCKET=my-bucket

# Endpoint (required for non-AWS providers)
AWS_S3_ENDPOINT=http://localhost:9000
```

---

## Performance Comparison

| Provider | Upload Speed | Download Speed | Pricing | Egress Fees |
|----------|-------------|----------------|---------|-------------|
| AWS S3 | Fast | Fast | $0.023/GB | $0.09/GB |
| MinIO (Local) | **Fastest** | **Fastest** | Free | Free |
| DigitalOcean Spaces | Fast | Fast | $5/250GB | $0.01/GB |
| Wasabi | Fast | Fast | $6.99/TB | **Free** |
| Cloudflare R2 | Fast | Fast | $0.015/GB | **Free** |
| Backblaze B2 | Medium | Medium | $5/TB | Free (3x) |

---

## Best Practices

1. **Local Development**: Use MinIO for local testing
2. **Production**: Choose based on your needs:
   - **Cost-effective**: Wasabi or Backblaze B2
   - **No egress fees**: Cloudflare R2 or Wasabi
   - **Simple & reliable**: AWS S3
   - **Self-hosted**: MinIO (on your servers)
   - **DigitalOcean users**: DigitalOcean Spaces

3. **Testing**: Always test with MinIO locally before deploying to production
4. **Fallback**: Consider multi-cloud strategy with automatic failover
5. **Security**: Always use HTTPS endpoints in production (except local MinIO)

---

## Troubleshooting

### Connection Refused

```bash
# Check if MinIO is running
docker ps | grep minio

# Check MinIO logs
docker logs minio
```

### Bucket Not Found

```bash
# Create bucket using mc client
mc mb local/my-bucket
```

### Access Denied

```bash
# Set public policy
mc anonymous set download local/my-bucket

# Or set bucket policy programmatically
```

### SSL Certificate Errors

For local MinIO with HTTP (not HTTPS), make sure to use `http://` in endpoint, not `https://`.

---

## Migration Between Providers

To migrate files between storage providers:

```go
// Example: Migrate from MinIO to AWS S3
func MigrateStorage(sourceClient, destClient *helper.S3Client, prefix string) error {
    ctx := context.Background()
    
    // List all files from source
    files, err := sourceClient.ListFiles(ctx, prefix, 1000)
    if err != nil {
        return err
    }
    
    // Copy each file to destination
    for _, key := range files {
        // Download from source
        // Upload to destination
    }
    
    return nil
}
```

---

## Conclusion

The implementation supports **any S3-compatible storage service** through the `Endpoint` configuration. Simply set the appropriate endpoint URL and credentials in your environment variables.

For local development, **MinIO** is highly recommended as it's:
- Free and open source
- Fast and lightweight
- 100% S3 API compatible
- Easy to run with Docker
