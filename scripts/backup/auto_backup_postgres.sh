#!/bin/bash

# PostgreSQL Automated Backup Script (Cron-friendly)
# Runs backup and logs output

# Load environment variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$PROJECT_DIR"

if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Configuration
LOG_FILE="${LOG_FILE:-./logs/backup_postgres.log}"
BACKUP_DIR="${BACKUP_DIR:-./backups}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"
DB_NAME="${DB_NAME}"

# Create log directory
mkdir -p "$(dirname "$LOG_FILE")"

# Log function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

log "========================================="
log "Starting automated PostgreSQL backup"
log "Database: $DB_NAME"

# Run backup script
"$SCRIPT_DIR/backup_postgres.sh" >> "$LOG_FILE" 2>&1

if [ $? -eq 0 ]; then
    log "Backup completed successfully"
    
    # Calculate backup statistics
    BACKUP_COUNT=$(find "$BACKUP_DIR" -name "backup_${DB_NAME}_*.sql.gz" | wc -l)
    TOTAL_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)
    
    log "Total backups: $BACKUP_COUNT"
    log "Total size: $TOTAL_SIZE"
    log "Retention: $RETENTION_DAYS days"
else
    log "❌ Backup failed! Check logs for details"
    
    # Optional: Send notification via webhook
    # curl -X POST https://your-webhook-url.com \
    #      -H "Content-Type: application/json" \
    #      -d "{\"message\": \"PostgreSQL backup failed for $DB_NAME\"}"
fi

log "========================================="
