#!/bin/sh
set -e

echo "=== Starting Application ==="

echo "Waiting for PostgreSQL to be ready..."
until pg_isready; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 2
done
echo "PostgreSQL is ready!"

echo "Running database migrations..."
goose up
echo "Migrations completed successfully!"

echo "Loading seed data..."
psql -f /app/seeds/cli_seed.sql
echo "Seed data loaded successfully!"

echo "Starting main application..."
exec /usr/local/bin/ranking_system serve --config configs/config.yaml