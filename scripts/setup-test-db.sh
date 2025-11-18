#!/bin/bash

set -e

echo "Setting up test database..."

DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
TEST_DB_NAME="pr_system_test"

export PGPASSWORD="$DB_PASSWORD"

echo "Creating test database: $TEST_DB_NAME"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "DROP DATABASE IF EXISTS $TEST_DB_NAME;"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "CREATE DATABASE $TEST_DB_NAME;"

echo "Running migrations..."
migrate -path migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$TEST_DB_NAME?sslmode=disable" up

echo "✅ Test database setup complete!"
echo "Connection string: postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$TEST_DB_NAME?sslmode=disable"
