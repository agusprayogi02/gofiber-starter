#!/bin/bash

# PostgreSQL Database Backup Script
# Creates compressed backup with automatic cleanup

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Configuration
BACKUP_DIR="${BACKUP_DIR:-./backups}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
DB_HOST="${DB_URL%%:*}"
DB_PORT="${DB_URL##*:}"
DB_NAME="${DB_NAME}"
DB_USER="${DB_USER}"
PGPASSWORD="${DB_PASS}"
BACKUP_FILE="${BACKUP_DIR}/backup_${DB_NAME}_${TIMESTAMP}.sql"
RETENTION_DAYS=7

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Export password for pg_dump
export PGPASSWORD

echo "🔄 Starting PostgreSQL database backup..."
echo "Database: $DB_NAME"
echo "Backup file: ${BACKUP_FILE}"

# Perform backup
pg_dump -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        --clean \
        --if-exists \
        --create \
        --no-owner \
        --no-acl \
        > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "📦 Compressing backup..."
    gzip "$BACKUP_FILE"
    
    echo "✅ Backup completed successfully!"
    echo "File: ${BACKUP_FILE}.gz"
    echo "Size: $(du -h "${BACKUP_FILE}.gz" | cut -f1)"
    
    # Cleanup old backups
    echo "🧹 Cleaning up backups older than ${RETENTION_DAYS} days..."
    find "$BACKUP_DIR" -name "backup_${DB_NAME}_*.sql.gz" -mtime +${RETENTION_DAYS} -delete
    
else
    echo "❌ Backup failed!"
    exit 1
fi

# Unset password
unset PGPASSWORD
