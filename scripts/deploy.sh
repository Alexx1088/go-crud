#!/bin/bash

# Exit on any error
set -e

echo "Applying database migrations..."

# Set the PostgreSQL password as an environment variable
export PGPASSWORD="$DB_PASSWORD"

# Apply all migration files in the migrations directory
for migration in /app/migrations/*.sql; do
    echo "Running migration: $migration"
    psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "$migration"
done

echo "Migrations applied successfully."