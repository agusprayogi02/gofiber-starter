#!/bin/bash

# Database Backup Script
# This script creates a backup of the MySQL database

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Configuration
BACKUP_DIR="${BACKUP_DIR:-./backups}"
DATE=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/backup_${DB_NAME}_${DATE}.sql"
COMPRESSED_FILE="${BACKUP_FILE}.gz"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

echo "🔄 Starting database backup..."
echo "Database: $DB_NAME"
echo "Backup file: $BACKUP_FILE"

# Create backup
if [ "$DB_TYPE" = "mysql" ]; then
    mysqldump -h "$DB_URL" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" \
        --single-transaction \
        --quick \
        --lock-tables=false \
        --routines \
        --triggers \
        --events \
        > "$BACKUP_FILE"
else
    echo "❌ Unsupported database type: $DB_TYPE"
    exit 1
fi

# Compress backup
echo "📦 Compressing backup..."
gzip "$BACKUP_FILE"

# Calculate size
SIZE=$(du -h "$COMPRESSED_FILE" | cut -f1)

echo "✅ Backup completed successfully!"
echo "File: $COMPRESSED_FILE"
echo "Size: $SIZE"

# Cleanup old backups (keep last 7 days)
echo "🧹 Cleaning up old backups (keeping last 7 days)..."
find "$BACKUP_DIR" -name "backup_*.sql.gz" -type f -mtime +7 -delete

echo "✨ Done!"
