#!/bin/bash
# Generate Atlas DB URL from .env variables
# This script reads existing ENV and constructs proper database URL

if [ ! -f .env ]; then
    echo "Error: .env file not found" >&2
    exit 1
fi

# Load .env
set -a
source .env
set +a

# Construct DSN based on DB_TYPE
case $DB_TYPE in
    mysql)
        echo "mysql://$DB_USER:$DB_PASS@$DB_URL/$DB_NAME?parseTime=true"
        ;;
    postgres)
        echo "postgres://$DB_USER:$DB_PASS@$DB_URL/$DB_NAME?sslmode=disable"
        ;;
    sqlserver)
        echo "sqlserver://$DB_USER:$DB_PASS@$DB_URL?database=$DB_NAME"
        ;;
    *)
        # Default to postgres
        echo "postgres://$DB_USER:$DB_PASS@$DB_URL/$DB_NAME?sslmode=disable"
        ;;
esac
