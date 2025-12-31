#!/bin/bash

# Automated Daily Backup Script
# Add this to crontab for automatic daily backups

set -e

# Configuration
BACKUP_DIR="./backups"
RETENTION_DAYS=30
LOG_FILE="./logs/backup.log"

# Create log directory
mkdir -p "$(dirname "$LOG_FILE")"

# Log function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

log "Starting automated backup..."

# Run backup
if ./scripts/backup/backup.sh >> "$LOG_FILE" 2>&1; then
    log "✅ Backup completed successfully"
else
    log "❌ Backup failed"
    exit 1
fi

# Cleanup old backups
log "Cleaning up backups older than $RETENTION_DAYS days..."
find "$BACKUP_DIR" -name "backup_*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete

# Get backup count and total size
BACKUP_COUNT=$(find "$BACKUP_DIR" -name "backup_*.sql.gz" -type f | wc -l)
TOTAL_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)

log "Current backups: $BACKUP_COUNT files, Total size: $TOTAL_SIZE"
log "✨ Automated backup completed"

# Optional: Send notification (uncomment to enable)
# curl -X POST "https://your-webhook-url" \
#   -H "Content-Type: application/json" \
#   -d "{\"text\":\"Database backup completed: $BACKUP_COUNT backups, $TOTAL_SIZE\"}"
