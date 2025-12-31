#!/bin/bash

# Database Restore Script
# This script restores a MySQL database from a backup file

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check if backup file is provided
if [ -z "$1" ]; then
    echo "❌ Error: No backup file specified"
    echo "Usage: $0 <backup_file.sql.gz>"
    echo ""
    echo "Available backups:"
    ls -lh ./backups/*.sql.gz 2>/dev/null || echo "  No backups found"
    exit 1
fi

BACKUP_FILE="$1"

# Check if file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Error: Backup file not found: $BACKUP_FILE"
    exit 1
fi

echo "⚠️  WARNING: This will REPLACE all data in database: $DB_NAME"
echo "Backup file: $BACKUP_FILE"
read -p "Are you sure you want to continue? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "❌ Restore cancelled"
    exit 0
fi

echo "🔄 Starting database restore..."

# Decompress if needed
if [[ "$BACKUP_FILE" == *.gz ]]; then
    echo "📦 Decompressing backup..."
    TEMP_FILE="${BACKUP_FILE%.gz}"
    gunzip -c "$BACKUP_FILE" > "$TEMP_FILE"
    RESTORE_FILE="$TEMP_FILE"
    CLEANUP_TEMP=true
else
    RESTORE_FILE="$BACKUP_FILE"
    CLEANUP_TEMP=false
fi

# Drop and recreate database
echo "🗑️  Dropping existing database..."
mysql -h "$DB_URL" -u "$DB_USER" -p"$DB_PASS" -e "DROP DATABASE IF EXISTS $DB_NAME;"

echo "📝 Creating new database..."
mysql -h "$DB_URL" -u "$DB_USER" -p"$DB_PASS" -e "CREATE DATABASE $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# Restore backup
echo "📥 Restoring database..."
mysql -h "$DB_URL" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" < "$RESTORE_FILE"

# Cleanup temp file
if [ "$CLEANUP_TEMP" = true ]; then
    rm "$TEMP_FILE"
fi

echo "✅ Database restored successfully!"
echo "Database: $DB_NAME"
echo "✨ Done!"
